package e2e

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/mew-ton/kex/internal/infrastructure/config"
	"gopkg.in/yaml.v3"
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

		if !strings.Contains(string(output), "failed to start") {
			t.Errorf("Expected 'failed to start' error, got: %s", output)
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

func TestKexStart_InvalidReference_ShouldNotFailIfDocsExist(t *testing.T) {
	// Setup
	tmpDir, err := os.MkdirTemp("", "kex-test-start-invalid-ref")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create valid content
	contentDir := filepath.Join(tmpDir, "contents")
	os.Mkdir(contentDir, 0755)
	os.WriteFile(filepath.Join(contentDir, "doc.md"), []byte("---\ntitle: Test Doc\n---\nHello"), 0644)

	// Create config with valid source AND invalid reference
	cfg := config.Config{
		Source: "contents",
		References: []string{
			"non-existent-path",
		},
		Logging: config.Logging{
			Level: "debug",
		},
	}
	data, _ := yaml.Marshal(cfg)
	os.WriteFile(filepath.Join(tmpDir, ".kex.yaml"), data, 0644)

	// Run kex start
	cmd := exec.Command(kexBinary, "start")
	cmd.Dir = tmpDir

	// Run asynchronously
	if err := cmd.Start(); err != nil {
		t.Fatalf("Failed to start command: %v", err)
	}

	// Wait a bit to ensure it keeps running (warning only)
	time.Sleep(1 * time.Second)

	if cmd.ProcessState != nil && cmd.ProcessState.Exited() {
		output, _ := cmd.CombinedOutput()
		t.Fatalf("Command exited unexpectedly! Output: %s", output)
	}

	// Cleanup
	cmd.Process.Kill()
}

func TestKexStart_AllSourcesInvalid_ShouldFail(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "kex-test-start-all-invalid")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Config with invalid source AND invalid reference
	cfg := config.Config{
		Source:     "invalid-source",
		References: []string{"invalid-ref"},
	}
	data, _ := yaml.Marshal(cfg)
	os.WriteFile(filepath.Join(tmpDir, ".kex.yaml"), data, 0644)

	cmd := exec.Command(kexBinary, "start")
	cmd.Dir = tmpDir

	// It should exit quickly with error
	if err := cmd.Start(); err != nil {
		t.Fatalf("Failed to start command: %v", err)
	}

	done := make(chan error, 1)
	go func() {
		done <- cmd.Wait()
	}()

	select {
	case err := <-done:
		if err == nil {
			t.Fatal("Command exited successfully (0) but expected failure due to no docs")
		}
		// Expected error exit
	case <-time.After(2 * time.Second):
		cmd.Process.Kill()
		t.Fatal("Command stuck running! Expected failure due to no docs")
	}
}

func TestKexStart_CwdFlag(t *testing.T) {
	t.Run("it should start successfully when using --cwd flag", func(t *testing.T) {
		tempDir := t.TempDir()
		projectDir := filepath.Join(tempDir, "project")
		os.MkdirAll(filepath.Join(projectDir, "docs"), 0755)

		// Create a valid configuration and document in a sub-folder
		os.WriteFile(filepath.Join(projectDir, ".kex.yaml"), []byte("source: docs\n"), 0644)
		os.WriteFile(filepath.Join(projectDir, "docs", "doc1.md"), []byte("---\nid: doc1\ntitle: Doc 1\n---\n"), 0644)

		// Run start from the parent tempDir (not inside projectDir)
		cmd := exec.Command(kexBinary, "start", "--cwd", projectDir)
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
		case <-time.After(1 * time.Second):
			// Success
			cmd.Process.Kill()
		}
	})
}

func TestKexStart_CLIReferences(t *testing.T) {
	t.Run("it should start successfully when using CLI args as references", func(t *testing.T) {
		tempDir := t.TempDir()

		// Project 1
		proj1 := filepath.Join(tempDir, "proj1")
		os.MkdirAll(filepath.Join(proj1, "docs"), 0755)
		// No .kex.yaml needed for purely reference-based if we point to docs dir directly
		os.WriteFile(filepath.Join(proj1, "docs", "doc1.md"), []byte("---\nid: doc1\ntitle: Doc 1\n---\n"), 0644)

		// Project 2
		proj2 := filepath.Join(tempDir, "proj2")
		os.MkdirAll(proj2, 0755)
		os.WriteFile(filepath.Join(proj2, "doc2.md"), []byte("---\nid: doc2\ntitle: Doc 2\n---\n"), 0644)

		// Run start with two arguments
		cmd := exec.Command(kexBinary, "start", filepath.Join(proj1, "docs"), proj2)
		cmd.Dir = tempDir

		if err := cmd.Start(); err != nil {
			t.Fatalf("Failed to start command: %v", err)
		}

		defer func() {
			if cmd.Process != nil {
				cmd.Process.Kill()
			}
		}()

		done := make(chan error, 1)
		go func() {
			done <- cmd.Wait()
		}()

		select {
		case err := <-done:
			if err != nil && !strings.Contains(err.Error(), "killed") {
				t.Errorf("Command exited unexpectedly: %v", err)
			}
		case <-time.After(1 * time.Second):
			cmd.Process.Kill()
		}
	})
}
