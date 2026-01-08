package domain

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseDocumentContent(t *testing.T) {
	tests := []struct {
		name    string
		content string
		wantID  string
		wantErr bool
	}{
		{
			name: "it should parse valid documents correctly",
			content: `---
title: Test Document
description: This is a test
keywords: [test, doc]
---
# Header
Body content`,
			wantID:  "test", // Derived from filename "test.md"
			wantErr: false,
		},
		{
			name:    "it should fail when frontmatter is missing",
			content: `# Just markdown`,
			wantID:  "",
			wantErr: true,
		},
		{
			name: "it should fail when title is missing",
			content: `---
description: No Title
---
Body`,
			wantID:  "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Use common root for this test case
			rootDir := t.TempDir()

			// Write to temp file
			tmpfile := filepath.Join(rootDir, "test.md")
			if err := os.WriteFile(tmpfile, []byte(tt.content), 0644); err != nil {
				t.Fatal(err)
			}

			doc, err := ParseDocument(tmpfile, rootDir)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseDocument() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && doc.ID != tt.wantID {
				t.Errorf("ParseDocument() ID = %v, want %v", doc.ID, tt.wantID)
			}
		})
	}

	t.Run("it should derive ID from scope", func(t *testing.T) {
		root := t.TempDir()
		// Structure: root/coding/go/style.md -> id: coding.go.style
		dir := filepath.Join(root, "coding", "go")
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatal(err)
		}

		content := `---
title: Style Guide
---
Body`
		path := filepath.Join(dir, "style.md")
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			t.Fatal(err)
		}

		doc, err := ParseDocument(path, root)
		if err != nil {
			t.Fatalf("ParseDocument failed: %v", err)
		}

		wantID := "coding.go.style"
		if doc.ID != wantID {
			t.Errorf("Expected ID %s, got %s", wantID, doc.ID)
		}

		// Verify Scopes
		if len(doc.Scopes) != 2 || doc.Scopes[0] != "coding" || doc.Scopes[1] != "go" {
			t.Errorf("Scopes wrong: %v", doc.Scopes)
		}
	})
}
