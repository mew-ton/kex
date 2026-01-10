package fs

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/mew-ton/kex/internal/domain"
	"github.com/mew-ton/kex/internal/infrastructure/logger"
)

func TestIndexer_Search(t *testing.T) {
	// Setup temp dir with sample docs
	tmpDir := t.TempDir()
	// Create coding directory to test scopes
	if err := os.Mkdir(filepath.Join(tmpDir, "coding"), 0755); err != nil {
		t.Fatal(err)
	}

	doc1 := `---
title: Doc 1
status: adopted
keywords: [apple, banana]
---
Content 1`

	// Doc 3 in coding directory
	doc3 := `---
title: Doc 3 Coding
status: adopted
keywords: [zebra]
---
Content 3`

	doc2 := `---
title: Doc 2
status: adopted
keywords: [banana, cherry]
---
Content 2`

	if err := os.WriteFile(filepath.Join(tmpDir, "doc1.md"), []byte(doc1), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(tmpDir, "doc2.md"), []byte(doc2), 0644); err != nil {
		t.Fatal(err)
	}

	if err := os.WriteFile(filepath.Join(tmpDir, "coding", "doc3.md"), []byte(doc3), 0644); err != nil {
		t.Fatal(err)
	}

	t.Run("should load valid documents", func(t *testing.T) {
		l := &logger.NoOpLogger{}
		provider := NewLocalProvider(tmpDir, l)
		idx := New(provider, l)
		err := idx.Load()
		if err != nil {
			t.Fatalf("Load failed: %v", err)
		}
	})

	// Initialize Indexer for search tests
	l := &logger.NoOpLogger{}
	provider := NewLocalProvider(tmpDir, l)
	idx := New(provider, l)
	if err := idx.Load(); err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	// Test Search
	tests := []struct {
		name            string
		keywords        []string
		exactScopeMatch bool
		wantIDs         []string
	}{
		{
			name:     "it should return relevant documents for apple",
			keywords: []string{"apple"},
			wantIDs:  []string{"doc1"},
		},
		{
			name:     "it should return all documents containing banana",
			keywords: []string{"banana"},
			wantIDs:  []string{"doc1", "doc2"},
		},
		{
			name:     "it should return no documents for unknown keywords",
			keywords: []string{"durian"},
			wantIDs:  nil,
		},

		{
			name:     "it should match documents by title words",
			keywords: []string{"Doc"},
			wantIDs:  []string{"doc1", "doc2", "coding.doc3"},
		},
		{
			name:            "it should return docs in scope when exactScopeMatch is true",
			keywords:        []string{"coding"},
			exactScopeMatch: true,
			wantIDs:         []string{"coding.doc3"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := idx.Search(tt.keywords, nil, tt.exactScopeMatch)
			if len(got) != len(tt.wantIDs) {
				t.Errorf("Search() got %d docs, want %d", len(got), len(tt.wantIDs))
			}
			// Check IDs presence (order doesn't matter strictly, but here manageable)
			for _, wantID := range tt.wantIDs {
				found := false
				for _, d := range got {
					if d.ID == wantID {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("Search() missing ID %s", wantID)
				}
			}
		})
	}
}

func TestIndexer_Load_DefaultStatus(t *testing.T) {
	// Setup temp dir with sample doc missing status
	tmpDir := t.TempDir()

	doc := `---
title: Doc No Status
keywords: [test]
---
Content`

	if err := os.WriteFile(filepath.Join(tmpDir, "doc.md"), []byte(doc), 0644); err != nil {
		t.Fatal(err)
	}

	l := &logger.NoOpLogger{}
	provider := NewLocalProvider(tmpDir, l)

	// We want to verify that the provider loads it with StatusAdopted
	// The Indexer.Load will use the provider's schema.
	// But we can check the provider output directly or through Indexer.
	// Let's check through Indexer to be sure end-to-end.

	idx := New(provider, l)
	// Default IncludeDrafts is false.
	// If it defaults to Draft, it won't be in Indexer.Documents if we rely on parseDocuments logic?
	// Wait, Indexer.Load calls Provider.Load which returns Schema.
	// detailed logic:
	// 1. Provider.Load calls collectMarkdownFiles -> ParseDocument
	// 2. Provider sets default status if missing.
	// 3. Provider returns schema.
	// 4. Indexer.Load converts Schema to Documents.

	if err := idx.Load(); err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	// It should be present and have StatusAdopted
	d, ok := idx.GetByID("doc")
	if !ok {
		// If it was Draft and IncludeDrafts=false (default), it might still be in Documents map?
		// Checking Indexer.Load implementation:
		// It iterates schema.Documents and adds them.
		// It does NOT filter by Draft status in the loop from schema!
		// Wait, let's check Indexer.Load again.
		// lines 56-75: iterates schema.Documents, creates domain.Document, calls i.addDocument.
		// DOES NOT CHECK IncludeDrafts there.
		// So it should be there regardless of status.
		t.Fatalf("Document not found in index")
	}

	if d.Status != domain.StatusAdopted {
		t.Errorf("expected status %q, got %q", domain.StatusAdopted, d.Status)
	}
}
