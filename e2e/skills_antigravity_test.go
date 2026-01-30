package e2e

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestAntigravitySkillsGeneration(t *testing.T) {
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
	runKex("init", "--agents=antigravity", "--scopes=documentation")

	// Debug content
	debugPath := filepath.Join(dir, "contents", "documentation", "kex", "choose-effective-keywords.md")
	debugContent, _ := os.ReadFile(debugPath)
	t.Logf("Generated Document Content:\n%s", string(debugContent))

	// 2. Configure .kex.yaml for Antigravity skills
	// We want to target documents that exist. 'kex init' creates 'contents/documentation/kex/...'
	// So we use keywords "kex".
	cfgData := `source: contents
update:
  ai-skills:
    targets: [antigravity]
    keywords: [documentation, kex]
`
	err := os.WriteFile(filepath.Join(dir, ".kex.yaml"), []byte(cfgData), 0644)
	if err != nil {
		t.Fatalf("Failed to write config: %v", err)
	}

	// 3. Run Update
	runKex("update")

	// 4. Verify Generation
	// We expect a skill for "contents/documentation/kex/choose-effective-keywords.md"
	// ID: documentation/kex/choose-effective-keywords
	// Template: .agent/skills/{{.SkillName}}/SKILL.md.template
	// Expected Output: .agent/skills/kex.documentation.kex.choose-effective-keywords/SKILL.md

	expectedPath := filepath.Join(dir, ".agent/skills/kex.documentation.kex.choose-effective-keywords/SKILL.md")
	if _, err := os.Stat(expectedPath); os.IsNotExist(err) {
		// Debug: list what WAS generated
		filepath.WalkDir(filepath.Join(dir, ".agent"), func(path string, d os.DirEntry, err error) error {
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

	// Check Frontmatter
	if !strings.Contains(contentStr, "name: documentation.kex.choose-effective-keywords") {
		t.Errorf("Content missing correct name in frontmatter. Got:\n%s", contentStr)
	}

	// Check Description (from document)
	if !strings.Contains(contentStr, "description:") {
		t.Errorf("Content missing description. Got:\n%s", contentStr)
	}
}
