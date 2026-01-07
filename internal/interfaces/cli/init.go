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
	printWelcome()

	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	agentType, scopes, err := determineAgentConfig(c)
	if err != nil {
		return err
	}

	pterm.Info.Printf("Initializing in: %s (Agent: %s, Scopes: %v)\n", cwd, agentType, scopes)

	saveConfig(agentType, scopes)

	return generateProject(cwd, agentType, scopes)
}

func printWelcome() {
	pterm.DefaultBigText.WithLetters(
		pterm.NewLettersFromString("KEX"),
	).Render()
	pterm.DefaultHeader.Println("Initializing Kex Repository")
}

func determineAgentConfig(c *cli.Context) (generator.AgentType, []string, error) {
	if c.IsSet("agent-type") {
		// Non-interactive Mode
		pterm.Info.Println("Agent Type provided via flag. Skipping interactive mode.")
		switch c.String("agent-type") {
		case string(generator.AgentTypeGeneral):
			return generator.AgentTypeGeneral, []string{"coding", "documentation"}, nil
		case string(generator.AgentTypeClaude):
			return generator.AgentTypeClaude, []string{"coding", "documentation"}, nil
		default:
			return "", nil, fmt.Errorf("invalid agent type: %s. Must be 'general' or 'claude'", c.String("agent-type"))
		}
	}

	// Interactive Mode
	selectedType, _ := pterm.DefaultInteractiveSelect.
		WithOptions([]string{string(generator.AgentTypeGeneral), string(generator.AgentTypeClaude)}).
		WithDefaultText("Select Agent Type").
		Show()

	selectedScopes, _ := pterm.DefaultInteractiveMultiselect.
		WithOptions([]string{"coding", "documentation"}).
		WithDefaultText("Select Scopes").
		WithFilter(false).
		WithDefaultOptions([]string{"coding", "documentation"}).
		Show()

	return generator.AgentType(selectedType), selectedScopes, nil
}

func saveConfig(agentType generator.AgentType, scopes []string) {
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
}

func generateProject(cwd string, agentType generator.AgentType, scopes []string) error {
	gen := generator.New(assets.Templates)
	spinner, _ := pterm.DefaultSpinner.Start("Generating project structure...")

	agentConfig := &config.Agent{Type: string(agentType), Scopes: scopes}

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
