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

	// 1. Resolve Project Root or Remote URL
	arg := c.Args().First()
	isRemote := strings.HasPrefix(arg, "http://") || strings.HasPrefix(arg, "https://")

	var repo *fs.Indexer
	var cfg config.Config
	var err error

	// Resolve Config & Logger
	// We try to load config early to get logging preferences
	projectRoot := "."
	if !isRemote && arg != "" {
		projectRoot = arg
	}
	cfg, err = config.Load(projectRoot)
	// Ignore error for now, will handle specific errors later

	// Init Logger
	var appLogger logger.Logger
	logFile := c.String("log-file")
	if logFile == "" {
		logFile = cfg.Logging.File
	}

	if logFile != "" {
		// Resolve absolute path if needed?
		// Ensure directory exists?
		// For now simple open.
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

	if isRemote {
		// Resolve Token: KEX_REMOTE_TOKEN > .kex.yaml (remoteToken)
		pathOrUrl := arg
		token := os.Getenv("KEX_REMOTE_TOKEN")

		// Check config only if token missing? Or always overrides?
		// Existing logic:
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

		provider := fs.NewRemoteProvider(pathOrUrl, token, appLogger)
		repo = fs.New(provider, appLogger)
		if err := repo.Load(); err != nil {
			return cli.Exit(fmt.Sprintf("Failed to load remote index: %v", err), 1)
		}
	} else {
		// Local Mode
		// projectRoot is already set above
		// cfg is already loaded above
		if err != nil { // Err from config.Load
			fmt.Fprintf(os.Stderr, "Warning: failed to load config: %v\n", err)
		}

		// Resolve Content Directory
		root := filepath.Join(projectRoot, cfg.Root)
		if c.String("root") != "" {
			root = c.String("root")
		}

		if _, err := os.Stat(root); os.IsNotExist(err) {
			return cli.Exit(fmt.Sprintf("Error: directory '%s' not found. Run 'kex init'?", root), 1)
		}

		provider := fs.NewLocalProvider(root, appLogger)
		repo = fs.New(provider, appLogger)
		if err := repo.Load(); err != nil {
			return cli.Exit(fmt.Sprintf("Fatal: failed to load documents: %v", err), 1)
		}
	}

	// 5. Strict validation on start (Use Case)
	// We use the Validator use case to determine validity
	// Note: Remote documents are assumed valid (or validated at build time).
	// But validator checks for structure/missing fields.
	// Initialize Validator with default rules for startup check
	v := validator.New([]validator.ValidationRule{
		&validator.IDRequiredRule{},
		&validator.TitleRequiredRule{},
		&validator.FilenameMatchRule{},
	})
	report := v.Validate(repo)

	// Check for Parse Errors (Critical for start)
	if len(report.GlobalErrors) > 0 {
		for _, e := range report.GlobalErrors {
			fmt.Fprintf(os.Stderr, "Error: %v\n", e)
		}
		return cli.Exit("Failed to start due to document errors. Run 'kex check' for details.", 1)
	}

	// Note: We might allow "AdoptedErrors" but block on "ParseErrors".
	// For strictness, let's block if Valid is false (excluding Drafts which are warnings)
	// But validator.Validate returns Valid=false if AdoptedErrors > 0.
	// As per previous design: "Failed to start due to document errors".
	if !report.Valid {
		// Print errors to stderr for debugging/user info
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

	// 4. Start Server (Interface)
	srv := mcp.New(repo, appLogger)
	fmt.Fprintf(os.Stderr, "Server listening on stdio...\n")
	if err := srv.Serve(); err != nil {
		return cli.Exit(fmt.Sprintf("Server error: %v", err), 1)
	}

	return nil
}
