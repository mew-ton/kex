package main

import (
	"log"
	"os"

	kexcli "github.com/mew-ton/kex/internal/interfaces/cli"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "kex",
		Usage: "Document Librarian Tool (MCP / Skills Management)",
		Commands: []*cli.Command{
			kexcli.InitCommand,
			kexcli.CheckCommand,
			kexcli.StartCommand,
			kexcli.GenerateCommand,
			kexcli.UpdateCommand,
			kexcli.AddCommand,
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
