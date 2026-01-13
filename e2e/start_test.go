package e2e

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestKexStart_Failure_MissingRoot(t *testing.T) {
	t.Run("it should fail when source directory is missing", func(t *testing.T) {
		tempDir := t.TempDir()
		// No contents dir created

		// Create config pointing to contents
		os.WriteFile(filepath.Join(tempDir, ".kex.yaml"), []byte("source: contents\n"), 0644)

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
		os.WriteFile(filepath.Join(tempDir, ".kex.yaml"), []byte("source: contents\n"), 0644)

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

func TestKexStart_References(t *testing.T) {
	t.Run("it should start successfully with references", func(t *testing.T) {
		tempDir := t.TempDir()

		// Project 1
		proj1 := filepath.Join(tempDir, "proj1")
		os.MkdirAll(filepath.Join(proj1, "docs"), 0755)
		os.WriteFile(filepath.Join(proj1, ".kex.yaml"), []byte("source: docs\n"), 0644)
		os.WriteFile(filepath.Join(proj1, "docs", "doc1.md"), []byte("---\nid: doc1\ntitle: Doc 1\n---\n"), 0644)

		// Project 2
		proj2 := filepath.Join(tempDir, "proj2")
		os.MkdirAll(filepath.Join(proj2, "notes"), 0755)
		os.WriteFile(filepath.Join(proj2, ".kex.yaml"), []byte("source: notes\n"), 0644)
		os.WriteFile(filepath.Join(proj2, "notes", "doc2.md"), []byte("---\nid: doc2\ntitle: Doc 2\n---\n"), 0644)

		// Main Config with references
		configContent := `references:
  - proj1
  - proj2
`
		os.WriteFile(filepath.Join(tempDir, ".kex.yaml"), []byte(configContent), 0644)

		// Run start without args
		cmd := exec.Command(kexBinary, "start")
		cmd.Dir = tempDir

		// Start the process
		if err := cmd.Start(); err != nil {
			t.Fatalf("Failed to start command: %v", err)
		}

		// Cleanup
		defer func() {
			if cmd.Process != nil {
				cmd.Process.Kill()
			}
		}()

		done := make(chan error, 1)
		go func() {
			done <- cmd.Wait()
		}()

		// Wait briefly to assume validation passed
		select {
		case err := <-done:
			if err != nil && !strings.Contains(err.Error(), "killed") {
				t.Errorf("Command exited unexpectedly: %v", err)
			}
		case <-time.After(500 * time.Millisecond):
			// Success
			cmd.Process.Kill()
		}
	})
}

func TestKexStart_LogFile(t *testing.T) {
	t.Run("it should create log file and record startup stats", func(t *testing.T) {
		tempDir := t.TempDir()
		contentsDir := filepath.Join(tempDir, "contents")
		os.MkdirAll(contentsDir, 0755)

		// Create a valid document
		doc := `---
id: log-doc
title: Log Doc
description: Test
---
Content`
		os.WriteFile(filepath.Join(contentsDir, "doc.md"), []byte(doc), 0644)
		os.WriteFile(filepath.Join(tempDir, ".kex.yaml"), []byte("source: contents\n"), 0644)

		logFile := filepath.Join(tempDir, "server.log")

		cmd := exec.Command(kexBinary, "start", "--log-file", logFile)
		cmd.Dir = tempDir

		if err := cmd.Start(); err != nil {
			t.Fatalf("Failed to start command: %v", err)
		}

		// Cleanup
		defer func() {
			if cmd.Process != nil {
				cmd.Process.Kill()
			}
		}()

		// Wait for logs to be written
		success := false
		for i := 0; i < 20; i++ {
			if _, err := os.Stat(logFile); err == nil {
				content, _ := os.ReadFile(logFile)
				if strings.Contains(string(content), "Kex Server Starting...") &&
					strings.Contains(string(content), "Documents Loaded: 1") {
					success = true
					break
				}
			}
			time.Sleep(100 * time.Millisecond)
		}

		if !success {
			// Read file content for debug
			content, _ := os.ReadFile(logFile)
			t.Errorf("Expected log file to contain startup stats. Content:\n%s", content)
		}
	})
}
