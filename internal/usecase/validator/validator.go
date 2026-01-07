package validator

import (
	"github.com/mew-ton/kex/internal/domain"
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

type Validator struct {
	Rules []ValidationRule
}

func New(rules []ValidationRule) *Validator {
	return &Validator{Rules: rules}
}

// Validate performs the validation logic using the repository
func (v *Validator) Validate(repo domain.DocumentRepository) ValidationReport {
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

		// Run all validation rules
		var docErrors []string
		for _, rule := range v.Rules {
			if err := rule.Validate(doc); err != nil {
				docErrors = append(docErrors, err.Error())
			}
		}

		if len(docErrors) > 0 {
			docReport.Errors = docErrors
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
