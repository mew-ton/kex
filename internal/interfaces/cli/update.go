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

	var agentConfig *config.Agent

	// Determine Agent Type
	// Priority: Flag > Config > Default
	var agentType generator.AgentType

	if cfg.Agent.Type != "" {
		agentType = generator.AgentType(cfg.Agent.Type)
		agentConfig = &cfg.Agent
	} else {
		agentType = generator.AgentTypeGeneral
	}

	if c.IsSet("agent-type") {
		// Override if flag is explicitly provided
		val := c.String("agent-type")
		if val == string(generator.AgentTypeGeneral) || val == string(generator.AgentTypeClaude) {
			agentType = generator.AgentType(val)
		} else {
			return fmt.Errorf("invalid agent type: %s", val)
		}
	}

	pterm.Info.Printf("Updating Kex resources in: %s (Agent: %s)\n", cwd, agentType)

	gen := generator.New(assets.Templates)

	spinner, _ := pterm.DefaultSpinner.Start("Updating files...")

	// Pass strategies from config
	strategies := cfg.Update.Strategies
	if strategies == nil {
		strategies = make(map[string]string)
	}

	if err := gen.Update(cwd, cfg.Root, agentType, strategies, agentConfig); err != nil {
		spinner.Fail(err.Error())
		return err
	}
	spinner.Success("Update complete")

	return nil
}
