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

		cmd := exec.Command(kexBinary, "init", "--agents=antigravity", "--scopes=coding", "--scopes=documentation")
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

		// Verify rule creation (default Antigravity)
		expectedRule := filepath.Join(tempDir, ".agent/rules/kex.use-kex.md")
		if _, err := os.Stat(expectedRule); os.IsNotExist(err) {
			t.Errorf(".agent/rules/kex.use-kex.md was not created")
		}

		// Verify .claude/rules/kex/use-kex.md does NOT exist (unselected)
		if _, err := os.Stat(filepath.Join(tempDir, ".claude", "rules", "kex", "kex.use-kex.md")); !os.IsNotExist(err) {
			t.Error(".claude/rules/kex/kex.use-kex.md was created but should have been ignored")
		}
	})
}
