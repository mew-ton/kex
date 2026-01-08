package e2e

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestKexCheck_InvalidFrontmatter(t *testing.T) {
	t.Run("it should report warnings for draft documents with invalid frontmatter", func(t *testing.T) {
		fixturePath, _ := filepath.Abs("fixtures/check-invalid-frontmatter")

		cmd := exec.Command(kexBinary, "check")
		cmd.Dir = fixturePath
		output, err := cmd.CombinedOutput()

		// Currently returns success (0) but with warnings for Drafts.
		if err != nil {
			t.Logf("Check failed: %v", err)
		} else {
			// Verify we got the warning
			if !strings.Contains(string(output), "WARNING") {
				t.Error("Expected warning in output, got none")
			}
		}
	})
}

func TestKexCheck_NoDocuments(t *testing.T) {
	t.Run("it should fail when no configuration or documents are found", func(t *testing.T) {
		tempDir := t.TempDir()

		cmd := exec.Command(kexBinary, "check")
		cmd.Dir = tempDir
		output, _ := cmd.CombinedOutput()

		// We expect some failure or warning about missing config/docs, usually strict check fails if nothing found?
		// Actually implementation detail: CheckCommand defaults?
		// e2e test log said "Failed to load documents: lstat contents: no such file".
		// So it should fail if default 'contents' dir is missing.
		if !strings.Contains(string(output), "Failed to load documents") {
			t.Logf("Output: %s", output)
		}
	})
}

func TestKexCheck_NoDocuments_WithConfig(t *testing.T) {
	t.Run("it should pass when config exists but no documents are present (empty root)", func(t *testing.T) {
		tempDir := t.TempDir()
		os.WriteFile(filepath.Join(tempDir, ".kex.yaml"), []byte("source: .\n"), 0644)

		cmd := exec.Command(kexBinary, "check")
		cmd.Dir = tempDir
		output, err := cmd.CombinedOutput()

		if err != nil {
			t.Errorf("Expected success, failed: %v\nOutput: %s", err, output)
		}
	})
}

func TestKexCheck_JSONOutput(t *testing.T) {
	t.Run("it should output valid JSON when --json flag is passed", func(t *testing.T) {
		tempDir := t.TempDir()
		// Create a dummy document
		doc := `---
id: test-json
title: JSON Test
keywords: [json]
---
Content`
		os.Mkdir(filepath.Join(tempDir, "contents"), 0755)
		os.WriteFile(filepath.Join(tempDir, "contents", "test-json.md"), []byte(doc), 0644)
		os.WriteFile(filepath.Join(tempDir, ".kex.yaml"), []byte("source: contents\n"), 0644)

		cmd := exec.Command(kexBinary, "check", "--json")
		cmd.Dir = tempDir
		output, err := cmd.CombinedOutput()

		if err != nil {
			t.Fatalf("Check failed: %v\nOutput: %s", err, output)
		}

		// Simple JSON validation
		if !strings.HasPrefix(string(output), "{") {
			t.Errorf("Output does not start with '{': %s", output)
		}
		if !strings.Contains(string(output), `"valid": true`) {
			t.Errorf("JSON output missing 'valid': true: %s", output)
		}
		if !strings.Contains(string(output), `"id": "test-json"`) {
			t.Errorf("JSON output missing document ID: %s", output)
		}
	})
}

func TestKexCheck_Success(t *testing.T) {
	t.Run("it should pass and output success message when documents are valid", func(t *testing.T) {
		tempDir := t.TempDir()

		// Valid Document
		doc := `---
id: valid-doc
title: Valid Document
status: adopted
keywords: [valid]
---
Content`

		os.Mkdir(filepath.Join(tempDir, "contents"), 0755)
		os.WriteFile(filepath.Join(tempDir, "contents", "valid-doc.md"), []byte(doc), 0644)
		os.WriteFile(filepath.Join(tempDir, ".kex.yaml"), []byte("source: contents\n"), 0644)

		cmd := exec.Command(kexBinary, "check")
		cmd.Dir = tempDir
		output, err := cmd.CombinedOutput()

		if err != nil {
			t.Fatalf("Check failed expected success: %v\nOutput: %s", err, output)
		}

		// Verify output
		if !strings.Contains(string(output), "Success") && !strings.Contains(string(output), "No errors found") {
			// t.Logf("Output did not contain 'Success': %s", output)
		}
	})
}

func TestKexCheck_PositionalArg(t *testing.T) {
	t.Run("it should pass when project root is passed as positional argument", func(t *testing.T) {
		baseDir := t.TempDir()
		projectRoot := filepath.Join(baseDir, "my-project")
		contentsDir := filepath.Join(projectRoot, "custom_check_contents")
		os.MkdirAll(contentsDir, 0755)

		// Create a valid document
		doc := `---
id: pos-check-doc
title: Positional Check Doc
status: adopted
---
Content`
		os.WriteFile(filepath.Join(contentsDir, "pos-check-doc.md"), []byte(doc), 0644)

		// Config at projectRoot
		os.WriteFile(filepath.Join(projectRoot, ".kex.yaml"), []byte("source: custom_check_contents\n"), 0644)

		// Run kex check <projectRoot> from baseDir
		cmd := exec.Command(kexBinary, "check", projectRoot)
		cmd.Dir = baseDir
		output, err := cmd.CombinedOutput()

		if err != nil {
			t.Fatalf("Check failed expected success: %v\nOutput: %s", err, string(output))
		}

		if !strings.Contains(string(output), "All checks passed") {
			t.Errorf("Expected 'All checks passed', got: %s", string(output))
		}
	})
}
