package commands

import (
	"fmt"
	"os"

	"kex/internal/indexer"
	"kex/internal/server"

	"github.com/urfave/cli/v2"
)

var StartCommand = &cli.Command{
	Name:   "start",
	Usage:  "Start MCP server",
	Action: runStart,
}

func runStart(c *cli.Context) error {
	// 1. Resolve configuration
	root := "contents"
	if _, err := os.Stat(root); os.IsNotExist(err) {
		return cli.Exit(fmt.Sprintf("Error: directory '%s' not found. Run 'kex init'?", root), 1)
	}

	// 2. Load Indexer
	idx := indexer.New(root)
	if err := idx.Load(); err != nil {
		return cli.Exit(fmt.Sprintf("Fatal: failed to load documents: %v", err), 1)
	}

	// 3. Strict validation on start (as per design)
	if len(idx.Errors) > 0 {
		for _, e := range idx.Errors {
			fmt.Fprintf(os.Stderr, "Error: %v\n", e)
		}
		return cli.Exit("Failed to start due to document errors. Run 'kex check' for details.", 1)
	}

	// 4. Start Server
	srv := server.New(idx)
	if err := srv.Serve(); err != nil {
		return cli.Exit(fmt.Sprintf("Server error: %v", err), 1)
	}

	return nil
}
