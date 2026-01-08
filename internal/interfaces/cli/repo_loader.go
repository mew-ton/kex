package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/mew-ton/kex/internal/infrastructure/config"
	"github.com/mew-ton/kex/internal/infrastructure/fs"
	"github.com/mew-ton/kex/internal/infrastructure/logger"
	"github.com/urfave/cli/v2"
)

// Helper to resolve Config, Logger, and Repository from CLI context
func resolveRepository(c *cli.Context) (*fs.Indexer, config.Config, logger.Logger, error) {
	// 1. Determine Main Project Root (for global config)
	// Priority: First arg that is a Project Root, or "."
	mainProjectRoot := "."
	mainConfigLoaded := false

	var contentRoots []string // Paths to actual content (sources)

	// Check positional args
	if c.Args().Present() {
		for _, arg := range c.Args().Slice() {
			// Check if arg is a directory and has .kex.yaml
			configPath := filepath.Join(arg, ".kex.yaml")
			if info, err := os.Stat(configPath); err == nil && !info.IsDir() {
				// FOUND PROJECT ROOT
				// 1. If this is the first project root found, set as Main
				if !mainConfigLoaded {
					mainProjectRoot = arg
					mainConfigLoaded = true
				}

				// 2. Load its config to find its content source
				subCfg, err := config.Load(arg)
				if err == nil && subCfg.Source != "" {
					resolvedSource := filepath.Join(arg, subCfg.Source)
					contentRoots = append(contentRoots, resolvedSource)
				} else {
					// Fallback if load fails? Treat arg as content root?
					// If yaml exists but invalid, ignoring it is safer than indexing raw root?
					// Let's fallback to "contents" default relative to arg.
					contentRoots = append(contentRoots, filepath.Join(arg, "contents"))
				}
			} else {
				// No .kex.yaml, treat as CONTENT ROOT directly
				contentRoots = append(contentRoots, arg)
			}
		}
	}

	// 2. Resolve Main Config and Logger
	cfg, appLogger, err := resolveConfigAndLogger(c, mainProjectRoot)
	if err != nil {
		// Log warning but proceed?
	}

	// 3. Resolve Remote Token (Env Override)
	if token := os.Getenv("KEX_REMOTE_TOKEN"); token != "" {
		cfg.RemoteToken = token
	}

	// 4. Create Repository
	// Note: Factory expects `cliSources` which override config.
	// But in Multi-Root mode, `contentRoots` IS the list of all sources.
	// If `contentRoots` is empty, Factory uses `cfg.Source`.
	// If `contentRoots` is populated (from args), Factory uses them.
	// If user ran `kex start` (no args), `contentRoots` is empty.
	// Factory uses `cfg.Source` (from mainProjectRoot="."). Correct.
	// If user ran `kex start .` -> `.` is Project Root. `contentRoots` gets `[./contents]`.
	// Factory uses `[./contents]`. Correct.
	factory := fs.NewFactory(appLogger)
	repo, err := factory.CreateRepository(cfg, contentRoots, mainProjectRoot)
	if err != nil {
		return nil, cfg, appLogger, err
	}

	// 5. Load
	if err := repo.Load(); err != nil {
		return nil, cfg, appLogger, fmt.Errorf("failed to load documents: %w", err)
	}

	return repo, cfg, appLogger, nil
}

// resolveConfigAndLogger loads config from projectRoot and sets up logger
func resolveConfigAndLogger(c *cli.Context, projectRoot string) (config.Config, logger.Logger, error) {
	cfg, err := config.Load(projectRoot)

	// Logger
	var appLogger logger.Logger
	logFile := c.String("log-file")
	if logFile == "" {
		logFile = cfg.Logging.File
	}

	if logFile != "" {
		fl, err := logger.NewFileLogger(logFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to open log file '%s': %v. Using stderr.\n", logFile, err)
			appLogger = logger.NewStderrLogger()
		} else {
			appLogger = fl
		}
	} else {
		appLogger = logger.NewStderrLogger()
	}

	return cfg, appLogger, err
}
