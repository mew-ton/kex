package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/mew-ton/kex/internal/infrastructure/config"
	"github.com/mew-ton/kex/internal/infrastructure/fs"
	"github.com/pterm/pterm"
	"github.com/urfave/cli/v2"
)

var GenerateCommand = &cli.Command{
	Name:   "generate",
	Usage:  "Generate static site (dist)",
	Flags:  []cli.Flag{},
	Action: runGenerate,
}

func runGenerate(c *cli.Context) error {
	pterm.DefaultSection.Println("Generating static site...")

	// 1. Resolve Repository
	repo, cfg, _, err := resolveRepository(c)
	if err != nil {
		pterm.Error.Printf("Failed to load documents: %v\n", err)
		return cli.Exit("", 1)
	}

	// 2. Export Schema
	schema, err := repo.Export()
	if err != nil {
		return cli.Exit(fmt.Sprintf("Failed to export: %v", err), 1)
	}

	// 3. Clean/Create dist
	outputDir := "dist"
	os.RemoveAll(outputDir)
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return cli.Exit(fmt.Sprintf("Failed to create dist: %v", err), 1)
	}

	// 4. Copy Files
	if err := copyDocuments(repo, schema, outputDir); err != nil {
		return cli.Exit(err.Error(), 1)
	}

	// 5. Transform Schema Paths if BaseURL is set
	transformPaths(cfg, schema)

	// 6. Write kex.json
	if err := writeIndex(schema, outputDir); err != nil {
		return cli.Exit(err.Error(), 1)
	}

	pterm.Success.Println("Generated static site in 'dist/'")
	return nil
}

func copyDocuments(repo *fs.Indexer, schema *fs.IndexSchema, outputDir string) error {
	copySpinner, _ := pterm.DefaultSpinner.Start("Copying files...")

	for _, doc := range schema.Documents {
		if err := copyDocument(repo, doc, outputDir); err != nil {
			copySpinner.Fail(fmt.Sprintf("Failed to copy %s: %v", doc.Path, err))
			return err
		}
	}
	copySpinner.Success("Files copied")
	return nil
}

func copyDocument(repo *fs.Indexer, doc *fs.DocumentSchema, outputDir string) error {
	// Destination path
	dstPath := filepath.Join(outputDir, doc.Path)

	if err := os.MkdirAll(filepath.Dir(dstPath), 0755); err != nil {
		return fmt.Errorf("failed to create dir: %w", err)
	}

	d, err := os.Create(dstPath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer d.Close()

	// Fetch Content via Repo/Provider
	domainDoc, ok := repo.GetByID(doc.ID)
	if !ok || domainDoc == nil {
		// Should not happen as schema comes from repo
		return nil
	}

	// Ensure Body is there and re-fetch if necessary
	if domainDoc.Body == "" {
		fetchedDoc, found := repo.GetByID(doc.ID)
		if found {
			domainDoc = fetchedDoc
		}
	}

	if _, err := d.WriteString(domainDoc.Body); err != nil {
		return fmt.Errorf("failed to write content: %w", err)
	}
	return nil
}

func transformPaths(cfg config.Config, schema *fs.IndexSchema) {
	if cfg.BaseURL != "" {
		for _, doc := range schema.Documents {
			base := cfg.BaseURL
			if base[len(base)-1] != '/' {
				base += "/"
			}
			doc.Path = base + doc.Path
		}
	}
}

func writeIndex(schema *fs.IndexSchema, outputDir string) error {
	f, err := os.Create(filepath.Join(outputDir, "kex.json"))
	if err != nil {
		return fmt.Errorf("failed to create kex.json: %w", err)
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	if err := enc.Encode(schema); err != nil {
		return fmt.Errorf("failed to encode kex.json: %w", err)
	}
	return nil
}
