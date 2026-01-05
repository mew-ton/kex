package e2e

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestKexInit(t *testing.T) {
	t.Run("it should create default config and content structure", func(t *testing.T) {
		tempDir := t.TempDir()

		cmd := exec.Command(kexBinary, "init")
		cmd.Dir = tempDir
		output, err := cmd.CombinedOutput()

		if err != nil {
			t.Fatalf("kex init failed: %v\nOutput: %s", err, output)
		}

		// Verify .kex.yaml exists
		if _, err := os.Stat(filepath.Join(tempDir, ".kex.yaml")); os.IsNotExist(err) {
			t.Error(".kex.yaml was not created")
		}

		// Verify contents directory exists
		if _, err := os.Stat(filepath.Join(tempDir, "contents")); os.IsNotExist(err) {
			t.Error("contents directory was not created")
		}

		// Verify contents/documentation/kex/write-concise-content.md exists (actual template)
		if _, err := os.Stat(filepath.Join(tempDir, "contents", "documentation", "kex", "write-concise-content.md")); os.IsNotExist(err) {
			t.Error("contents/documentation/kex/write-concise-content.md was not extracted")
		}
	})
}
