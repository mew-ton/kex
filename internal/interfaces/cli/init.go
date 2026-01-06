package cli

import (
	"fmt"
	"os"

	"github.com/mew-ton/kex/assets"
	"github.com/mew-ton/kex/internal/usecase/generator"

	"github.com/pterm/pterm"
	"github.com/urfave/cli/v2"
)

var InitCommand = &cli.Command{
	Name:  "init",
	Usage: "Initialize kex repository",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "agent-type",
			Usage:   "Agent type for guidelines (general, claude)",
			Value:   "general",
			Aliases: []string{"a"},
		},
	},
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

	var agentType generator.AgentType
	switch c.String("agent-type") {
	case string(generator.AgentTypeGeneral):
		agentType = generator.AgentTypeGeneral
	case string(generator.AgentTypeClaude):
		agentType = generator.AgentTypeClaude
	default:
		return fmt.Errorf("invalid agent type: %s. Must be 'general' or 'claude'", c.String("agent-type"))
	}

	pterm.Info.Printf("Initializing in: %s (Agent: %s)\n", cwd, agentType)

	gen := generator.New(assets.Templates)

	spinner, _ := pterm.DefaultSpinner.Start("Generating project structure...")
	if err := gen.Generate(cwd, agentType); err != nil {
		spinner.Fail(err.Error())
		return err
	}
	spinner.Success("Project structure generated")

	pterm.Println() // Spacer
	pterm.DefaultSection.Println("Initialization complete!")
	pterm.Info.Println("Run 'kex check' to validate your documents.")

	return nil
}
