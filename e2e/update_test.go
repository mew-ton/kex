package e2e

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestKexUpdate(t *testing.T) {
	// Setup
	dir := t.TempDir()

	// Helper to run kex
	runKex := func(args ...string) {
		cmd := exec.Command(kexBinary, args...)
		cmd.Dir = dir
		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("kex %v failed: %v\nOutput: %s", args, err, output)
		}
	}

	// 1. Init
	runKex("init", "--agent-type=general")

	// 2. Modify a system document (should be overwritten)
	targetDoc := filepath.Join(dir, "contents/documentation/kex/choose-effective-keywords.md")
	err := os.WriteFile(targetDoc, []byte("Modified Content"), 0644)
	if err != nil {
		t.Fatalf("Failed to modify file: %v", err)
	}

	// 3. Create a dummy AGENTS.md with markers
	agentsPath := filepath.Join(dir, "AGENTS.md")
	initialAgentsContent := `
# My Custom Rules
Do whatever.

<!-- kex: auto-update start -->
Old Rules
<!-- kex: auto-update end -->
`
	err = os.WriteFile(agentsPath, []byte(initialAgentsContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write AGENTS.md: %v", err)
	}

	// 4. Run Update
	runKex("update")

	// 5. Verify system document is reverted (Overwritten)
	content, err := os.ReadFile(targetDoc)
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}
	if string(content) == "Modified Content" {
		t.Errorf("System document was NOT overwritten/updated")
	}
	// The template content changes based on source, but we check if it's NOT modified content.
	// And ideally it contains something from the template.

	// 6. Verify AGENTS.md is updated (Marker logic)
	agentsContent, err := os.ReadFile(agentsPath)
	if err != nil {
		t.Fatalf("Failed to read AGENTS.md: %v", err)
	}
	agentsStr := string(agentsContent)

	// Custom content should be preserved
	if !strings.Contains(agentsStr, "# My Custom Rules") {
		t.Errorf("User content in AGENTS.md was lost")
	}

	// Marker content should be updated (Old Rules replaced)
	if strings.Contains(agentsStr, "Old Rules") {
		t.Errorf("Content between markers was NOT updated")
	}

	// Check for new content (Default Scopes: Coding + Documentation)
	if !strings.Contains(agentsStr, "Design & Implementation Phase") {
		t.Errorf("Coding guidelines were not injected into AGENTS.md")
	}
	if !strings.Contains(agentsStr, "Documentation Phase") {
		t.Errorf("Documentation guidelines were not injected into AGENTS.md")
	}
	// Header is outside markers, so it should NOT be injected if not present

	// 7. Verify Configural Update (Change Scopes)
	// Modify .kex.yaml to only have "coding" scope
	configData := `root: contents
agent:
  type: general
  scopes: ["coding"]
`
	if err := os.WriteFile(filepath.Join(dir, ".kex.yaml"), []byte(configData), 0644); err != nil {
		t.Fatal(err)
	}

	// Run Update again
	runKex("update")

	// Verify AGENTS.md again
	agentsContent, _ = os.ReadFile(agentsPath)
	agentsStr = string(agentsContent)

	if !strings.Contains(agentsStr, "Design & Implementation Phase") {
		t.Errorf("Coding guidelines should still be present")
	}
	if strings.Contains(agentsStr, "Documentation Phase") {
		t.Errorf("Documentation guidelines should have been removed based on config")
	}
}
