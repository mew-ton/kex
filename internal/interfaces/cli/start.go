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

	// 1. Resolve Working Directory
	cwd, err := resolveCwd(c)
	if err != nil {
		return err
	}

	// 2. Load Configuration
	cfg, err := loadConfiguration(cwd)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: config load failed: %v. Using defaults.\n", err)
	}

	// 3. Setup Logger
	appLogger, err := setupAppLogger(c, cfg, cwd)
	if err != nil {
		return err
	}
	logger.SetGeneric(appLogger)

	// 4. Create and Prepare Repository
	repo, loadedRoots, err := createRepository(cfg, appLogger, cwd)
	if err != nil {
		return cli.Exit(err.Error(), 1)
	}

	// 5. Log Startup Stats & Checks
	logStartupStats(repo, loadedRoots)

	if err := checkRepositoryState(repo); err != nil {
		return err
	}

	defer logger.Info("Kex Server Stopping...")

	return startServer(repo)
}

func resolveCwd(c *cli.Context) (string, error) {
	cwd := c.String("cwd")
	if cwd == "" {
		return os.Getwd()
	}
	return filepath.Abs(cwd)
}

func loadConfiguration(cwd string) (config.Config, error) {
	return config.Load(cwd)
}

func setupAppLogger(c *cli.Context, cfg config.Config, cwd string) (logger.Logger, error) {
	appLogger, err := resolveLogger(c, cfg, cwd)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: logger setup failed: %v. Using stderr.\n", err)
		return logger.NewStderrLogger(), nil
	}
	return appLogger, nil
}

func createRepository(cfg config.Config, l logger.Logger, cwd string) (*fs.Indexer, []string, error) {
	providers, loadedRoots, err := loadProviders(cfg, l, cwd)
	if err != nil {
		return nil, nil, err
	}

	if len(providers) == 0 {
		return nil, nil, fmt.Errorf("no valid sources configured. Please check your .kex.yaml")
	}

	compositeProvider := fs.NewCompositeProvider(providers)
	repo := fs.New(compositeProvider, l)

	if err := repo.Load(); err != nil {
		return nil, nil, fmt.Errorf("fatal: failed to load documents: %w", err)
	}

	if err := validateRepository(repo); err != nil {
		return nil, nil, err
	}

	return repo, loadedRoots, nil
}

func logStartupStats(repo *fs.Indexer, loadedRoots []string) {
	logger.Info("Kex Server Starting...")
	logger.Info("Roots: %v", loadedRoots)

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
}

func checkRepositoryState(repo *fs.Indexer) error {
	if len(repo.Documents) == 0 {
		return cli.Exit("Error: No documents found in any sources. Please check your source/references path.", 1)
	}
	return nil
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
		return fmt.Errorf("failed to start due to document errors. Run 'kex check' for details")
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
		return fmt.Errorf("validation failed (documents contain errors). Run 'kex check' for details")
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
	factory := fs.NewProviderFactory(cfg, l)

	// Helper to add provider
	addProvider := func(pathOrURL string, isReference bool) {
		provider, resolvedPath, err := factory.CreateProvider(pathOrURL, isReference, cwd)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to load source/reference '%s': %v\n", pathOrURL, err)
			return
		}

		sourceType := "Local"
		if isURL(pathOrURL) {
			sourceType = "Remote"
		}
		fmt.Fprintf(os.Stderr, "Source: %s (%s)\n", sourceType, resolvedPath)

		providers = append(providers, provider)
		loadedRoots = append(loadedRoots, resolvedPath)
	}

	// Load Local Source
	if cfg.Source != "" {
		addProvider(cfg.Source, false)
	}

	// Load References
	for _, ref := range cfg.References {
		addProvider(ref, true)
	}

	return providers, loadedRoots, nil
}
