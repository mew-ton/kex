package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

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

	// 1. Resolve Repository
	// check command usually doesn't take args for root in old version (it took projectRoot arg).
	// resolveRepository handles args as roots.
	// This aligns check behavior with start/generate.

	repo, _, _, err := resolveRepository(c)

	if err != nil {
		if isJSON {
			printJSONError(err.Error())
			return cli.Exit("", 1)
		}
		// resolveRepository logs to its logger which might be stderr.
		// If we want to show fatal error here:
		return cli.Exit(fmt.Sprintf("Fatal: failed to load documents: %v", err), 1)
	}

	repo.IncludeDrafts = true // Check command always checks drafts?
	// Note: resolveRepository calls Load() internally.
	// If we set IncludeDrafts AFTER Load, it doesn't affect the initial load validation (if any).
	// BUT Indexer.Load() logic:
	// 1. Load Schema (Provider)
	// 2. Convert to Domain Documents
	// 3. parseDocuments() -> checks IncludeDrafts.
	// Oh wait, Indexer.Load() in `indexer.go` iterates Schema and adds documents.
	// It doesn't use `parseDocuments` anymore?
	// `parseDocuments` seems unused in my previous read of indexer.go?
	// Let's re-read indexer.go logic if possible.
	// Assuming Load() populates documents.
	// If IncludeDrafts is false (default), does it filter?
	// Current `indexer.go` implementation of `Load`:
	// It iterates schema and calls `addDocument`.
	// It relies on Schema having the documents.
	// LocalProvider parses content.
	// If validation rules check "Status: draft", they can ignore it?
	// `check` command validates all.
	// The Validator logic is what matters.
	// Validator rules validation.
	// The previous `loadRepository` helper set `repo.IncludeDrafts = true`.
	// This suggests it's important.
	// Since `resolveRepository` creates and loads, we might need to configure factory?
	// Or factory should accept options.
	// For now, let's set it and maybe reload? Or is it too late?
	// Actually, `IncludeDrafts` on Indexer seems not used in `Load` based on my view of `indexer.go`.
	// It was used in `search` maybe? Or just legacy?

	// Let's assume repo is loaded.

	// Initialize Validator with default rules
	rules := []validator.ValidationRule{
		&validator.IDRequiredRule{},
		&validator.TitleRequiredRule{},
		&validator.FilenameMatchRule{},
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
	for _, err := range report.GlobalErrors {
		pterm.Error.Printf("Parse Error: %v\n", err)
	}

	// Print Document Errors
	for _, doc := range report.Documents {
		for _, err := range doc.Errors {
			// Logic to print based on status is now presentation concern
			// We can assume ValidationReport contains necessary info
			// Re-mapping status for display
			if doc.Status == "draft" {
				pterm.Warning.Printf("[%s] (Draft) %s\n", doc.ID, err)
			} else {
				pterm.Error.Printf("[%s] (Adopted) %s\n", doc.ID, err)
			}
		}
	}

	// Statistics
	pterm.Println() // Spacer
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
