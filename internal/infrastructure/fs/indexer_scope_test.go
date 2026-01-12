package fs

import (
	"testing"
)

func TestIndexer_isSubset(t *testing.T) {
	tests := []struct {
		name        string
		docScopes   []string
		queryScopes map[string]struct{}
		want        bool
	}{
		{
			name:      "Exact Match",
			docScopes: []string{"coding", "go"},
			queryScopes: map[string]struct{}{
				"coding": {},
				"go":     {},
			},
			want: true,
		},
		{
			name:      "Superset Query (Match)",
			docScopes: []string{"coding"},
			queryScopes: map[string]struct{}{
				"coding":    {},
				"guideline": {},
			},
			want: true,
		},
		{
			name:      "Subset Query (Exclude - Missing 1 scope)",
			docScopes: []string{"coding", "go"},
			queryScopes: map[string]struct{}{
				"coding": {},
			},
			want: false,
		},
		{
			name:      "Disjoint (Exclude)",
			docScopes: []string{"documentation"},
			queryScopes: map[string]struct{}{
				"coding": {},
			},
			want: false,
		},
		{
			name:      "Empty Doc Scopes (Match - Root doc)",
			docScopes: []string{},
			queryScopes: map[string]struct{}{
				"coding": {},
			},
			want: true,
		},
		{
			name:      "Case Insensitive Handling Check (Should be handled by caller, but good to verify assumption if logic changes)",
			docScopes: []string{"Coding"}, // Caller standardizes strictly, but if logic compares strings...
			queryScopes: map[string]struct{}{
				"coding": {}, // Assuming map keys are lowercase
			},
			// isSubset usually assumes normalized input?
			// In our implementation plan, keys are normalized. Doc scopes in domain might be raw?
			// Let's assume isSubset normalizes doc scopes or expects them normalized?
			// Indexer.Load normalizes? Actually Indexer stores exact strings.
			// So isSubset should probably normalize doc scopes to be safe.
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a dummy indexer to access the method, or make it a helper function.
			// Since we want to test the logic, we can make it a private method on Indexer or just a helper.
			// For now, let's implement the logic here directly to verify expectations,
			// or assume we will act on the Indexer method.
			// Let's assume we implement `isSubset` as a helper or method.
			// We can export it for test or put test in `fs` package (which it is).

			got := isSubset(tt.docScopes, tt.queryScopes)
			if got != tt.want {
				t.Errorf("isSubset() = %v, want %v", got, tt.want)
			}
		})
	}
}
