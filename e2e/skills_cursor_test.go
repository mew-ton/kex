package e2e

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestCursorSkillsGeneration(t *testing.T) {
	// Setup
	dir := t.TempDir()

	// Helper to run kex
	runKex := func(args ...string) {
		cmd := exec.Command(kexBinary, args...)
		cmd.Dir = dir
		output, err := cmd.CombinedOutput()
		t.Logf("kex %v Output:\n%s", args, output)
		if err != nil {
			t.Fatalf("kex %v failed: %v", args, err)
		}
	}

	// 1. Init (generates default content)
	runKex("init", "--agents=cursor", "--scopes=documentation")

	// 2. Configure .kex.yaml for Cursor skills
	cfgData := `source: contents
update:
  ai:
    targets: [cursor]
    keywords: [documentation, kex]
`
	err := os.WriteFile(filepath.Join(dir, ".kex.yaml"), []byte(cfgData), 0644)
	if err != nil {
		t.Fatalf("Failed to write config: %v", err)
	}

	// 3. Run Update
	runKex("update")

	// 4. Verify Generation
	// Expected Output: .cursor/skills/kex.documentation.kex.choose-effective-keywords/SKILL.md

	expectedPath := filepath.Join(dir, ".cursor/skills/kex.documentation.kex.choose-effective-keywords/SKILL.md")
	if _, err := os.Stat(expectedPath); os.IsNotExist(err) {
		// Debug: list what WAS generated
		filepath.WalkDir(filepath.Join(dir, ".cursor"), func(path string, d os.DirEntry, err error) error {
			t.Logf("Found: %s", path)
			return nil
		})
		t.Fatalf("Expected skill file not found at: %s", expectedPath)
	}

	// 5. Verify Content
	content, err := os.ReadFile(expectedPath)
	if err != nil {
		t.Fatalf("Failed to read skill file: %v", err)
	}
	contentStr := string(content)

	// Check Frontmatter (description)
	if !strings.Contains(contentStr, "description:") {
		t.Errorf("Content missing description. Got:\n%s", contentStr)
	}
}
