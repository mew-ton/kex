package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/mew-ton/kex/internal/infrastructure/config"
	"github.com/mew-ton/kex/internal/infrastructure/fs"
	"github.com/mew-ton/kex/internal/infrastructure/logger"
	"github.com/pterm/pterm"
	"github.com/urfave/cli/v2"
)

var GenerateCommand = &cli.Command{
	Name:   "generate",
	Usage:  "Generate static site (dist)",
	Action: runGenerate,
}

func runGenerate(c *cli.Context) error {
	pterm.DefaultSection.Println("Generating static site...")

	// 1. Resolve Config
	projectRoot := c.Args().First()
	if projectRoot == "" {
		projectRoot = "."
	}
	cfg := loadConfig(projectRoot)

	// 2. Scan Documents
	schema, err := scanDocuments(projectRoot, cfg)
	if err != nil {
		return cli.Exit(err.Error(), 1)
	}

	// 3. Prepare Dist Directory
	outputDir := filepath.Join(projectRoot, "dist")
	if err := prepareDistDir(outputDir); err != nil {
		return cli.Exit(err.Error(), 1)
	}

	// 4. Copy Contents
	if err := copyContents(projectRoot, outputDir, cfg, schema); err != nil {
		return cli.Exit(err.Error(), 1)
	}

	// 5. Apply BaseURL (if needed) and Write Manifest
	if err := writeManifest(outputDir, cfg, schema); err != nil {
		return cli.Exit(err.Error(), 1)
	}

	pterm.Success.Println("Generated static site in 'dist/'")
	return nil
}

// Helpers

func loadConfig(projectRoot string) config.Config {
	cfg, err := config.Load(projectRoot)
	if err != nil {
		pterm.Warning.Printf("Failed to load config: %v\n", err)
		return config.Config{} // Return empty/default config
	}
	return cfg
}

func scanDocuments(projectRoot string, cfg config.Config) (*fs.IndexSchema, error) {
	// For generate, we might want to support multiple sources?
	// For now, let's assume we scan all sources defined?
	// But copyContents assumes a single structure?
	// The CLI "generate" is typically for generating a static site from specific docs.
	// Issue #34 says "multiplexing" is for start (MCP).
	// kex generate documentation usually processes one set of docs.
	// But if config defines multiple sources...
	// Let's iterate all sources and aggregate?
	// or just take the first one as default "main" docs?
	// Given "generate" creates a site structure, multiplexing sources into one site
	// would require collision handling etc.
	// For backward compatibility and simplicity, let's iterate all and merge into one list,
	// checking for collisions or simply overriding.

	// Scan all sources
	var collectedDocs []*fs.DocumentSchema

	l := logger.NewStderrLogger()
	spinner, _ := pterm.DefaultSpinner.Start("Scanning documents...")

	for _, source := range cfg.Sources {
		root := filepath.Join(projectRoot, source)
		if _, err := os.Stat(root); os.IsNotExist(err) {
			continue
		}

		provider := fs.NewLocalProvider(root, l)
		repo := fs.New(provider, l)

		if err := repo.Load(); err != nil {
			// Warn but maybe continue?
			pterm.Warning.Printf("Failed to load source %s: %v\n", source, err)
			continue
		}

		schema, err := repo.Export()
		if err != nil {
			pterm.Warning.Printf("Failed to export schema for %s: %v\n", source, err)
			continue
		}
		collectedDocs = append(collectedDocs, schema.Documents...)
	}

	spinner.Success("Documents scanned")

	// Construct aggregate schema
	schema := &fs.IndexSchema{
		Documents: collectedDocs,
	}

	return schema, nil
}

func prepareDistDir(outputDir string) error {
	if err := os.RemoveAll(outputDir); err != nil {
		return fmt.Errorf("failed to clean dist: %w", err)
	}
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create dist: %w", err)
	}
	return nil
}

func copyContents(projectRoot, outputDir string, cfg config.Config, schema *fs.IndexSchema) error {
	copySpinner, _ := pterm.DefaultSpinner.Start("Copying files...")

	// To copy correctly, we need to know WHICH source a doc came from.
	// But schema.Documents only has relative Path (e.g. "foo.md").
	// And we are flattening everything into Dist.
	// If we have multiple sources:
	// source1/foo.md
	// source2/bar.md
	// We need to look up where "foo.md" exists.

	for _, doc := range schema.Documents {
		// Find source for this doc
		// Check each source for the file
		var foundSrc string
		for _, source := range cfg.Sources {
			srcPath := filepath.Join(projectRoot, source, doc.Path)
			if _, err := os.Stat(srcPath); err == nil {
				foundSrc = srcPath
				break
			}
		}

		if foundSrc == "" {
			pterm.Warning.Printf("File not found for %s\n", doc.Path)
			continue
		}

		dstPath := filepath.Join(outputDir, doc.Path)

		if err := copyFile(foundSrc, dstPath); err != nil {
			copySpinner.Fail(fmt.Sprintf("Failed to copy %s: %v", foundSrc, err))
			return fmt.Errorf("failed to copy file: %w", err)
		}
	}
	copySpinner.Success("Files copied")
	return nil
}

func writeManifest(outputDir string, cfg config.Config, schema *fs.IndexSchema) error {
	// Transform Schema Paths if BaseURL is set
	if cfg.BaseURL != "" {
		base := cfg.BaseURL
		if base[len(base)-1] != '/' {
			base += "/"
		}
		for _, doc := range schema.Documents {
			doc.Path = base + doc.Path
		}
	}

	f, err := os.Create(filepath.Join(outputDir, "kex.json"))
	if err != nil {
		return fmt.Errorf("failed to create kex.json: %v", err)
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	if err := enc.Encode(schema); err != nil {
		return fmt.Errorf("failed to encode kex.json: %w", err)
	}
	return nil
}

func copyFile(src, dst string) error {
	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return err
	}

	s, err := os.Open(src)
	if err != nil {
		return err
	}
	defer s.Close()

	d, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer d.Close()

	_, err = io.Copy(d, s)
	return err
}
