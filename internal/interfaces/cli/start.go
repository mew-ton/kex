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
	"github.com/mew-ton/kex/internal/usecase/validator"

	"github.com/urfave/cli/v2"
)

var StartCommand = &cli.Command{
	Name:  "start",
	Usage: "Start MCP server",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "root",
			Usage:   "Path to guidelines directory",
			Aliases: []string{"r"},
		},
		&cli.StringFlag{
			Name:  "log-file",
			Usage: "Path to log file",
		},
	},
	Action: runStart,
}

func runStart(c *cli.Context) error {
	fmt.Fprintf(os.Stderr, "Starting Kex Server...\n")

	arg := c.Args().First()
	isRemote := strings.HasPrefix(arg, "http://") || strings.HasPrefix(arg, "https://")

	// 1. Resolve Project Root
	projectRoot := "."
	if !isRemote && arg != "" {
		projectRoot = arg
	}

	// 2. Load Config
	cfg, err := config.Load(projectRoot)
	if err != nil {
		// Proceed with default config if load fails (consistent with previous behavior)
	}

	// 3. Setup Logger
	appLogger, err := resolveLogger(c, cfg, projectRoot)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: logger setup failed: %v. Using stderr.\n", err)
		appLogger = logger.NewStderrLogger()
	}
	logger.SetGeneric(appLogger)

	repo, err := createRepository(c, cfg, appLogger, arg, isRemote, projectRoot)
	if err != nil {
		return err
	}

	if err := validateRepository(repo); err != nil {
		return err
	}

	// 4. Log Startup Stats
	logger.Info("Kex Server Starting...")
	logger.Info("Root: %s", projectRoot)
	logger.Info("Mode: %s", func() string {
		if isRemote {
			return "Remote"
		}
		return "Local"
	}())
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

func createRepository(c *cli.Context, cfg config.Config, l logger.Logger, arg string, isRemote bool, projectRoot string) (*fs.Indexer, error) {
	if isRemote {
		pathOrUrl := arg
		token := os.Getenv("KEX_REMOTE_TOKEN")
		if token == "" {
			if cfg.RemoteToken != "" {
				token = cfg.RemoteToken
			}
		}

		fmt.Fprintf(os.Stderr, "Source: Remote (%s)\n", pathOrUrl)
		if token != "" {
			fmt.Fprintf(os.Stderr, "Auth: Token provided\n")
		} else {
			fmt.Fprintf(os.Stderr, "Auth: None\n")
		}

		provider := fs.NewRemoteProvider(pathOrUrl, token, l)
		repo := fs.New(provider, l)
		if err := repo.Load(); err != nil {
			return nil, cli.Exit(fmt.Sprintf("Failed to load remote index: %v", err), 1)
		}
		return repo, nil
	}

	// Local Mode
	root := filepath.Join(projectRoot, cfg.Root)
	if c.String("root") != "" {
		root = c.String("root")
	}

	if _, err := os.Stat(root); os.IsNotExist(err) {
		return nil, cli.Exit(fmt.Sprintf("Error: directory '%s' not found. Run 'kex init'?", root), 1)
	}

	provider := fs.NewLocalProvider(root, l)
	repo := fs.New(provider, l)
	if err := repo.Load(); err != nil {
		return nil, cli.Exit(fmt.Sprintf("Fatal: failed to load documents: %v", err), 1)
	}
	return repo, nil
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
	srv := mcp.New(repo)
	fmt.Fprintf(os.Stderr, "Server listening on stdio...\n")
	if err := srv.Serve(); err != nil {
		logger.Error("Server error: %v", err)
		return cli.Exit(fmt.Sprintf("Server error: %v", err), 1)
	}
	return nil
}
