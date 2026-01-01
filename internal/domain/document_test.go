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
id: test-doc
title: Test Document
description: This is a test
keywords: [test, doc]
---
# Header
Body content`,
			wantID:  "test-doc",
			wantErr: false,
		},
		{
			name:    "it should fail when frontmatter is missing",
			content: `# Just markdown`,
			wantID:  "",
			wantErr: true,
		},
		{
			name: "it should fail when the ID is missing",
			content: `---
title: No ID
---
Body`,
			wantID:  "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Write to temp file
			tmpfile := filepath.Join(t.TempDir(), "test.md")
			if err := os.WriteFile(tmpfile, []byte(tt.content), 0644); err != nil {
				t.Fatal(err)
			}

			doc, err := ParseDocument(tmpfile)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseDocument() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && doc.ID != tt.wantID {
				t.Errorf("ParseDocument() ID = %v, want %v", doc.ID, tt.wantID)
			}
		})
	}
}
