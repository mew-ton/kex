package cli

import (
	"encoding/json"
	"fmt"
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
	// Scan single source
	l := logger.NewStderrLogger()
	spinner, _ := pterm.DefaultSpinner.Start("Scanning documents...")

	source := cfg.Source
	root := filepath.Join(projectRoot, source)
	if _, err := os.Stat(root); os.IsNotExist(err) {
		return nil, fmt.Errorf("source directory %s does not exist", source)
	}

	provider := fs.NewLocalProvider(root, l)
	repo := fs.New(provider, l)

	if err := repo.Load(); err != nil {
		return nil, fmt.Errorf("failed to load documents: %w", err)
	}

	schema, err := repo.Export()
	if err != nil {
		return nil, fmt.Errorf("failed to export schema: %w", err)
	}

	spinner.Success("Documents scanned")

	return schema, nil
}

func prepareDistDir(outputDir string) error {
	if err := os.RemoveAll(outputDir); err != nil {
		return fmt.Errorf("failed to clean dist: %w", err)
	}
	if err := fs.EnsureDir(outputDir); err != nil {
		return fmt.Errorf("failed to create dist: %w", err)
	}
	return nil
}

func copyContents(projectRoot, outputDir string, cfg config.Config, schema *fs.IndexSchema) error {
	copySpinner, _ := pterm.DefaultSpinner.Start("Copying files...")

	source := cfg.Source

	for _, doc := range schema.Documents {
		srcPath := filepath.Join(projectRoot, source, doc.Path)
		dstPath := filepath.Join(outputDir, doc.Path)

		if err := fs.CopyFile(srcPath, dstPath); err != nil {
			copySpinner.Fail(fmt.Sprintf("Failed to copy %s: %v", srcPath, err))
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
