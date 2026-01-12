package cli

import (
	"os"
	"sort"
	"strings"

	"github.com/mew-ton/kex/assets"
	"github.com/mew-ton/kex/internal/infrastructure/config"
	"github.com/mew-ton/kex/internal/usecase/generator"
	"gopkg.in/yaml.v3"

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
		// Add skills flag for completeness?
		&cli.StringSliceFlag{
			Name:  "skills",
			Usage: "Keywords for Skills (e.g. go, typescript)",
		},
	},
}

func runInit(c *cli.Context) error {
	printWelcome()

	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	var selectedMcpAgents map[string]bool
	var selectedMcpScopes []string
	var selectedSkillsAgents map[string]bool
	var skillsKeywords []string

	// Check if flags are used for non-interactive mode
	if c.IsSet("agents") || c.IsSet("scopes") || c.IsSet("skills") {
		// Non-interactive mode
		selectedMcpAgents = make(map[string]bool)
		for _, a := range c.StringSlice("agents") {
			selectedMcpAgents[strings.ToLower(a)] = true
		}

		selectedMcpScopes = c.StringSlice("scopes")

		// If custom skills are provided, we enable skills for Claude by default?
		// Or assume agents flag covers it?
		// For backward compat (and simplicity), let's say if "claude" is in agents, we enable MCP rules.
		// If "skills" are provided, we enable skills for Claude (since currently only Claude supports it).
		if c.IsSet("skills") {
			selectedSkillsAgents = map[string]bool{"claude": true}
			skillsKeywords = c.StringSlice("skills")
		}
	} else {
		// Interactive Mode
		// 1. Select Capabilities
		capabilities, err := selectCapabilities()
		if err != nil {
			return err
		}

		selectedMcpAgents = make(map[string]bool)
		selectedSkillsAgents = make(map[string]bool)

		hasMcpCapability := false
		hasSkillsCapability := false

		for _, cap := range capabilities {
			if strings.Contains(cap, "(MCP Rules)") {
				hasMcpCapability = true
				agentName := strings.Split(cap, " ")[0]
				selectedMcpAgents[strings.ToLower(agentName)] = true
			}
			if strings.Contains(cap, "(Skills)") {
				hasSkillsCapability = true
				agentName := strings.Split(cap, " ")[0]
				selectedSkillsAgents[strings.ToLower(agentName)] = true
			}
		}

		// 2. Select Scopes for MCP Rules
		if hasMcpCapability {
			scopes, err := selectMcpScopes()
			if err != nil {
				return err
			}
			selectedMcpScopes = scopes
		}

		// 3. Input Skills Keywords
		if hasSkillsCapability {
			keywords, err := inputSkillsKeywords()
			if err != nil {
				return err
			}
			skillsKeywords = keywords
		}
	}

	pterm.Info.Printf("Initializing in: %s\n", cwd)

	if err := saveConfig(selectedMcpAgents, selectedMcpScopes, selectedSkillsAgents, skillsKeywords); err != nil {
		return err
	}

	// Run Update logic
	return runUpdate(c)
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

	// Dynamically build capabilities list from manifest
	// For v2, we hardcode known logic mapping for now based on the feedback
	// "Select Agent Capabilities to enable"
	// [x] Antigravity (MCP Rules)
	// [x] Claude (MCP Rules)
	// [ ] Claude (Skills)
	// [ ] Cursor (MCP Rules)

	var options []string
	// Sort agent keys to have consistent order?
	var agentKeys []string
	for k := range manifest.AiAgents {
		agentKeys = append(agentKeys, k)
	}
	sort.Strings(agentKeys)

	// Custom mapping logic for display
	// Verify agent exists in manifest before adding option
	if _, ok := manifest.AiAgents["antigravity"]; ok {
		options = append(options, "Antigravity (MCP Rules)")
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

	return Multiselect("Select Agent Capabilities to enable", options, preSelected)
}

func selectMcpScopes() ([]string, error) {
	options := []string{"coding", "documentation"}
	preSelected := map[string]bool{
		"coding":        true,
		"documentation": true,
	}
	return Multiselect("Select Scopes for MCP Rules (What logic should be enforced via MCP?)", options, preSelected)
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
		// Defaults
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

func saveConfig(mcpAgents map[string]bool, mcpScopes []string, skillsAgents map[string]bool, skillsKeywords []string) error {
	// Build AiMcpRules
	var mcpTargets []string
	for agent := range mcpAgents {
		mcpTargets = append(mcpTargets, agent)
	}
	sort.Strings(mcpTargets)

	// Build AiSkills
	var skillsTargets []string
	for agent := range skillsAgents {
		skillsTargets = append(skillsTargets, agent)
	}
	sort.Strings(skillsTargets)

	// Documents defaults
	docs := map[string]string{
		"kex": "all",
	}

	// Create Config
	cfg := config.Config{
		Sources: []string{"contents"},
		Update: config.UpdateConfig{
			Documents: docs,
			AiMcpRules: config.AiMcpRules{
				Targets: mcpTargets,
				Scopes:  mcpScopes,
			},
			AiSkills: config.AiSkills{
				Targets:  skillsTargets,
				Keywords: skillsKeywords,
			},
		},
	}

	data, err := yaml.Marshal(&cfg)
	if err != nil {
		return err
	}

	if err := os.WriteFile(".kex.yaml", data, 0644); err != nil {
		pterm.Warning.Printf("Failed to save .kex.yaml: %v\n", err)
		return err
	}
	return nil
}
