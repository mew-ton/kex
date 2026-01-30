package fs

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/mew-ton/kex/internal/infrastructure/logger"
)

func TestIndexer_StrictScope_Subset(t *testing.T) {
	// Setup temp dir with hierarchical docs
	tmpDir := t.TempDir()

	// 1. coding/rule.md (Scope: [coding])
	// Keywords: [guideline]
	if err := os.MkdirAll(filepath.Join(tmpDir, "coding"), 0755); err != nil {
		t.Fatal(err)
	}
	docCoding := `---
type: indicator
title: Coding Rule
keywords: [guideline]
---
Content`
	if err := os.WriteFile(filepath.Join(tmpDir, "coding", "rule.md"), []byte(docCoding), 0644); err != nil {
		t.Fatal(err)
	}

	// 2. coding/go/rule.md (Scope: [coding, go])
	// Keywords: [guideline]
	if err := os.MkdirAll(filepath.Join(tmpDir, "coding", "go"), 0755); err != nil {
		t.Fatal(err)
	}
	docGo := `---
type: indicator
title: Go Rule
keywords: [guideline]
---
Content`
	if err := os.WriteFile(filepath.Join(tmpDir, "coding", "go", "rule.md"), []byte(docGo), 0644); err != nil {
		t.Fatal(err)
	}

	// 3. documentation/rule.md (Scope: [documentation])
	// Keywords: [guideline]
	if err := os.MkdirAll(filepath.Join(tmpDir, "documentation"), 0755); err != nil {
		t.Fatal(err)
	}
	docDoc := `---
type: indicator
title: Documentation Rule
keywords: [guideline]
---
Content`
	if err := os.WriteFile(filepath.Join(tmpDir, "documentation", "rule.md"), []byte(docDoc), 0644); err != nil {
		t.Fatal(err)
	}

	// 4. Root doc (Scope: [])
	// Keywords: [guideline]
	docRoot := `---
type: indicator
title: Root Rule
keywords: [guideline]
---
Content`
	if err := os.WriteFile(filepath.Join(tmpDir, "rule.md"), []byte(docRoot), 0644); err != nil {
		t.Fatal(err)
	}

	// Initialize Indexer
	l := &logger.NoOpLogger{}
	provider := NewLocalProvider(tmpDir, l)
	idx := New(provider, l)
	if err := idx.Load(); err != nil {
		t.Fatal(err)
	}

	type matchExpectation struct {
		Title       string
		ShouldMatch bool
	}

	tests := []struct {
		name     string
		keywords []string
		expected []matchExpectation
	}{
		{
			// Coding, Go match because 'coding' is subset of [coding, guideline].
			// Wait: Doc [coding, go] scopes are {coding, go}. Query scopes: {coding}.
			// {coding, go} is NOT subset of {coding}. So Go Rule should be EXCLUDED.
			// Doc [coding] scopes are {coding}. Query scopes: {coding}.
			// {coding} is subset of {coding}. So Coding Rule should be INCLUDED.
			name:     "Query [coding, guideline]",
			keywords: []string{"coding", "guideline"},
			expected: []matchExpectation{
				{"Coding Rule", true},
				{"Go Rule", false},            // Missing 'go' in query
				{"Documentation Rule", false}, // Missing 'documentation'
				{"Root Rule", true},           // Empty scope is subset of any
			},
		},
		{
			// Query [coding, go, guideline]
			// Coding Rule (scope {coding}) -> Subset of {coding, go} -> Match
			// Go Rule (scope {coding, go}) -> Subset of {coding, go} -> Match
			name:     "Query [coding, go, guideline]",
			keywords: []string{"coding", "go", "guideline"},
			expected: []matchExpectation{
				{"Coding Rule", true},
				{"Go Rule", true},
				{"Documentation Rule", false},
				{"Root Rule", true},
			},
		},
		{
			// Query [guideline] (Implicit scopes: {})
			// Coding Rule -> {coding} NOT subset of {} -> Exclude
			// Go Rule -> {coding, go} NOT subset of {} -> Exclude
			// Doc Rule -> {documentation} NOT subset of {} -> Exclude
			// Root Rule -> {} matches {} -> Match
			name:     "Query [guideline]",
			keywords: []string{"guideline"},
			expected: []matchExpectation{
				{"Coding Rule", false},
				{"Go Rule", false},
				{"Documentation Rule", false},
				{"Root Rule", true},
			},
		},
		{
			// Query [go] (Implicit scopes: {go})
			// Coding Rule -> {coding} NOT subset of {go} -> Exclude
			// Go Rule -> {coding, go} NOT subset of {go} -> Exclude (Missing 'coding')
			name:     "Query [go]",
			keywords: []string{"go"},
			expected: []matchExpectation{
				{"Coding Rule", false},
				{"Go Rule", false},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := idx.Search(tt.keywords, nil, false)

			for _, exp := range tt.expected {
				found := false
				for _, d := range got {
					if d.Title == exp.Title {
						found = true
						break
					}
				}
				if found != exp.ShouldMatch {
					t.Errorf("Title %q: match = %v, want %v", exp.Title, found, exp.ShouldMatch)
				}
			}
		})
	}
}
