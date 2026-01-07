package e2e

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/mew-ton/kex/internal/infrastructure/fs"
)

func TestKexStart_Remote(t *testing.T) {
	// Setup Mock Server
	// We need to serve a kex.json and a markdown file
	mux := http.NewServeMux()

	docData := "---" + "\n" +
		"id: remote-doc" + "\n" +
		"title: Remote Doc" + "\n" +
		"description: A valid remote document" + "\n" +
		"keywords: [remote]" + "\n" +
		"status: adopted" + "\n" +
		"---" + "\n" +
		"Remote Content"

	mux.HandleFunc("/contents/remote-doc.md", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/markdown")
		fmt.Fprint(w, docData)
	})

	mux.HandleFunc("/kex.json", func(w http.ResponseWriter, r *http.Request) {
		schema := fs.IndexSchema{
			GeneratedAt: time.Now(),
			Documents: []*fs.DocumentSchema{
				{
					ID:          "remote-doc",
					Title:       "Remote Doc",
					Description: "A valid remote document",
					Keywords:    []string{"remote"},
					Path:        "contents/remote-doc.md", // Relative to kex.json
				},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(schema)
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	t.Run("it should start with remote URL", func(t *testing.T) {
		// Run kex start <url>
		cmd := exec.Command(kexBinary, "start", server.URL)

		// Wait, start is long running. CombinedOutput waits for exit.
		// If start runs successfully, it blocks forever.
		// We need to capture pipes.

		// Revert to Start/Wait but capture pipes? Or just log file?
		// StdoutPipe/StderrPipe.
		/*
			stdout, _ := cmd.StdoutPipe()
			stderr, _ := cmd.StderrPipe()
			if err := cmd.Start(); err != nil { ... }
			// Read from pipes in goroutine?
		*/

		// For now, let's just inspect output if it exits early.
		// If it blocks, we kill it.
		// Use a buffer?
		// Simpler: use file for output?

		outfile, _ := os.Create(filepath.Join(t.TempDir(), "output.log"))
		defer outfile.Close()
		cmd.Stdout = outfile
		cmd.Stderr = outfile

		// Start process
		if err := cmd.Start(); err != nil {
			t.Fatalf("Failed to start: %v", err)
		}
		defer cmd.Process.Kill()

		// Wait briefly for validation
		done := make(chan error, 1)
		go func() {
			done <- cmd.Wait()
		}()

		// Give it time to fetch and validate
		time.Sleep(500 * time.Millisecond)

		if err := cmd.Process.Kill(); err != nil {
			t.Logf("Kill error: %v", err)
		}

		err := <-done
		if err != nil && !strings.Contains(err.Error(), "killed") && !strings.Contains(err.Error(), "interrupt") {
			// Check if it exited with validation error?
			debugOut, _ := os.ReadFile(outfile.Name())
			t.Errorf("Remote start failed: %v\nOutput:\n%s", err, string(debugOut))
		}
	})

	// TODO: Verify that it actually loaded the document?
	// Hard to check stdio without an MCP client simulation.
	// But if validation failed (fetch error), it would exit non-zero.
}
