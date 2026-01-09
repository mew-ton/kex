package search

import (
	"reflect"
	"testing"

	"github.com/mew-ton/kex/internal/domain"
)

// MockRepository for testing
type MockRepository struct {
	SearchFunc  func(keywords []string, scopes []string, exactScopeMatch bool) []*domain.Document
	GetByIDFunc func(id string) (*domain.Document, bool)
}

func (m *MockRepository) Search(keywords []string, scopes []string, exactScopeMatch bool) []*domain.Document {
	if m.SearchFunc != nil {
		return m.SearchFunc(keywords, scopes, exactScopeMatch)
	}
	return nil
}

// Unused in this test but required by interface
func (m *MockRepository) GetAll() []*domain.Document { return nil }
func (m *MockRepository) GetErrors() []error         { return nil }
func (m *MockRepository) GetByID(id string) (*domain.Document, bool) {
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(id)
	}
	return nil, false
}
func (m *MockRepository) Load() error { return nil }

func TestUseCase_Execute(t *testing.T) {
	tests := []struct {
		name           string
		keywords       []string
		filePath       string
		expectedScopes []string
	}{
		{
			name:           "it should include typescript scopes for .ts file",
			keywords:       []string{"testing"},
			filePath:       "src/app.ts",
			expectedScopes: []string{"coding", "typescript", "javascript", "frontend", "vcs"},
		},
		{
			name:           "it should include go scopes for .go file",
			keywords:       []string{"testing"},
			filePath:       "main.go",
			expectedScopes: []string{"coding", "go", "backend", "vcs"},
		},
		{
			name:           "it should include documentation scopes for .md file",
			keywords:       []string{"deploy"},
			filePath:       "README.md",
			expectedScopes: []string{"documentation", "vcs"},
		},
		{
			name:           "it should fallback to default scopes for unknown extension",
			keywords:       []string{"fix"},
			filePath:       "unknown.xyz",
			expectedScopes: []string{"coding", "vcs"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockRepository{
				SearchFunc: func(keywords []string, scopes []string, exactScopeMatch bool) []*domain.Document {
					// Verify that the correct scopes are passed to the repository
					if !reflect.DeepEqual(scopes, tt.expectedScopes) {
						t.Errorf("Execute() passed scopes = %v, want %v", scopes, tt.expectedScopes)
					}
					return []*domain.Document{{ID: "doc-1", Title: "Result 1"}}
				},
			}

			uc := New(mockRepo)
			result := uc.Execute(tt.keywords, tt.filePath, false)

			if len(result.Documents) != 1 {
				t.Errorf("Execute() expected 1 document, got %d", len(result.Documents))
			}
		})
	}
}
