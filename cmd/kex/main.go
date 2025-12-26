package main

import (
	"log"
	"os"

	"kex/internal/commands"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "kex",
		Usage: "Document Librarian MCP",
		Commands: []*cli.Command{
			commands.InitCommand,
			commands.CheckCommand,
			commands.StartCommand,
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
