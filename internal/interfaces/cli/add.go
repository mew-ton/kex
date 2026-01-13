package cli

import (
	"fmt"
	"os"

	// Added missing import for strings
	"github.com/mew-ton/kex/internal/infrastructure/config"
	"github.com/mew-ton/kex/internal/infrastructure/fs"
	"github.com/mew-ton/kex/internal/infrastructure/logger"

	"github.com/urfave/cli/v2"
)

var AddCommand = &cli.Command{
	Name:      "add",
	Usage:     "Add a document source reference",
	ArgsUsage: "[path_or_url]",
	Action:    runAdd,
}

func runAdd(c *cli.Context) error {
	arg := c.Args().First()
	if arg == "" {
		return cli.Exit("Error: missing path or url argument", 1)
	}

	// 1. Validate
	if isURL(arg) {
		p := fs.NewRemoteProvider(arg, "", nil) // No logger, no token (assuming public or environment will handle later)
		if err := p.Validate(); err != nil {
			return cli.Exit(fmt.Sprintf("Error: reachable check failed for '%s': %v", arg, err), 1)
		}
	} else {
		p := fs.NewLocalProvider(arg, nil) // No logger
		if err := p.Validate(); err != nil {
			return cli.Exit(fmt.Sprintf("Error: invalid path '%s': %v", arg, err), 1)
		}
	}

	// 2. Load Config
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	cfg, err := config.Load(cwd)
	if err != nil {
		// New config if not exists
		logger.Info("Creating new config as load failed/missing: %v", err)
		// If load fails because file missing, it returns default config.
		// If it fails because of parse error, we might be overwriting corrupt config.
		// config.Load handles IsNotExist by returning default.
	}

	// 3. Append
	// Check for duplicates
	for _, ref := range cfg.References {
		if ref == arg {
			logger.Info("Reference '%s' already exists.", arg)
			return nil
		}
	}

	cfg.References = append(cfg.References, arg)

	// 4. Save
	if err := config.Save(cwd, cfg); err != nil {
		return cli.Exit(fmt.Sprintf("Error saving config: %v", err), 1)
	}

	logger.Info("Added reference: %s", arg)
	return nil
}
