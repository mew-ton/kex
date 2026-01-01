package indexer

import (
	"os"
	"path/filepath"
	"testing"
)

func TestIndexer_Search(t *testing.T) {
	// Setup temp dir with sample docs
	tmpDir := t.TempDir()

	doc1 := `---
id: doc1
title: Doc 1
keywords: [apple, banana]
---
Content 1`

	doc2 := `---
id: doc2
title: Doc 2
keywords: [banana, cherry]
---
Content 2`

	if err := os.WriteFile(filepath.Join(tmpDir, "doc1.md"), []byte(doc1), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(tmpDir, "doc2.md"), []byte(doc2), 0644); err != nil {
		t.Fatal(err)
	}

	// Initialize Indexer
	idx := New(tmpDir)
	if err := idx.Load(); err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	// Test Search
	tests := []struct {
		name     string
		keywords []string
		wantIDs  []string
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := idx.Search(tt.keywords, nil)
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
