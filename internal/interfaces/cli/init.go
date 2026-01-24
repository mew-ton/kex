package cli

import (
	"os"
	"sort"
	"strings"

	"github.com/mew-ton/kex/assets"
	"github.com/mew-ton/kex/internal/infrastructure/config"
	"github.com/mew-ton/kex/internal/infrastructure/ui"
	"github.com/mew-ton/kex/internal/usecase/generator"

	"github.com/pterm/pterm"
	"github.com/urfave/cli/v2"
)

var InitCommand = &cli.Command{
	Name:   "init",
	Usage:  "Initialize kex repository",
	Action: runInit,
	Flags: []cli.Flag{
		&cli.StringSliceFlag{
			Name:    "agents",
			Aliases: []string{"a"},
			Usage:   "AI Agents to enable (Antigravity, Cursor, Claude)",
		},
		&cli.StringSliceFlag{
			Name:    "scopes",
			Aliases: []string{"s"},
			Usage:   "Scopes to enable (coding, documentation)",
		},
		&cli.StringSliceFlag{
			Name:  "skills",
			Usage: "Keywords for Skills (e.g. go, typescript)",
		},
	},
}

type initSelection struct {
	mcpAgents      map[string]bool
	mcpScopes      []string
	skillsAgents   map[string]bool
	skillsKeywords []string
}

func runInit(c *cli.Context) error {
	printWelcome()

	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	selection, err := resolveSelection(c)
	if err != nil {
		return err
	}

	pterm.Info.Printf("Initializing in: %s\n", cwd)

	if err := saveConfig(cwd, selection); err != nil {
		return err
	}

	return runUpdate(c)
}

func resolveSelection(c *cli.Context) (*initSelection, error) {
	if c.IsSet("agents") || c.IsSet("scopes") || c.IsSet("skills") {
		return resolveNonInteractiveSelection(c)
	}
	return resolveInteractiveSelection()
}

func resolveNonInteractiveSelection(c *cli.Context) (*initSelection, error) {
	sel := &initSelection{
		mcpAgents: make(map[string]bool),
	}

	for _, a := range c.StringSlice("agents") {
		sel.mcpAgents[strings.ToLower(a)] = true
	}

	sel.mcpScopes = c.StringSlice("scopes")

	if c.IsSet("skills") {
		sel.skillsAgents = map[string]bool{"claude": true}
		sel.skillsKeywords = c.StringSlice("skills")
	}

	return sel, nil
}

func resolveInteractiveSelection() (*initSelection, error) {
	capabilities, err := selectCapabilities()
	if err != nil {
		return nil, err
	}

	sel := &initSelection{
		mcpAgents:    make(map[string]bool),
		skillsAgents: make(map[string]bool),
	}

	hasMcpCapability := false
	hasSkillsCapability := false

	for _, cap := range capabilities {
		if strings.Contains(cap, "(MCP Rules)") {
			hasMcpCapability = true
			agentName := strings.Split(cap, " ")[0]
			sel.mcpAgents[strings.ToLower(agentName)] = true
		}
		if strings.Contains(cap, "(Skills)") {
			hasSkillsCapability = true
			agentName := strings.Split(cap, " ")[0]
			sel.skillsAgents[strings.ToLower(agentName)] = true
		}
	}

	if hasMcpCapability {
		sel.mcpScopes = append(sel.mcpScopes, "coding")
		enableDocs, err := confirmDocumentationSupport()
		if err != nil {
			return nil, err
		}
		if enableDocs {
			sel.mcpScopes = append(sel.mcpScopes, "documentation")
		}
	}

	if hasSkillsCapability {
		keywords, err := inputSkillsKeywords()
		if err != nil {
			return nil, err
		}
		sel.skillsKeywords = keywords
	}

	return sel, nil
}

func printWelcome() {
	pterm.DefaultBigText.WithLetters(
		pterm.NewLettersFromString("KEX"),
	).Render()
	pterm.DefaultHeader.Println("Initializing Kex Repository")
}

func selectCapabilities() ([]string, error) {
	manifest, err := generator.LoadManifest(assets.Assets)
	if err != nil {
		return nil, err
	}

	var options []string
	if _, ok := manifest.AiAgents["antigravity"]; ok {
		options = append(options, "Antigravity (MCP Rules)")
		options = append(options, "Antigravity (Skills)")
	}
	if _, ok := manifest.AiAgents["claude"]; ok {
		options = append(options, "Claude (MCP Rules)")
		options = append(options, "Claude (Skills)")
	}
	if _, ok := manifest.AiAgents["cursor"]; ok {
		options = append(options, "Cursor (MCP Rules)")
	}

	pterm.Info.Println("MCP Rules are static guidelines enforced by the AI. Skills are dynamic knowledge retrieved by keywords.")

	preSelected := map[string]bool{
		"Antigravity (MCP Rules)": true,
		"Claude (MCP Rules)":      true,
	}

	return ui.Multiselect("Select Agent Capabilities to enable", options, preSelected)
}

func confirmDocumentationSupport() (bool, error) {
	pterm.Println()
	pterm.Info.Println("Indexable documents allow the AI to answer questions about this repository's business logic, architecture, etc.")

	return pterm.DefaultInteractiveConfirm.
		WithDefaultText("Do you want to provide indexable documents in this repository?").
		WithDefaultValue(false).
		Show()
}

func inputSkillsKeywords() ([]string, error) {
	prompt := pterm.DefaultInteractiveTextInput
	prompt.DefaultText = "go, typescript"
	prompt.Delimiter = ": "

	pterm.Println()
	pterm.Print(pterm.Cyan("? ") + "Enter keywords for Skills (comma separated) (Examples: coding, documentation, kex)\n")

	result, err := prompt.Show("> ")
	if err != nil {
		return nil, err
	}

	if strings.TrimSpace(result) == "" {
		return []string{"coding", "documentation", "kex"}, nil
	}

	parts := strings.Split(result, ",")
	var keywords []string
	for _, p := range parts {
		clean := strings.TrimSpace(p)
		if clean != "" {
			keywords = append(keywords, clean)
		}
	}
	return keywords, nil
}

func saveConfig(cwd string, sel *initSelection) error {
	var mcpTargets []string
	for agent := range sel.mcpAgents {
		mcpTargets = append(mcpTargets, agent)
	}
	sort.Strings(mcpTargets)

	var skillsTargets []string
	for agent := range sel.skillsAgents {
		skillsTargets = append(skillsTargets, agent)
	}
	sort.Strings(skillsTargets)

	docs := buildDocumentsConfig(sel.mcpScopes)

	cfg := config.Config{
		Source: "contents",
		Update: config.UpdateConfig{
			Documents: docs,
			AiMcpRules: config.AiMcpRules{
				Targets: mcpTargets,
				Scopes:  sel.mcpScopes,
			},
			AiSkills: config.AiSkills{
				Targets:  skillsTargets,
				Keywords: sel.skillsKeywords,
			},
		},
	}

	if err := config.Save(cwd, cfg); err != nil {
		pterm.Warning.Printf("Failed to save .kex.yaml: %v\n", err)
		return err
	}
	return nil
}

func buildDocumentsConfig(scopes []string) map[string]string {
	docs := make(map[string]string)
	for _, s := range scopes {
		if s == "documentation" {
			docs["kex"] = "all"
			break
		}
	}
	return docs
}
