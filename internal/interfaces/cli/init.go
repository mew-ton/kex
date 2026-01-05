package cli

import (
	"github.com/mew-ton/kex/assets"
	"github.com/mew-ton/kex/internal/usecase/generator"
	"os"

	"github.com/pterm/pterm"
	"github.com/urfave/cli/v2"
)

var InitCommand = &cli.Command{
	Name:   "init",
	Usage:  "Initialize kex repository",
	Action: runInit,
}

func runInit(c *cli.Context) error {
	// Welcome Message
	pterm.DefaultBigText.WithLetters(
		pterm.NewLettersFromString("KEX"),
	).Render()

	pterm.DefaultHeader.Println("Initializing Kex Repository")

	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	pterm.Info.Printf("Initializing in: %s\n", cwd)

	gen := generator.New(assets.Templates)

	spinner, _ := pterm.DefaultSpinner.Start("Generating project structure...")
	if err := gen.Generate(cwd); err != nil {
		spinner.Fail(err.Error())
		return err
	}
	spinner.Success("Project structure generated")

	pterm.Println() // Spacer
	pterm.DefaultSection.Println("Initialization complete!")
	pterm.Info.Println("Run 'kex check' to validate your documents.")

	return nil
}
