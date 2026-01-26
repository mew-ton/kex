package cli

import (
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/mew-ton/kex/assets"
	"github.com/mew-ton/kex/internal/infrastructure/config"
	"github.com/mew-ton/kex/internal/infrastructure/ui"
	"github.com/mew-ton/kex/internal/usecase/generator"
	"github.com/pterm/pterm"
	"github.com/urfave/cli/v2"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
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
	agents         map[string]bool
	keywords       []string
	createContents bool
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

	// Ensure contents directory exists if selected
	if selection.createContents {
		if err := os.MkdirAll(filepath.Join(cwd, "contents"), 0755); err != nil {
			return err
		}
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
		agents:         make(map[string]bool),
		createContents: true, // Default to true for non-interactive for now, or check flags
	}

	for _, a := range c.StringSlice("agents") {
		sel.agents[strings.ToLower(a)] = true
	}

	// Default to Antigravity and Claude if no agents specified
	if len(sel.agents) == 0 {
		sel.agents["antigravity"] = true
		sel.agents["claude"] = true
	}

	if c.IsSet("skills") {
		sel.keywords = c.StringSlice("skills")
	}

	return sel, nil
}

func resolveInteractiveSelection() (*initSelection, error) {
	selectedAgents, err := selectAgents()
	if err != nil {
		return nil, err
	}

	sel := &initSelection{
		agents: make(map[string]bool),
	}

	for _, a := range selectedAgents {
		sel.agents[strings.ToLower(a)] = true
	}

	// If no agents selected, we can skip other prompts or just return empty
	if len(sel.agents) == 0 {
		return sel, nil
	}

	createContents, err := confirmContentsCreation()
	if err != nil {
		return nil, err
	}
	sel.createContents = createContents

	// Only ask for keywords if we are creating/maintaining contents (Skills source)
	if !sel.createContents {
		keywords, err := inputSkillsKeywords()
		if err != nil {
			return nil, err
		}
		sel.keywords = keywords
	}

	return sel, nil
}

func printWelcome() {
	pterm.DefaultBigText.WithLetters(
		pterm.NewLettersFromString("KEX"),
	).Render()
	pterm.DefaultHeader.Println("Initializing Kex Repository")
}

func selectAgents() ([]string, error) {
	manifest, err := generator.LoadManifest(assets.Assets)
	if err != nil {
		return nil, err
	}

	var options []string
	for agentName := range manifest.AiAgents {
		// Capitalize for display
		options = append(options, cases.Title(language.English).String(agentName))
	}
	sort.Strings(options)

	pterm.Println()
	pterm.Info.Println("Select the AI Agents you want to target (Guidelines & Skills).")

	// Default to Antigravity if available, others optional
	preSelected := map[string]bool{
		"Antigravity": true,
		"Claude":      true,
	}

	return ui.Multiselect("Target Agents", options, preSelected)
}

func confirmContentsCreation() (bool, error) {
	pterm.Println() // Add newline for better spacing
	// Requirement 7: Ask "Do you want to maintain common guidelines..."
	return pterm.DefaultInteractiveConfirm.
		WithDefaultText("Do you want to maintain common guidelines (creates contents/ directory)?").
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
	var targets []string
	for agent := range sel.agents {
		targets = append(targets, agent)
	}
	sort.Strings(targets)

	// Documents config: If creating contents, enable 'kex' docs.
	// If not, keep empty or minimal.
	docs := make(map[string]string)
	source := ""
	if sel.createContents {
		docs["kex"] = "all"
		source = "contents"
	}

	cfg := config.Config{
		Source: source,
		Update: config.UpdateConfig{
			Documents: docs,
			Ai: config.AiConfig{
				Targets:  targets,
				Keywords: sel.keywords,
			},
		},
	}

	if err := config.Save(cwd, cfg); err != nil {
		pterm.Warning.Printf("Failed to save .kex.yaml: %v\n", err)
		return err
	}
	return nil
}

// Deprecated or Unused helper functions removed
