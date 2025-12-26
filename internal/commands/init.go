package commands

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"kex/assets"

	"github.com/pterm/pterm"
	"github.com/urfave/cli/v2"
)

var InitCommand = &cli.Command{
	Name:   "init",
	Usage:  "Initialize kex repository",
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

	pterm.Info.Printf("Initializing in: %s\n", cwd)

	// 1. Create directory structure
	spinner, _ := pterm.DefaultSpinner.Start("Creating directory structure...")
	contentsDir := filepath.Join(cwd, "contents")
	if err := os.MkdirAll(contentsDir, 0755); err != nil {
		spinner.Fail("Failed to create contents directory")
		return fmt.Errorf("failed to create contents directory: %w", err)
	}
	spinner.Success("Directory structure created")

	// 2. Extract templates
	spinner, _ = pterm.DefaultSpinner.Start("Extracting templates...")
	err = fs.WalkDir(assets.Templates, "templates", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel("templates", path)
		if err != nil {
			return err
		}

		targetPath := filepath.Join(contentsDir, relPath)

		if d.IsDir() {
			return os.MkdirAll(targetPath, 0755)
		}

		data, err := fs.ReadFile(assets.Templates, path)
		if err != nil {
			return err
		}

		// Don't overwrite if exists
		if _, err := os.Stat(targetPath); err == nil {
			// Just log info, don't fail spinner
			pterm.Info.Printf("Skipping existing file: %s\n", relPath)
			return nil
		}

		// Using Printf inside spinner might be messy, but for detailed verbose we can.
		// For now, let's keep it clean.
		return os.WriteFile(targetPath, data, 0644)
	})

	if err != nil {
		spinner.Fail("Failed to extract templates")
		return fmt.Errorf("failed to extract templates: %w", err)
	}
	spinner.Success("Templates extracted")

	// 3. Create .kex.yaml (Simple default)
	configPath := filepath.Join(cwd, ".kex.yaml")
	defaultConfig := []byte(`root: contents
`)
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		if err := os.WriteFile(configPath, defaultConfig, 0644); err != nil {
			return fmt.Errorf("failed to create config: %w", err)
		}
		pterm.Success.Println("Created .kex.yaml")
	}

	pterm.Println() // Spacer
	pterm.DefaultSection.Println("Initialization complete!")
	pterm.Info.Println("Run 'kex check' to validate your documents.")

	return nil
}
