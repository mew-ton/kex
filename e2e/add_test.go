package e2e

import (
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestKexAdd(t *testing.T) {
	t.Run("it should add local reference to .kex.yaml", func(t *testing.T) {
		dir := t.TempDir()

		// Create a target directory to reference
		refDir := filepath.Join(dir, "my-library")
		os.MkdirAll(refDir, 0755)

		// Run kex add my-library
		cmd := exec.Command(kexBinary, "add", "my-library")
		cmd.Dir = dir
		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("kex add failed: %v\nOutput: %s", err, output)
		}

		// Verify .kex.yaml
		content, err := os.ReadFile(filepath.Join(dir, ".kex.yaml"))
		if err != nil {
			t.Fatalf("Failed to read .kex.yaml: %v", err)
		}

		if !strings.Contains(string(content), "references:") {
			t.Errorf(".kex.yaml missing references key")
		}
		if !strings.Contains(string(content), "- my-library") {
			t.Errorf(".kex.yaml missing my-library in references")
		}
	})

	t.Run("it should fail when adding invalid local path", func(t *testing.T) {
		dir := t.TempDir()

		cmd := exec.Command(kexBinary, "add", "non-existent")
		cmd.Dir = dir
		err := cmd.Run()
		if err == nil {
			t.Error("Expected kex add non-existent to fail, but it succeeded")
		}
	})

	t.Run("it should add remote URL reference to .kex.yaml", func(t *testing.T) {
		dir := t.TempDir()

		// Start a local test server to simulate a remote Kex source
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/kex.json" {
				w.WriteHeader(http.StatusOK)
				// Serve a minimal valid kex.json
				w.Write([]byte(`{"documents": []}`))
				return
			}
			w.WriteHeader(http.StatusNotFound)
		}))
		defer server.Close()

		url := server.URL

		cmd := exec.Command(kexBinary, "add", url)
		cmd.Dir = dir
		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("kex add failed: %v\nOutput: %s", err, output)
		} else {
			content, _ := os.ReadFile(filepath.Join(dir, ".kex.yaml"))
			if !strings.Contains(string(content), url) {
				t.Errorf(".kex.yaml missing url %s", url)
			}
		}
	})

	t.Run("it should fail when adding unreachable URL", func(t *testing.T) {
		dir := t.TempDir()

		url := "https://invalid.domain.example.com"

		cmd := exec.Command(kexBinary, "add", url)
		cmd.Dir = dir
		err := cmd.Run()
		if err == nil {
			t.Error("Expected kex add invalid-url to fail, but it succeeded")
		}
	})
}
