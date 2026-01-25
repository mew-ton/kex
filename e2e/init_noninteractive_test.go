package e2e

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestInitNonInteractiveSkills(t *testing.T) {
	// Setup
	dir := t.TempDir()

	// Run kex init with both agents and skills
	// This should enable 'cursor' for skills, NOT just 'claude'
	cmd := exec.Command(kexBinary, "init", "--agents=cursor", "--skills=go")
	cmd.Dir = dir
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("kex init failed: %v\nOutput: %s", err, output)
	}

	// Read .kex.yaml
	configPath := filepath.Join(dir, ".kex.yaml")
	content, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("Failed to read .kex.yaml: %v", err)
	}
	contentStr := string(content)

	// Verify cursor is in skills targets
	// We expect cursor to be present in the file (simple check for now)
	if !strings.Contains(contentStr, "cursor") {
		t.Errorf("Expected cursor to be present in config, got:\n%s", contentStr)
	}

	// Verify backup behavior (no agents specified)
	// Run in a new dir
	dir2 := t.TempDir()
	cmd2 := exec.Command(kexBinary, "init", "--skills=go")
	cmd2.Dir = dir2
	output2, err := cmd2.CombinedOutput()
	if err != nil {
		t.Fatalf("kex init failed: %v\nOutput: %s", err, output2)
	}

	configPath2 := filepath.Join(dir2, ".kex.yaml")
	content2, err := os.ReadFile(configPath2)
	if err != nil {
		t.Fatalf("Failed to read .kex.yaml: %v", err)
	}
	contentStr2 := string(content2)

	// Verify claude is default
	if !strings.Contains(contentStr2, "claude") {
		t.Errorf("Expected claude to be present in config for skills fallback, got:\n%s", contentStr2)
	}
}
