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
	var scopes []string

	if c.IsSet("agent-type") {
		// Non-interactive Mode (Flag provided)
		pterm.Info.Println("Agent Type provided via flag. Skipping interactive mode.")
		switch c.String("agent-type") {
		case string(generator.AgentTypeGeneral):
			agentType = generator.AgentTypeGeneral
		case string(generator.AgentTypeClaude):
			agentType = generator.AgentTypeClaude
		default:
			return fmt.Errorf("invalid agent type: %s. Must be 'general' or 'claude'", c.String("agent-type"))
		}
		// Default scopes for non-interactive
		scopes = []string{"coding", "documentation"}
	} else {
		// Interactive Mode
		selectedType, _ := pterm.DefaultInteractiveSelect.
			WithOptions([]string{string(generator.AgentTypeGeneral), string(generator.AgentTypeClaude)}).
			WithDefaultText("Select Agent Type").
			Show()
		agentType = generator.AgentType(selectedType)

		selectedScopes, _ := pterm.DefaultInteractiveMultiselect.
			WithOptions([]string{"coding", "documentation"}).
			WithDefaultText("Select Scopes").
			WithFilter(false).
			WithDefaultOptions([]string{"coding", "documentation"}).
			Show()
		scopes = selectedScopes
	}

	pterm.Info.Printf("Initializing in: %s (Agent: %s, Scopes: %v)\n", cwd, agentType, scopes)

	// Save Config
	// Manual string construction to avoid complexity of a full Save method for now
	// Ensure scopes are formatted as YAML list [ "a", "b" ]
	scopesYaml := "["
	for i, s := range scopes {
		if i > 0 {
			scopesYaml += ", "
		}
		scopesYaml += fmt.Sprintf("\"%s\"", s)
	}
	scopesYaml += "]"

	configData := fmt.Sprintf("root: contents\nagent:\n  type: %s\n  scopes: %s\n", agentType, scopesYaml)
	if err := os.WriteFile(".kex.yaml", []byte(configData), 0644); err != nil {
		pterm.Warning.Printf("Failed to save .kex.yaml: %v\n", err)
	}

	gen := generator.New(assets.Templates)

	spinner, _ := pterm.DefaultSpinner.Start("Generating project structure...")

	// We need to pass agent config to Generate as well?
	// The current Generate signature only takes AgentType.
	// But AGENTS.md generation now depends on scopes.
	// We should update Generator.Generate to take *config.Agent
	agentConfig := &config.Agent{Type: string(agentType), Scopes: scopes}

	// Update Generator.Generate signature in separate step or assume it uses defaults?
	// Wait, we updated Generator.Update but not Generator.Generate.
	// We need to update Generator.Generate too.

	// For this step, let's assume we will update Generate signature next.
	if err := gen.Generate(cwd, agentType, agentConfig); err != nil {
		spinner.Fail(err.Error())
		return err
	}
	spinner.Success("Project structure generated")

	pterm.Println() // Spacer
	pterm.DefaultSection.Println("Initialization complete!")
	pterm.Info.Println("Run 'kex check' to validate your documents.")

	return nil
}
