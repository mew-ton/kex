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
	},
	Action: runStart,
}

func runStart(c *cli.Context) error {
	fmt.Fprintf(os.Stderr, "Starting Kex Server...\n")

	args := c.Args().Slice()
	if len(args) == 0 {
		args = []string{"."}
	}

	// 1. Setup Logger (Use config from first local path or CWD)
	// We make a best-effort attempt to load config for logging purposes
	primaryRoot := "."
	for _, arg := range args {
		if !isURL(arg) {
			primaryRoot = arg
			break
		}
	}

	// Load primary config for logging/global settings
	primaryCfg, _ := config.Load(primaryRoot)

	// Setup Logger
	appLogger, err := resolveLogger(c, primaryCfg, primaryRoot)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: logger setup failed: %v. Using stderr.\n", err)
		appLogger = logger.NewStderrLogger()
	}
	logger.SetGeneric(appLogger)

	// 2. Create Providers
	var providers []fs.DocumentProvider
	var loadedRoots []string

	for _, arg := range args {
		if isURL(arg) {
			// Remote Provider
			token := os.Getenv("KEX_REMOTE_TOKEN")
			// Try to use primary config's token as fallback if available
			if token == "" && primaryCfg.RemoteToken != "" {
				token = primaryCfg.RemoteToken
			}

			fmt.Fprintf(os.Stderr, "Source: Remote (%s)\n", arg)
			p := fs.NewRemoteProvider(arg, token, appLogger)
			providers = append(providers, p)
			loadedRoots = append(loadedRoots, arg)
			continue
		}

		// Local Provider (Project Root)
		projectRoot := arg
		if _, err := os.Stat(projectRoot); os.IsNotExist(err) {
			return cli.Exit(fmt.Sprintf("Error: directory '%s' not found.", projectRoot), 1)
		}

		cfg, err := config.Load(projectRoot)
		if err != nil {
			logger.Error("Failed to load config for %s: %v", projectRoot, err)
			// Continue with defaults? Or fail?
			// config.Load returns defaults if file missing, so this error is for corrupt yaml etc.
			// Let's warn and continue with defaults returned by Load?
			// But Load returns error only if yaml unmarshal fails.
		}

		// Create a provider for each Source defined in this project
		for _, source := range cfg.Sources {
			fullSourcePath := filepath.Join(projectRoot, source)
			if _, err := os.Stat(fullSourcePath); os.IsNotExist(err) {
				// Warn but don't fail hard? Or fail? The user configured it.
				// Consistent with old behavior: explicit root must exist.
				return cli.Exit(fmt.Sprintf("Error: source directory '%s' not found in project '%s'.", source, projectRoot), 1)
			}

			fmt.Fprintf(os.Stderr, "Source: Local (%s)\n", fullSourcePath)
			p := fs.NewLocalProvider(fullSourcePath, appLogger)
			providers = append(providers, p)
		}
		loadedRoots = append(loadedRoots, projectRoot)
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
