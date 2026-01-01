package e2e

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestKexInit(t *testing.T) {
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

	// NOTE: We rely on exit code (err == nil) and file existence for success.
	// Avoid asserting exact stdout text to prevent flakiness.
}
