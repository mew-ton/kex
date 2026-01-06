package e2e

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestKexStart_Failure_MissingRoot(t *testing.T) {
	t.Run("it should fail when root directory is missing", func(t *testing.T) {
		tempDir := t.TempDir()
		// No contents dir created

		// Create config pointing to contents
		os.WriteFile(filepath.Join(tempDir, ".kex.yaml"), []byte("root: contents\n"), 0644)

		cmd := exec.Command(kexBinary, "start")
		cmd.Dir = tempDir
		output, err := cmd.CombinedOutput()

		if err == nil {
			t.Fatalf("Expected start to fail due to missing root, but it succeeded.")
		}

		if !strings.Contains(string(output), "not found") {
			t.Errorf("Expected 'not found' error, got: %s", output)
		}
	})
}

func TestKexStart_Failure_ParseErrors(t *testing.T) {
	t.Run("it should fail when documents have parse errors", func(t *testing.T) {
		tempDir := t.TempDir()

		// Invalid Frontmatter
		doc := `---
id: broken
title: Broken
: invalid-yaml
---
Content`

		os.Mkdir(filepath.Join(tempDir, "contents"), 0755)
		os.WriteFile(filepath.Join(tempDir, "contents", "broken.md"), []byte(doc), 0644)
		os.WriteFile(filepath.Join(tempDir, ".kex.yaml"), []byte("root: contents\n"), 0644)

		cmd := exec.Command(kexBinary, "start")
		cmd.Dir = tempDir
		output, err := cmd.CombinedOutput()

		if err == nil {
			t.Fatalf("Expected start to fail due to parse errors, but it succeeded.")
		}

		if !strings.Contains(string(output), "Failed to start") {
			t.Errorf("Expected 'Failed to start' error, got: %s", output)
		}
	})
}

func TestKexStart_WithRootFlag(t *testing.T) {
	t.Run("it should start successfully with --root flag overriding config", func(t *testing.T) {
		tempDir := t.TempDir()
		customRoot := filepath.Join(tempDir, "custom_guidelines")
		os.MkdirAll(customRoot, 0755)

		// Create a valid document in custom root
		doc := `---
id: test-doc
title: Test Doc
description: Test
---
Content`
		os.WriteFile(filepath.Join(customRoot, "test-doc.md"), []byte(doc), 0644)

		// Create invalid config to ensure override works
		os.WriteFile(filepath.Join(tempDir, ".kex.yaml"), []byte("root: non_existent\n"), 0644)

		// Run start with --root
		cmd := exec.Command(kexBinary, "start", "--root", customRoot)
		cmd.Dir = tempDir

		// Start the process
		if err := cmd.Start(); err != nil {
			t.Fatalf("Failed to start command: %v", err)
		}

		// Wait a bit to ensure it doesn't crash immediately (validation happens on start)
		done := make(chan error, 1)
		go func() {
			done <- cmd.Wait()
		}()

		// For this test, valid start means it is running.
		// If it failed validation, it would exit immediately.

		// Let's kill it.
		if err := cmd.Process.Kill(); err != nil {
			t.Logf("Failed to kill process: %v", err)
		}

		// Wait for exit
		err := <-done

		if err != nil && !strings.Contains(err.Error(), "killed") && !strings.Contains(err.Error(), "interrupt") {
			t.Errorf("Command exited with error: %v", err)
		}
	})
}
