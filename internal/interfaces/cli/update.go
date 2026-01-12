package cli

import (
	"os"

	"github.com/mew-ton/kex/assets"
	"github.com/mew-ton/kex/internal/infrastructure/config"
	"github.com/mew-ton/kex/internal/usecase/generator"

	"github.com/pterm/pterm"
	"github.com/urfave/cli/v2"
)

var UpdateCommand = &cli.Command{
	Name:   "update",
	Usage:  "Update kex documentation and configuration",
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

	pterm.Info.Printf("Updating Kex resources in: %s\n", cwd)

	gen := generator.New(assets.Assets)

	spinner, _ := pterm.DefaultSpinner.Start("Updating files...")

	// Pass strategies from config (map[string]string)
	// Strategies are populated by config.Load defaults
	updateConfig := cfg.Update

	// Update uses the first source as the primary location for system docs (e.g. kex/*)
	// This might need refinement for multiplexing, but "kex update" is mainly for
	// the primary documentation repo.
	sourceRoot := "contents"
	if cfg.Source != "" {
		sourceRoot = cfg.Source
	}

	if err := gen.Update(cwd, sourceRoot, updateConfig); err != nil {
		spinner.Fail(err.Error())
		return err
	}
	spinner.Success("Update complete")

	return nil
}
