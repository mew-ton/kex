package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"kex/internal/config"
	"kex/internal/domain"
	"kex/internal/indexer"

	"github.com/pterm/pterm"
	"github.com/urfave/cli/v2"
)

var CheckCommand = &cli.Command{
	Name:   "check",
	Usage:  "Validate documents",
	Action: runCheck,
}

func runCheck(c *cli.Context) error {
	pterm.DefaultSection.Println("Checking documents...")

	// 1. Resolve configuration (root directory)
	// 1. Resolve configuration
	cfg, err := config.Load()
	if err != nil {
		pterm.Warning.Printf("Failed to load config, using defaults: %v\n", err)
	}
	root := cfg.Root

	if _, err := os.Stat(root); os.IsNotExist(err) {
		pterm.Error.Printf("Error: directory '%s' not found. Run 'kex init'?\n", root)
		return cli.Exit("", 1)
	}

	// 2. Load Indexer
	spinner, _ := pterm.DefaultSpinner.Start("Loading documents...")
	idx := indexer.New(root)
	idx.IncludeDrafts = true
	if err := idx.Load(); err != nil {
		spinner.Fail("Failed to load documents")
		return cli.Exit(fmt.Sprintf("Fatal: failed to load documents: %v", err), 1)
	}
	spinner.Success("Documents loaded")

	// 3. Validation Logic
	stats := struct {
		Adopted       int
		Draft         int
		AdoptedErrors int
		DraftWarnings int // Draft errors are warnings
		ParseErrors   int
	}{}

	// 3.1 Report Parsing Errors (These are fatal to the specific file)
	for _, err := range idx.Errors {
		pterm.Error.Printf("Parse Error: %v\n", err)
		stats.ParseErrors++
	}

	// 3.2 Validate Documents
	for _, doc := range idx.Documents {
		if doc.Status == domain.StatusDraft {
			stats.Draft++
		} else {
			stats.Adopted++
		}

		if err := validateDocument(doc); err != nil {
			if doc.Status == domain.StatusDraft {
				pterm.Warning.Printf("[%s] (Draft) %v\n", doc.ID, err)
				stats.DraftWarnings++
			} else {
				pterm.Error.Printf("[%s] (Adopted) %v\n", doc.ID, err)
				stats.AdoptedErrors++
			}
		}
	}

	// 4. Statistics Table
	pterm.Println() // Spacer
	pterm.DefaultSection.Println("Statistics")

	tableData := [][]string{
		{"Metric", "Count"},
		{"Total Documents", fmt.Sprintf("%d", len(idx.Documents))},
		{"Adopted", fmt.Sprintf("%d", stats.Adopted)},
		{"Draft", fmt.Sprintf("%d", stats.Draft)},
		{"Parse Errors", fmt.Sprintf("%d", stats.ParseErrors)},
		{"Adopted Errors", fmt.Sprintf("%d", stats.AdoptedErrors)},
		{"Draft Warnings", fmt.Sprintf("%d", stats.DraftWarnings)},
	}

	// Render table
	// We use a simple table render
	pterm.DefaultTable.WithHasHeader().WithData(tableData).Render()

	pterm.Println()

	// 5. Exit Code Logic
	if stats.AdoptedErrors > 0 || stats.ParseErrors > 0 {
		pterm.Error.Println("Check failed. Please fix errors.")
		return cli.Exit("", 1)
	}

	pterm.Success.Println("All checks passed.")
	return nil
}

func validateDocument(doc *domain.Document) error {
	// Note: We no longer fail simply for being Draft.
	// We only check for inconsistencies (like filename mismatch)

	// Check Filename matches ID
	filename := filepath.Base(doc.Path)
	ext := filepath.Ext(filename)
	basename := strings.TrimSuffix(filename, ext)

	if basename != doc.ID {
		return fmt.Errorf("filename must match id (filename: %s, id: %s)", filename, doc.ID)
	}

	return nil
}
