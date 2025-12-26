package commands

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"kex/assets"

	"github.com/urfave/cli/v2"
)

var InitCommand = &cli.Command{
	Name:   "init",
	Usage:  "Initialize kex repository",
	Action: runInit,
}

func runInit(c *cli.Context) error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	fmt.Printf("Initializing kex in %s\n", cwd)

	// 1. Create directory structure
	contentsDir := filepath.Join(cwd, "contents")
	if err := os.MkdirAll(contentsDir, 0755); err != nil {
		return fmt.Errorf("failed to create contents directory: %w", err)
	}

	// 2. Extract templates
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
			fmt.Printf("Skipping existing file: %s\n", relPath)
			return nil
		}

		fmt.Printf("Creating %s\n", relPath)
		return os.WriteFile(targetPath, data, 0644)
	})

	if err != nil {
		return fmt.Errorf("failed to extract templates: %w", err)
	}

	// 3. Create .kex.yaml (Simple default)
	configPath := filepath.Join(cwd, ".kex.yaml")
	defaultConfig := []byte(`root: contents
`)
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		if err := os.WriteFile(configPath, defaultConfig, 0644); err != nil {
			return fmt.Errorf("failed to create config: %w", err)
		}
		fmt.Println("Created .kex.yaml")
	}

	fmt.Println("Initialization complete.")
	return nil
}
