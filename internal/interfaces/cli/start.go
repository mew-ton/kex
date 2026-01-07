package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/mew-ton/kex/internal/infrastructure/config"
	"github.com/mew-ton/kex/internal/infrastructure/fs"
	"github.com/mew-ton/kex/internal/interfaces/mcp"
	"github.com/mew-ton/kex/internal/usecase/validator"

	"github.com/urfave/cli/v2"
)

var StartCommand = &cli.Command{
	Name:  "start",
	Usage: "Start MCP server",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "root",
			Usage:   "Path to guidelines directory",
			Aliases: []string{"r"},
		},
	},
	Action: runStart,
}

func runStart(c *cli.Context) error {
	fmt.Fprintf(os.Stderr, "Starting Kex Server...\n")

	// 1. Resolve Project Root
	projectRoot := c.Args().First()
	if projectRoot == "" {
		projectRoot = "."
	}

	// 2. Resolve configuration
	cfg, err := config.Load(projectRoot)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: failed to load config: %v\n", err)
	}

	// 3. Resolve Content Directory
	// Default: projectRoot/cfg.Root
	// Override: --root flag
	root := filepath.Join(projectRoot, cfg.Root)
	if c.String("root") != "" {
		// If --root is absolute, use it directly.
		// If relative, treat it as relative to CWD (standard CLI behavior)
		// Or should it be relative to projectRoot?
		// Issue says "--root option should specify CWD", which is confusing.
		// But "standard" way is: --root overrides everything.
		// Let's assume --root is explicit path.
		root = c.String("root")
	}

	if _, err := os.Stat(root); os.IsNotExist(err) {
		return cli.Exit(fmt.Sprintf("Error: directory '%s' not found. Run 'kex init'?", root), 1)
	}

	// 4. Load Indexer (Infrastructure)
	repo := fs.New(root)
	if err := repo.Load(); err != nil {
		return cli.Exit(fmt.Sprintf("Fatal: failed to load documents: %v", err), 1)
	}

	// 5. Strict validation on start (Use Case)
	// We use the Validator use case to determine validity
	report := validator.Validate(repo)

	// Check for Parse Errors (Critical for start)
	if len(report.GlobalErrors) > 0 {
		for _, e := range report.GlobalErrors {
			fmt.Fprintf(os.Stderr, "Error: %v\n", e)
		}
		return cli.Exit("Failed to start due to document errors. Run 'kex check' for details.", 1)
	}

	// Note: We might allow "AdoptedErrors" but block on "ParseErrors".
	// For strictness, let's block if Valid is false (excluding Drafts which are warnings)
	// But validator.Validate returns Valid=false if AdoptedErrors > 0.
	// As per previous design: "Failed to start due to document errors".
	if !report.Valid {
		return cli.Exit("Validation failed (documents contain errors). Run 'kex check' for details.", 1)
	}

	// 4. Start Server (Interface)
	srv := mcp.New(repo)
	fmt.Fprintf(os.Stderr, "Server listening on stdio...\n")
	if err := srv.Serve(); err != nil {
		return cli.Exit(fmt.Sprintf("Server error: %v", err), 1)
	}

	return nil
}
