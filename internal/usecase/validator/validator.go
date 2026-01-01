package validator

import (
	"fmt"
	"path/filepath"
	"strings"

	"kex/internal/domain"
)

// Data Structures for Reporting

type ValidationReport struct {
	Valid        bool             `json:"valid"`
	Stats        CheckStats       `json:"stats"`
	Documents    []DocumentReport `json:"documents,omitempty"`
	GlobalErrors []string         `json:"global_errors,omitempty"`
}

type CheckStats struct {
	Total         int `json:"total"`
	Adopted       int `json:"adopted"`
	Draft         int `json:"draft"`
	AdoptedErrors int `json:"adopted_errors"`
	DraftWarnings int `json:"draft_warnings"`
	ParseErrors   int `json:"parse_errors"`
}

type DocumentReport struct {
	ID     string   `json:"id"`
	Path   string   `json:"path"`
	Status string   `json:"status"`
	Errors []string `json:"errors,omitempty"`
}

// Validate performs the validation logic using the repository
func Validate(repo domain.DocumentRepository) ValidationReport {
	docs := repo.GetAll()
	loadErrors := repo.GetErrors()

	report := ValidationReport{
		Valid: true,
		Stats: CheckStats{
			Total: len(docs),
		},
	}

	// Report Parsing Errors
	for _, err := range loadErrors {
		report.GlobalErrors = append(report.GlobalErrors, err.Error())
		report.Stats.ParseErrors++
	}

	// Validate Documents
	for _, doc := range docs {
		docReport := DocumentReport{
			ID:     doc.ID,
			Path:   doc.Path,
			Status: string(doc.Status),
		}

		if doc.Status == domain.StatusDraft {
			report.Stats.Draft++
		} else {
			report.Stats.Adopted++
		}

		if err := validateDocument(doc); err != nil {
			docReport.Errors = append(docReport.Errors, err.Error())

			if doc.Status == domain.StatusDraft {
				report.Stats.DraftWarnings++
			} else {
				report.Stats.AdoptedErrors++
			}
		}

		report.Documents = append(report.Documents, docReport)
	}

	if report.Stats.AdoptedErrors > 0 || report.Stats.ParseErrors > 0 {
		report.Valid = false
	}

	return report
}

func validateDocument(doc *domain.Document) error {
	filename := filepath.Base(doc.Path)
	ext := filepath.Ext(filename)
	basename := strings.TrimSuffix(filename, ext)

	if basename != doc.ID {
		return fmt.Errorf("filename must match id (filename: %s, id: %s)", filename, doc.ID)
	}
	return nil
}
