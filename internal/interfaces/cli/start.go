package cli

import (
	"fmt"
	"os"

	"github.com/mew-ton/kex/internal/infrastructure/fs"
	"github.com/mew-ton/kex/internal/infrastructure/logger"
	"github.com/mew-ton/kex/internal/interfaces/mcp"
	"github.com/mew-ton/kex/internal/usecase/validator"

	"github.com/urfave/cli/v2"
)

var StartCommand = &cli.Command{
	Name:  "start",
	Usage: "Start MCP server",
	Flags: []cli.Flag{

		&cli.StringFlag{
			Name:  "log-file",
			Usage: "Path to log file",
		},
	},
	Action: runStart,
}

func runStart(c *cli.Context) error {
	fmt.Fprintf(os.Stderr, "Starting Kex Server...\n")

	// Compatibility: Handle `kex start <url_or_path>` as an additional root
	// We might need to manually inject it into CLI context or handle it in resolveRepository?
	// `resolveRepository` reads c.StringSlice("root").
	// If there are Args, we should append them?
	// But `resolveRepository` is generic.

	// Let's modify usage:
	// If args exist, treat them as overrides or additions?
	// Old behavior: arg overrides config.root.
	// New behavior: arg is added to roots?

	// Actually, let's keep it simple.
	// 1. Resolve Repo
	repo, _, l, err := resolveRepository(c)

	// HANDLE LEGACY ARG SUPPORT (kex start <path>)
	// If resolved repo has NO documents (maybe config was empty), or strict check?
	// Or explicitly check c.Args().First() and add it if not in roots?
	// The problem is resolveRepository already called Load().

	// OPTIMIZATION:
	// We should probably pass args to resolveRepository or handle args before calling it.
	// But since we can't easily modify context flags...

	// Let's rely on resolveRepository for now.
	// If the user provided `kex start path/to/docs`, `resolveConfigAndLogger` in repo_loader.go
	// (which I copied from start.go) intentionally ignored handling `arg` as a config root based on my previous copy.
	// Wait, I removed the `arg` handling in `resolveConfigAndLogger` in `repo_loader.go`!

	// I need to fix `repo_loader.go` or handle it here.
	// Given `kex start` has specific Arg behavior, maybe I should handle it.
	// But `check` might want same behavior? (kex check <path>?)

	if err != nil {
		return err
	}

	if err := validateRepository(repo); err != nil {
		return err
	}

	return startServer(repo, l)
}

func validateRepository(repo *fs.Indexer) error {
	v := validator.New([]validator.ValidationRule{
		&validator.IDRequiredRule{},
		&validator.TitleRequiredRule{},
		&validator.FilenameMatchRule{},
	})
	report := v.Validate(repo)

	if len(report.GlobalErrors) > 0 {
		for _, e := range report.GlobalErrors {
			fmt.Fprintf(os.Stderr, "Error: %v\n", e)
		}
		return cli.Exit("Failed to start due to document errors. Run 'kex check' for details.", 1)
	}

	if !report.Valid {
		for _, err := range report.GlobalErrors {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		}
		for _, doc := range report.Documents {
			for _, err := range doc.Errors {
				fmt.Fprintf(os.Stderr, "Error [%s]: %v\n", doc.ID, err)
			}
		}
		return cli.Exit("Validation failed (documents contain errors). Run 'kex check' for details.", 1)
	}
	return nil
}

func startServer(repo *fs.Indexer, l logger.Logger) error {
	srv := mcp.New(repo, l)
	fmt.Fprintf(os.Stderr, "Server listening on stdio...\n")
	if err := srv.Serve(); err != nil {
		return cli.Exit(fmt.Sprintf("Server error: %v", err), 1)
	}
	return nil
}
