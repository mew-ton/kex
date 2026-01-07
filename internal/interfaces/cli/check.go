package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/mew-ton/kex/internal/infrastructure/config"
	"github.com/mew-ton/kex/internal/infrastructure/fs"
	"github.com/mew-ton/kex/internal/infrastructure/logger"
	"github.com/mew-ton/kex/internal/usecase/validator"

	"github.com/pterm/pterm"
	"github.com/urfave/cli/v2"
)

var CheckCommand = &cli.Command{
	Name:  "check",
	Usage: "Validate documents",
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:  "json",
			Usage: "Output results in JSON format",
		},
	},
	Action: runCheck,
}

func runCheck(c *cli.Context) error {
	isJSON := c.Bool("json")

	if !isJSON {
		pterm.DefaultSection.Println("Checking documents...")
	}

	// 1. Resolve Project Root
	projectRoot := c.Args().First()
	if projectRoot == "" {
		projectRoot = "."
	}

	cfg, err := resolveConfig(projectRoot)
	if err != nil {
		if !isJSON {
			pterm.Warning.Printf("Failed to load config, using defaults: %v\n", err)
		}
	}

	// 2. Resolve Content Directory
	root := filepath.Join(projectRoot, cfg.Root)

	repo, err := loadRepository(root, !isJSON)
	if err != nil {
		if isJSON {
			printJSONError(err.Error())
			return cli.Exit("", 1)
		}
		return cli.Exit(fmt.Sprintf("Fatal: failed to load documents: %v", err), 1)
	}

	report := validator.Validate(repo)

	if isJSON {
		printJSONReport(report)
	} else {
		printHumanReport(report, report.Stats.Total)
	}

	if !report.Valid {
		return cli.Exit("", 1)
	}

	return nil
}

func resolveConfig(projectRoot string) (config.Config, error) {
	return config.Load(projectRoot)
}

func loadRepository(root string, showSpinner bool) (*fs.Indexer, error) {
	var spinner *pterm.SpinnerPrinter
	if showSpinner {
		spinner, _ = pterm.DefaultSpinner.Start("Loading documents...")
	}

	// Use NoOpLogger for Check command to avoid clutter
	l := &logger.NoOpLogger{}
	provider := fs.NewLocalProvider(root, l)
	repo := fs.New(provider, l)
	repo.IncludeDrafts = true
	if err := repo.Load(); err != nil {
		if spinner != nil {
			spinner.Fail("Failed to load documents")
		}
		return nil, err
	}
	if spinner != nil {
		spinner.Success("Documents loaded")
	}
	return repo, nil
}

// Presentation Logic

func printJSONReport(report validator.ValidationReport) {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	encoder.Encode(report)
}

func printJSONError(msg string) {
	// Minimal JSON error output
	fmt.Printf(`{"valid": false, "global_errors": ["%s"]}`+"\n", strings.ReplaceAll(msg, "\"", "\\\""))
}

func printHumanReport(report validator.ValidationReport, totalDocs int) {
	// Print Global Errors
	for _, err := range report.GlobalErrors {
		pterm.Error.Printf("Parse Error: %v\n", err)
	}

	// Print Document Errors
	for _, doc := range report.Documents {
		for _, err := range doc.Errors {
			// Logic to print based on status is now presentation concern
			// We can assume ValidationReport contains necessary info
			// Re-mapping status for display
			if doc.Status == "draft" {
				pterm.Warning.Printf("[%s] (Draft) %s\n", doc.ID, err)
			} else {
				pterm.Error.Printf("[%s] (Adopted) %s\n", doc.ID, err)
			}
		}
	}

	// Statistics
	pterm.Println() // Spacer
	pterm.DefaultSection.Println("Statistics")

	tableData := [][]string{
		{"Metric", "Count"},
		{"Total Documents", fmt.Sprintf("%d", totalDocs)},
		{"Adopted", fmt.Sprintf("%d", report.Stats.Adopted)},
		{"Draft", fmt.Sprintf("%d", report.Stats.Draft)},
		{"Parse Errors", fmt.Sprintf("%d", report.Stats.ParseErrors)},
		{"Adopted Errors", fmt.Sprintf("%d", report.Stats.AdoptedErrors)},
		{"Draft Warnings", fmt.Sprintf("%d", report.Stats.DraftWarnings)},
	}

	pterm.DefaultTable.WithHasHeader().WithData(tableData).Render()
	pterm.Println()

	if !report.Valid {
		pterm.Error.Println("Check failed. Please fix errors.")
	} else {
		pterm.Success.Println("All checks passed.")
	}
}
