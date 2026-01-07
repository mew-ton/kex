package cli

import (
	"fmt"
	"os"

	"github.com/mew-ton/kex/assets"
	"github.com/mew-ton/kex/internal/infrastructure/config"
	"github.com/mew-ton/kex/internal/usecase/generator"

	"github.com/pterm/pterm"
	"github.com/urfave/cli/v2"
)

var UpdateCommand = &cli.Command{
	Name:  "update",
	Usage: "Update kex documentation and configuration",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "agent-type",
			Usage:   "Agent type for guidelines (general, claude)",
			Value:   "general",
			Aliases: []string{"a"},
		},
	},
	Action: runUpdate,
}

func runUpdate(c *cli.Context) error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	// Load Config to get update strategies
	cfg, err := config.Load(cwd)
	if err != nil {
		// If config fails, we might still proceed with defaults, but let's warn
		pterm.Warning.Printf("Failed to load .kex.yaml: %v. Using default update strategies.\n", err)
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

	pterm.Info.Printf("Updating Kex resources in: %s\n", cwd)

	gen := generator.New(assets.Templates)

	spinner, _ := pterm.DefaultSpinner.Start("Updating files...")

	// Pass strategies from config
	strategies := cfg.Update.Strategies
	if strategies == nil {
		strategies = make(map[string]string)
	}

	if err := gen.Update(cwd, agentType, strategies); err != nil {
		spinner.Fail(err.Error())
		return err
	}
	spinner.Success("Update complete")

	return nil
}
