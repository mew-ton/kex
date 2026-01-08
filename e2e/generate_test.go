package e2e

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/mew-ton/kex/internal/infrastructure/fs"
)

func TestKexGenerate_Success(t *testing.T) {
	t.Run("it should generate dist directory with kex.json and mirrored contents", func(t *testing.T) {
		tempDir := t.TempDir()

		// Setup Project
		contentsDir := filepath.Join(tempDir, "contents")
		os.MkdirAll(filepath.Join(contentsDir, "coding"), 0755)

		// Create Documents
		doc1 := `---
id: doc-1
title: Doc 1
status: adopted
---
Content 1`
		os.WriteFile(filepath.Join(contentsDir, "coding", "doc1.md"), []byte(doc1), 0644)

		docDraft := `---
id: doc-draft
title: Draft
status: draft
---
Draft Content`
		os.WriteFile(filepath.Join(contentsDir, "coding", "draft.md"), []byte(docDraft), 0644)

		os.WriteFile(filepath.Join(tempDir, ".kex.yaml"), []byte("root: contents\n"), 0644)

		// Run Generate
		cmd := exec.Command(kexBinary, "generate", tempDir)
		cmd.Dir = tempDir
		output, err := cmd.CombinedOutput()

		if err != nil {
			t.Fatalf("Generate failed: %v\nOutput: %s", err, output)
		}

		// Verify dist existence
		distDir := filepath.Join(tempDir, "dist")
		if _, err := os.Stat(distDir); os.IsNotExist(err) {
			t.Fatalf("dist directory not created")
		}

		// Verify kex.json
		kexJSONPath := filepath.Join(distDir, "kex.json")
		data, err := os.ReadFile(kexJSONPath)
		if err != nil {
			t.Fatalf("failed to read kex.json: %v", err)
		}

		var schema fs.IndexSchema
		if err := json.Unmarshal(data, &schema); err != nil {
			t.Fatalf("failed to parse kex.json: %v", err)
		}

		// Check Documents count (should exclude draft)
		if len(schema.Documents) != 1 {
			t.Errorf("Expected 1 document, got %d", len(schema.Documents))
		}
		// Dynamic ID: coding/doc1.md -> coding.doc1
		if schema.Documents[0].ID != "coding.doc1" {
			t.Errorf("Expected coding.doc1, got %s", schema.Documents[0].ID)
		}
		// Check Path relative to root or dist?
		// Our logic keeps it relative to root (e.g. "contents/coding/doc1.md") or relative to contents?
		// Logic: rel from cfg.Root. No, logic says doc.Path from doc.Path which is relative to... wait.
		// indexer.Loader -> ParseDocument -> doc.Path is Absolute by default.
		// Indexer.Export -> filepath.Rel(i.Root, doc.Path).
		// So path should be "coding/doc1.md" if i.Root is ".../contents".
		expectedPath := "coding/doc1.md"
		if schema.Documents[0].Path != expectedPath {
			t.Errorf("Expected path %s, got %s", expectedPath, schema.Documents[0].Path)
		}

		// Verify File Mirroring
		// Dist layout: dist/coding/doc1.md
		mirroredPath := filepath.Join(distDir, expectedPath)
		if _, err := os.Stat(mirroredPath); os.IsNotExist(err) {
			t.Errorf("Mirrored file not found at %s", mirroredPath)
		}

		// Verify Draft not mirrored?
		// Logic iterates schema documents. Since schema only has adopted, draft shouldn't be copied.
		draftPath := filepath.Join(distDir, "coding", "draft.md")
		if _, err := os.Stat(draftPath); !os.IsNotExist(err) {
			t.Errorf("Draft file should not be mirrored")
		}
	})
}

func TestKexGenerate_WithBaseURL(t *testing.T) {
	t.Run("it should prepend BaseURL to paths in kex.json", func(t *testing.T) {
		tempDir := t.TempDir()
		contentsDir := filepath.Join(tempDir, "contents")
		os.MkdirAll(contentsDir, 0755)

		doc := `---
title: Doc Base
status: adopted
---
Content`
		os.WriteFile(filepath.Join(contentsDir, "doc.md"), []byte(doc), 0644)

		config := "root: contents\nbaseURL: https://example.com/docs/"
		os.WriteFile(filepath.Join(tempDir, ".kex.yaml"), []byte(config), 0644)

		cmd := exec.Command(kexBinary, "generate", tempDir)
		cmd.Dir = tempDir
		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("Generate failed: %v\nOutput: %s", err, output)
		}

		kexJSONPath := filepath.Join(tempDir, "dist", "kex.json")
		data, _ := os.ReadFile(kexJSONPath)
		var schema fs.IndexSchema
		json.Unmarshal(data, &schema)

		expectedURL := "https://example.com/docs/doc.md"
		if schema.Documents[0].Path != expectedURL {
			t.Errorf("Expected URL %s, got %s", expectedURL, schema.Documents[0].Path)
		}
	})
}
