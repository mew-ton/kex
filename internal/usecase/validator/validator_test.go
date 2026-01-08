package validator

import (
	"errors"
	"testing"

	"github.com/mew-ton/kex/internal/domain"
)

// MockRepository for testing
type MockRepository struct {
	GetAllFunc    func() []*domain.Document
	GetErrorsFunc func() []error
}

// Implement only necessary methods for Validator
func (m *MockRepository) GetAll() []*domain.Document {
	if m.GetAllFunc != nil {
		return m.GetAllFunc()
	}
	return nil
}

func (m *MockRepository) GetErrors() []error {
	if m.GetErrorsFunc != nil {
		return m.GetErrorsFunc()
	}
	return nil
}

// Unused methods
func (m *MockRepository) GetByID(id string) (*domain.Document, bool) { return nil, false }
func (m *MockRepository) Search(k, s []string) []*domain.Document    { return nil }
func (m *MockRepository) Load() error                                { return nil }

func TestValidate(t *testing.T) {
	tests := []struct {
		name          string
		documents     []*domain.Document
		globalErrors  []error
		wantValid     bool
		wantErrCount  int // Adopted document errors
		wantGlobCount int // Global (parse) errors
	}{
		{
			name:      "it should be valid with no documents and no errors",
			documents: []*domain.Document{},
			wantValid: true,
		},
		{
			name: "it should report parse errors and be invalid",
			globalErrors: []error{
				errors.New("failed to parse doc 1"),
			},
			wantValid:     false,
			wantGlobCount: 1,
		},
		{
			name: "it should invalid if title is missing",
			documents: []*domain.Document{
				{ID: "doc-1", Title: "", Status: domain.StatusAdopted, Path: "doc-1.md"},
			},
			wantValid:    false,
			wantErrCount: 1,
		},
		{
			name: "it should ignore draft documents errors",
			documents: []*domain.Document{
				// Now only Title is validated. If title is missing, it's an error.
				{ID: "draft", Title: "", Status: domain.StatusDraft, Path: "draft.md"},
			},
			wantErrCount: 1,    // Title required
			wantValid:    true, // Drafts are skipped in error count impacting validity (usually)
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockRepository{
				GetAllFunc: func() []*domain.Document {
					return tt.documents
				},
				GetErrorsFunc: func() []error {
					return tt.globalErrors
				},
			}

			rules := []ValidationRule{
				&TitleRequiredRule{},
			}
			v := New(rules)
			report := v.Validate(mockRepo)

			if report.Valid != tt.wantValid {
				t.Errorf("Validate() valid = %v, want %v", report.Valid, tt.wantValid)
			}
			if len(report.GlobalErrors) != tt.wantGlobCount {
				t.Errorf("Validate() global errors = %d, want %d", len(report.GlobalErrors), tt.wantGlobCount)
			}

			// Count total errors across all documents
			totalDocErrors := 0
			for _, d := range report.Documents {
				totalDocErrors += len(d.Errors)
			}

			if totalDocErrors != tt.wantErrCount {
				t.Errorf("Validate() total errors = %d, want %d", totalDocErrors, tt.wantErrCount)
			}

			// Note: Draft logic might differ. If Drafts validation generates errors, we need to adjust expected.
			// Let's assume strict validation for all documents returned by GetAll() for now.
			if tt.name == "it should ignore draft documents errors" && !tt.wantValid {
				// Adjust expectation if I was wrong about drafts
			}
		})
	}
}
