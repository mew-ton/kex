package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"kex/internal/domain"
	"kex/internal/indexer"

	"github.com/urfave/cli/v2"
)

var CheckCommand = &cli.Command{
	Name:   "check",
	Usage:  "Validate documents",
	Action: runCheck,
}

func runCheck(c *cli.Context) error {
	// 1. Resolve configuration (root directory)
	// TODO: Load .kex.yaml properly. For now we assume defaults or check contents/
	root := "contents"
	configPath := ".kex.yaml" // Just check existence check?
	if _, err := os.Stat(configPath); err == nil {
		// Parse config if we had a config parser.
		// For v1 we can just default to contents/ unless we really implement config loading.
	}

	if _, err := os.Stat(root); os.IsNotExist(err) {
		return cli.Exit(fmt.Sprintf("Error: directory '%s' not found. Run 'kex init'?", root), 1)
	}

	fmt.Printf("Checking documents in: %s\n", root)

	// 2. Load Indexer (Build Phase)
	idx := indexer.New(root)
	if err := idx.Load(); err != nil {
		return cli.Exit(fmt.Sprintf("Fatal: failed to load documents: %v", err), 1)
	}

	// 3. Validation Logic
	errorCount := 0

	// 3.1 Report Parsing Errors
	for _, err := range idx.Errors {
		fmt.Printf("Error: %v\n", err)
		errorCount++
	}

	// 3.2 Validate Valid Documents
	for _, doc := range idx.Documents {
		if err := validateDocument(doc); err != nil {
			fmt.Printf("Error: [%s] %v\n", doc.Path, err)
			errorCount++
		}
	}

	if errorCount > 0 {
		return cli.Exit(fmt.Sprintf("\nFound %d errors.", errorCount), 1)
	}

	fmt.Println("All checks passed.")
	return nil
}

func validateDocument(doc *domain.Document) error {
	// Check Status
	if doc.Status == domain.StatusDraft {
		return fmt.Errorf("document is in 'draft' status")
	}

	// Check Filename matches ID
	filename := filepath.Base(doc.Path)
	ext := filepath.Ext(filename)
	basename := strings.TrimSuffix(filename, ext)

	if basename != doc.ID {
		return fmt.Errorf("filename must match id (filename: %s, id: %s)", filename, doc.ID)
	}

	return nil
}
