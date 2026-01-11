package cli

import (
	"fmt"
	"os"
	"sort"

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
	},
}

func runInit(c *cli.Context) error {
	printWelcome()

	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	selectedAgents, err := selectAgents(c)
	if err != nil {
		return err
	}

	selectedScopes, err := selectScopes(c)
	if err != nil {
		return err
	}

	pterm.Info.Printf("Initializing in: %s (Agents: %v, Scopes: %v)\n", cwd, selectedAgents, selectedScopes)

	if err := saveConfig(selectedAgents, selectedScopes); err != nil {
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

func selectAgents(c *cli.Context) ([]string, error) {
	manifest, err := generator.LoadManifest(assets.Assets)
	if err != nil {
		return nil, err
	}

	// Build options map
	nameToKey := make(map[string]string) // Name -> Key
	keyToKey := make(map[string]string)  // Key -> Key (validation)
	var options []string

	for key, def := range manifest.AiAgents {
		options = append(options, def.Name)
		nameToKey[def.Name] = key
		keyToKey[key] = key
	}
	sort.Strings(options)

	var inputs []string
	if c.IsSet("agents") {
		inputs = c.StringSlice("agents")
	} else {
		preSelected := map[string]bool{
			"Antigravity": true,
			"Cursor":      true,
		}

		resultNames, err := Multiselect("Select AI Agents", options, preSelected)
		if err != nil {
			return nil, err
		}
		inputs = resultNames
	}

	if len(inputs) == 0 {
		return nil, fmt.Errorf("at least one agent must be selected")
	}

	var validKeys []string
	for _, input := range inputs {
		// 1. Check if it matches a Key
		if _, ok := keyToKey[input]; ok {
			validKeys = append(validKeys, input)
			continue
		}
		// 2. Check if it matches a Name
		if key, ok := nameToKey[input]; ok {
			validKeys = append(validKeys, key)
			continue
		}
		// 3. Fallback: Case-insensitive check against Name or Key?
		// For robustness, let's assume if user types "antigravity" (key) or "Antigravity" (name), it works.
		// Previous checks covered exact matches.
		// Let's add simple normalization failure.
		return nil, fmt.Errorf("unknown agent: %s", input)
	}

	return validKeys, nil
}

func selectScopes(c *cli.Context) ([]string, error) {
	if c.IsSet("scopes") {
		return c.StringSlice("scopes"), nil
	}

	options := []string{"coding", "documentation"}
	preSelected := map[string]bool{
		"coding":        true,
		"documentation": true,
	}

	return Multiselect("Select Scopes", options, preSelected)
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func saveConfig(agents []string, scopes []string) error {
	// Create Config with strategies map
	strategies := make(map[string]string)

	// System Docs: Always Enable
	strategies["kex"] = "all"

	// Agents
	hasCoding := contains(scopes, "coding")
	hasDoc := contains(scopes, "documentation")

	mode := "none"
	if hasCoding && hasDoc {
		mode = "all"
	} else if hasCoding {
		mode = "coding-only"
	} else if hasDoc {
		mode = "documentation-only"
	}

	if mode != "none" {
		for _, agentKey := range agents {
			strategies[agentKey] = mode
		}
	}

	cfg := config.Config{
		Root: "contents",
		Update: config.Update{
			Strategies: strategies,
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
