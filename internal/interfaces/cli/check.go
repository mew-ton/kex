package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/mew-ton/kex/internal/infrastructure/config"
	"github.com/mew-ton/kex/internal/infrastructure/fs"
	"github.com/mew-ton/kex/internal/infrastructure/logger"
	"github.com/mew-ton/kex/internal/usecase/validator"

	"github.com/pterm/pterm"
	"github.com/urfave/cli/v2"
)

var CheckCommand = &cli.Command{
	Name:  "check",
	Usage: "Validate documents",
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:  "json",
			Usage: "Output results in JSON format",
		},
	},
	Action: runCheck,
}

func runCheck(c *cli.Context) error {
	isJSON := c.Bool("json")

	if !isJSON {
		pterm.DefaultSection.Println("Checking documents...")
	}

	// 1. Resolve Project Root
	projectRoot := c.Args().First()
	if projectRoot == "" {
		projectRoot = "."
	}

	cfg, err := resolveConfig(projectRoot)
	if err != nil {
		if !isJSON {
			pterm.Warning.Printf("Failed to load config, using defaults: %v\n", err)
		}
	}

	// 2. Resolve Content Directory
	// For Check, we iterate all sources and check the aggregate?
	// Or strict check each source?
	// Let's use CompositeProvider logic here too for consistency with Start.

	repo, err := loadRepository(projectRoot, cfg, !isJSON)
	if err != nil {
		if isJSON {
			printJSONError(err.Error())
			return cli.Exit("", 1)
		}
		return cli.Exit(fmt.Sprintf("Fatal: failed to load documents: %v", err), 1)
	}

	// Initialize Validator with default rules
	rules := []validator.ValidationRule{
		&validator.TitleRequiredRule{},
	}
	v := validator.New(rules)
	report := v.Validate(repo)

	if isJSON {
		printJSONReport(report)
	} else {
		printHumanReport(report, report.Stats.Total)
	}

	if !report.Valid {
		return cli.Exit("", 1)
	}

	return nil
}

func resolveConfig(projectRoot string) (config.Config, error) {
	return config.Load(projectRoot)
}

func loadRepository(projectRoot string, cfg config.Config, showSpinner bool) (*fs.Indexer, error) {
	var spinner *pterm.SpinnerPrinter
	if showSpinner {
		spinner, _ = pterm.DefaultSpinner.Start("Loading documents...")
	}

	// Use NoOpLogger for Check command to avoid clutter
	l := &logger.NoOpLogger{}

	var providers []fs.DocumentProvider

	// Helper to add provider
	addProvider := func(pathOrURL string, isReference bool) error {
		if strings.HasPrefix(pathOrURL, "http://") || strings.HasPrefix(pathOrURL, "https://") {
			// Remote Provider
			token := os.Getenv("KEX_REMOTE_TOKEN")
			if token == "" && cfg.RemoteToken != "" {
				token = cfg.RemoteToken
			}
			providers = append(providers, fs.NewRemoteProvider(pathOrURL, token, l))
			return nil
		}

		// Local Provider
		var fullPath string
		if filepath.IsAbs(pathOrURL) {
			fullPath = pathOrURL
		} else {
			fullPath = filepath.Join(projectRoot, pathOrURL)
		}

		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			if isReference {
				return fmt.Errorf("reference '%s' not found", pathOrURL)
			}
			return fmt.Errorf("source '%s' not found", pathOrURL)
		}

		providers = append(providers, fs.NewLocalProvider(fullPath, l))
		return nil
	}

	// Load Local Source
	if cfg.Source != "" {
		if err := addProvider(cfg.Source, false); err != nil {
			if spinner != nil {
				spinner.Fail(err.Error())
			}
			return nil, err
		}
	}

	// Load References
	for _, ref := range cfg.References {
		if err := addProvider(ref, true); err != nil {
			if spinner != nil {
				spinner.Fail(err.Error())
			}
			return nil, err
		}
	}

	if len(providers) == 0 {
		if spinner != nil {
			spinner.Fail("No valid content directories found")
		}
		return nil, fmt.Errorf("no valid sources found in config (source or references)")
	}

	composite := fs.NewCompositeProvider(providers)
	repo := fs.New(composite, l)
	repo.IncludeDrafts = true
	if err := repo.Load(); err != nil {
		if spinner != nil {
			spinner.Fail("Failed to load documents")
		}
		return nil, err
	}
	if spinner != nil {
		spinner.Success("Documents loaded")
	}
	return repo, nil
}

// Presentation Logic

func printJSONReport(report validator.ValidationReport) {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	encoder.Encode(report)
}

func printJSONError(msg string) {
	// Minimal JSON error output
	fmt.Printf(`{"valid": false, "global_errors": ["%s"]}`+"\n", strings.ReplaceAll(msg, "\"", "\\\""))
}

func printHumanReport(report validator.ValidationReport, totalDocs int) {
	// Print Global Errors
	if len(report.GlobalErrors) > 0 {
		pterm.DefaultSection.Println("Parse Errors")
		for _, err := range report.GlobalErrors {
			pterm.Error.Printf("- %v\n", err)
		}
		pterm.Println()
	}

	// Filter documents
	var drafts []validator.DocumentReport
	var errorDocs []validator.DocumentReport

	for _, doc := range report.Documents {
		if doc.Status == "draft" {
			drafts = append(drafts, doc)
		}
		if len(doc.Errors) > 0 {
			errorDocs = append(errorDocs, doc)
		}
	}

	// Print Drafts
	if len(drafts) > 0 {
		pterm.DefaultSection.Println("Drafts")
		for _, doc := range drafts {
			pterm.Println(fmt.Sprintf("- %s", doc.ID))
		}
		pterm.Println()
	}

	// Print Document Errors
	if len(errorDocs) > 0 {
		pterm.DefaultSection.Println("Validation Errors")
		for _, doc := range errorDocs {
			for _, err := range doc.Errors {
				if doc.Status == "draft" {
					pterm.Warning.Printf("[%s] (Draft) %s\n", doc.ID, err)
				} else {
					pterm.Error.Printf("[%s] (Adopted) %s\n", doc.ID, err)
				}
			}
		}
		pterm.Println()
	}

	// Statistics
	pterm.DefaultSection.Println("Statistics")

	tableData := [][]string{
		{"Metric", "Count"},
		{"Total Documents", fmt.Sprintf("%d", totalDocs)},
		{"Adopted", fmt.Sprintf("%d", report.Stats.Adopted)},
		{"Draft", fmt.Sprintf("%d", report.Stats.Draft)},
		{"Parse Errors", fmt.Sprintf("%d", report.Stats.ParseErrors)},
		{"Adopted Errors", fmt.Sprintf("%d", report.Stats.AdoptedErrors)},
		{"Draft Warnings", fmt.Sprintf("%d", report.Stats.DraftWarnings)},
	}

	pterm.DefaultTable.WithHasHeader().WithData(tableData).Render()
	pterm.Println()

	if !report.Valid {
		pterm.Error.Println("Check failed. Please fix errors.")
	} else {
		pterm.Success.Println("All checks passed.")
	}
}
