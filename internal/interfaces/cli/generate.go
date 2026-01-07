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

	// 1. Resolve Project Root
	projectRoot := c.Args().First()
	if projectRoot == "" {
		projectRoot = "."
	}

	cfg, err := config.Load(projectRoot)
	if err != nil {
		pterm.Warning.Printf("Failed to load config: %v\n", err)
	}

	root := filepath.Join(projectRoot, cfg.Root)
	outputDir := filepath.Join(projectRoot, "dist")

	// 2. Load Indexer (Local Scan)
	spinner, _ := pterm.DefaultSpinner.Start("Scanning documents...")
	l := logger.NewStderrLogger()
	provider := fs.NewLocalProvider(root, l)
	repo := fs.New(provider, l)
	if err := repo.Load(); err != nil {
		spinner.Fail(fmt.Sprintf("Failed to load: %v", err))
		return cli.Exit("", 1)
	}
	spinner.Success("Documents scanned")

	// 3. Export Schema
	schema, err := repo.Export()
	if err != nil {
		return cli.Exit(fmt.Sprintf("Failed to export: %v", err), 1)
	}

	// 4. Apply BaseURL Transformation if needed
	// Note: We modifying Schema in place or treating paths?
	// Schema.Documents[i].Path is relative.
	// If Config.BaseURL is set, we might want to make them absolute URLs in kex.json
	// BUT, if we do that, mirroring files logic needs to know original relative path.
	// Export() returns relative paths (or abs if failed).
	// Let's copy files first, then update schema paths if needed.

	// 5. Clean/Create dist
	os.RemoveAll(outputDir)
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return cli.Exit(fmt.Sprintf("Failed to create dist: %v", err), 1)
	}

	// 6. Copy Files
	copySpinner, _ := pterm.DefaultSpinner.Start("Copying files...")
	for _, doc := range schema.Documents {
		// doc.Path is relative to cfg.Root usually
		srcPath := filepath.Join(root, doc.Path)
		dstPath := filepath.Join(outputDir, "contents", doc.Path) // Mirror structure under dist/contents?
		// Or mirror under dist/ root?
		// Issue requirement says: "dist ディレクトリを作り、そのルートには kex.json を置く. markdown は配下に構造維持してコピーする"
		// If doc.Path involves folders, we mirror them.
		// Let's mirror starting from root.
		// If doc.Path is "coding/style.md", dst is "dist/coding/style.md".
		dstPath = filepath.Join(outputDir, doc.Path)

		if err := copyFile(srcPath, dstPath); err != nil {
			copySpinner.Fail(fmt.Sprintf("Failed to copy %s: %v", srcPath, err))
			return cli.Exit("", 1)
		}
	}
	copySpinner.Success("Files copied")

	// 7. Transform Schema Paths if BaseURL is set
	if cfg.BaseURL != "" {
		for _, doc := range schema.Documents {
			// doc.Path becomes absolute URL
			// Ensure BaseURL doesn't end with slash and Path doesn't start with slash (relative)
			// But ensure we process all.
			// Is doc.Path URL encoded? Assume simple paths.
			// doc.Path is "coding/style.md"
			// URL: BaseURL + "/" + doc.Path
			// Or just simple string concat.
			// Use simple string concat ensuring separator.
			base := cfg.BaseURL
			if base[len(base)-1] != '/' {
				base += "/"
			}
			doc.Path = base + doc.Path
		}
	}

	// 8. Write kex.json
	f, err := os.Create(filepath.Join(outputDir, "kex.json"))
	if err != nil {
		return cli.Exit(fmt.Sprintf("Failed to create kex.json: %v", err), 1)
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	if err := enc.Encode(schema); err != nil {
		return cli.Exit("Failed to encode kex.json", 1)
	}

	pterm.Success.Println("Generated static site in 'dist/'")
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
