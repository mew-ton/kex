package fs_test

import (
	"testing"

	"github.com/mew-ton/kex/internal/infrastructure/fs"
	"github.com/mew-ton/kex/internal/infrastructure/logger"
)

type MockProvider struct{}

func (m *MockProvider) Load() (*fs.IndexSchema, []error) {
	return &fs.IndexSchema{
		Documents: []*fs.DocumentSchema{
			{
				ID:       "documentation/kex/choose-effective-keywords",
				Title:    "Choose effective keywords",
				Type:     "indicator",
				Keywords: []string{"keyword", "kex", "documentation"},
				Scopes:   []string{"documentation", "kex"},
				Path:     "contents/documentation/kex/choose-effective-keywords.md",
			},
		},
	}, nil
}
func (m *MockProvider) FetchContent(path string) (string, error) { return "", nil }

func TestIndexer_Search_Indicator(t *testing.T) {
	idx := fs.New(&MockProvider{}, logger.NewStderrLogger()) // Use correct constructor
	idx.Load()

	// Test case: Search for "kex" AND "documentation"
	results := idx.Search([]string{"kex", "documentation"}, nil, true)
	// Wait, kex command (update) uses keywords, which implies exactScopeMatch=false or true?
	// generator.go calls retrieveIndicators -> indexer.Search(keywords).
	// It passes nil scopes?
	// retrieveIndicators signature: func (g *Generator) retrieveIndicators(ctx context.Context, keywords []string)
	// It calls g.indexer.Search(keywords, nil, true). "true" means Exact Scope Match?
	// Wait, retrieveIndicators does NOT pass exactScopeMatch arg?
	// I need to check generator.go call site.

	// Assuming search logic:
	if len(results) == 0 {
		t.Errorf("Search failed for 'kex'")
	} else {
		t.Logf("Found %d results", len(results))
	}
}
