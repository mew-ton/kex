package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/mew-ton/kex/internal/infrastructure/config"
	"github.com/mew-ton/kex/internal/infrastructure/fs"
	"github.com/mew-ton/kex/internal/infrastructure/logger"
	"github.com/mew-ton/kex/internal/interfaces/mcp"
	"github.com/mew-ton/kex/internal/usecase/retrieve"
	"github.com/mew-ton/kex/internal/usecase/search"
	"github.com/mew-ton/kex/internal/usecase/validator"

	"github.com/urfave/cli/v2"
)

var StartCommand = &cli.Command{
	Name:  "start",
	Usage: "Start MCP server",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "log-file",
			Usage: "Path to log file",
		},
		&cli.StringFlag{
			Name:  "cwd",
			Usage: "Working directory for the operation",
		},
	},
	Action: runStart,
}

func runStart(c *cli.Context) error {
	fmt.Fprintf(os.Stderr, "Starting Kex Server...\n")

	// 1. Setup Logger
	cwd := c.String("cwd")
	var err error
	if cwd == "" {
		cwd, err = os.Getwd()
		if err != nil {
			return err
		}
	} else {
		cwd, err = filepath.Abs(cwd)
		if err != nil {
			return fmt.Errorf("invalid cwd path: %w", err)
		}
	}

	// Load config for logging & sources
	cfg, _ := config.Load(cwd)

	// Setup Logger
	appLogger, err := resolveLogger(c, cfg, cwd)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: logger setup failed: %v. Using stderr.\n", err)
		appLogger = logger.NewStderrLogger()
	}
	logger.SetGeneric(appLogger)

	// 2. Load Documents
	providers, loadedRoots, err := loadProviders(cfg, appLogger, cwd)
	if err != nil {
		return cli.Exit(err.Error(), 1)
	}

	if len(providers) == 0 {
		return cli.Exit("No valid sources configured. Please check your .kex.yaml.", 1)
	}

	// 3. Create Composite Repository
	compositeProvider := fs.NewCompositeProvider(providers)
	repo := fs.New(compositeProvider, appLogger)

	if err := repo.Load(); err != nil {
		return cli.Exit(fmt.Sprintf("Fatal: failed to load documents: %v", err), 1)
	}

	if err := validateRepository(repo); err != nil {
		return err
	}

	// 4. Log Startup Stats
	logger.Info("Kex Server Starting...")
	logger.Info("Roots: %v", loadedRoots)

	// Stats
	var loadedIDs []string
	for id := range repo.Documents {
		loadedIDs = append(loadedIDs, id)
	}

	logger.Info("Documents Loaded: %d, IDs=%v", len(repo.Documents), loadedIDs)

	// Fail if no documents found (per user request)
	if len(repo.Documents) == 0 {
		return cli.Exit("Error: No documents found in any sources. Please check your source/references path.", 1)
	}

	if len(repo.Errors) > 0 {
		logger.Info("Load Errors: %d", len(repo.Errors))
	} else {
		logger.Info("Load Status: OK")
	}

	defer logger.Info("Kex Server Stopping...")

	return startServer(repo)
}

func isURL(s string) bool {
	return strings.HasPrefix(s, "http://") || strings.HasPrefix(s, "https://")
}

func resolveLogger(c *cli.Context, cfg config.Config, projectRoot string) (logger.Logger, error) {
	logFile := c.String("log-file")
	if logFile == "" {
		logFile = cfg.Logging.File
	}

	if logFile != "" && !filepath.IsAbs(logFile) && projectRoot != "." {
		logFile = filepath.Join(projectRoot, logFile)
	}

	if logFile != "" {
		fl, err := logger.NewFileLogger(logFile)
		if err != nil {
			return nil, fmt.Errorf("failed to open log file '%s': %w", logFile, err)
		}
		return fl, nil
	}
	return logger.NewStderrLogger(), nil
}

func validateRepository(repo *fs.Indexer) error {
	v := validator.New([]validator.ValidationRule{
		&validator.TitleRequiredRule{},
	})
	report := v.Validate(repo)

	if len(report.GlobalErrors) > 0 {
		for _, e := range report.GlobalErrors {
			fmt.Fprintf(os.Stderr, "Error: %v\n", e)
		}
		return cli.Exit("Failed to start due to document errors. Run 'kex check' for details.", 1)
	}

	if !report.Valid {
		for _, err := range report.GlobalErrors {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		}
		for _, doc := range report.Documents {
			for _, err := range doc.Errors {
				fmt.Fprintf(os.Stderr, "Error [%s]: %v\n", doc.ID, err)
			}
		}
		return cli.Exit("Validation failed (documents contain errors). Run 'kex check' for details.", 1)
	}
	return nil
}

func startServer(repo *fs.Indexer) error {
	searchUC := search.New(repo)
	retrieveUC := retrieve.New(repo)
	srv := mcp.New(searchUC, retrieveUC)
	fmt.Fprintf(os.Stderr, "Server listening on stdio...\n")
	if err := srv.Serve(); err != nil {
		logger.Error("Server error: %v", err)
		return cli.Exit(fmt.Sprintf("Server error: %v", err), 1)
	}
	return nil
}

func loadProviders(cfg config.Config, l logger.Logger, cwd string) ([]fs.DocumentProvider, []string, error) {
	var providers []fs.DocumentProvider
	var loadedRoots []string

	// Helper to add provider
	addProvider := func(pathOrURL string, isReference bool) error {
		if isURL(pathOrURL) {
			// Remote Provider
			token := os.Getenv("KEX_REMOTE_TOKEN")
			if token == "" && cfg.RemoteToken != "" {
				token = cfg.RemoteToken
			}

			fmt.Fprintf(os.Stderr, "Source: Remote (%s)\n", pathOrURL)
			providers = append(providers, fs.NewRemoteProvider(pathOrURL, token, l))
			loadedRoots = append(loadedRoots, pathOrURL)
			return nil
		}

		// Local Provider
		var fullPath string
		if filepath.IsAbs(pathOrURL) {
			fullPath = pathOrURL
		} else {
			fullPath = filepath.Join(cwd, pathOrURL)
		}

		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			if isReference {
				return fmt.Errorf("reference '%s' not found", pathOrURL)
			}
			return fmt.Errorf("source '%s' not found", pathOrURL)
		}

		fmt.Fprintf(os.Stderr, "Source: Local (%s)\n", fullPath)
		providers = append(providers, fs.NewLocalProvider(fullPath, l))
		loadedRoots = append(loadedRoots, fullPath)
		return nil
	}

	// Load Local Source
	if cfg.Source != "" {
		if err := addProvider(cfg.Source, false); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to load source '%s': %v\n", cfg.Source, err)
		}
	}

	// Load References
	for _, ref := range cfg.References {
		if err := addProvider(ref, true); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to load reference '%s': %v\n", ref, err)
		}
	}

	return providers, loadedRoots, nil
}
