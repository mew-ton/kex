package e2e

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mew-ton/kex/assets"
	"github.com/mew-ton/kex/internal/infrastructure/config"
	"github.com/mew-ton/kex/internal/usecase/generator"
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
	runKex("init", "--agents=antigravity", "--scopes=coding", "--scopes=documentation")

	// 2. Modify a system document (should be overwritten)
	targetDoc := filepath.Join(dir, "contents/documentation/kex/choose-effective-keywords.md")
	err := os.WriteFile(targetDoc, []byte("Modified Content"), 0644)
	if err != nil {
		t.Fatalf("Failed to modify file: %v", err)
	}

	// 3. Setup user content in .agent/rules/use-kex.md
	rulePath := filepath.Join(dir, ".agent/rules/use-kex.md")

	initialContent, _ := os.ReadFile(rulePath)
	// Inject a custom header line
	customHeader := "# My Custom Rules\n" + string(initialContent)
	err = os.WriteFile(rulePath, []byte(customHeader), 0644)
	if err != nil {
		t.Fatalf("Failed to write use-kex.md: %v", err)
	}

	// 4. Run Update (Antigravity should overwrite by default based on kex init setting overwrite? No, wait.)
	// kex init generates .kex.yaml with "overwrite" for selected agents.
	// So if we run update, it should overwrite.

	runKex("update")

	// 5. Verify system document is overwritten (default)
	content, err := os.ReadFile(targetDoc)
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}
	if string(content) == "Modified Content" {
		t.Errorf("System document was NOT overwritten/updated")
	}

	// 6. Verify use-kex.md is OVERWRITTEN (Because kex init sets it to overwrite)

	contentOverwrittenInit, _ := os.ReadFile(rulePath)
	if strings.Contains(string(contentOverwrittenInit), "# My Custom Rules") {
		// It was preserved! Why?
		// Maybe init didn't set overwrite?
		// Check .kex.yaml content
		configFile, _ := os.ReadFile(filepath.Join(dir, ".kex.yaml"))
		t.Logf("Config Content: %s", string(configFile))
		t.Errorf("Rule file should have been overwritten by default init (overwrite strategy), but was preserved")
	}

	// 7. Test Agent Switch (Enable Claude)
	// We need to update .kex.yaml to enable Claude
	// Let's hack .kex.yaml

	// 8. Test Default Fallback (Ignore/Nil)
	// Inject custom content again to verify preservation
	customContentPreserved := "# Should Be Preserved\n"
	err = os.WriteFile(rulePath, []byte(customContentPreserved), 0644)
	if err != nil {
		t.Fatalf("Failed to write custom content: %v", err)
	}

	configFallback := config.UpdateConfig{
		Documents: make(map[string]string),
		// Empty Ai effective means no targets
	}
	// If logic is "Ignore", then Update() does nothing. File preserves state.

	gen := generator.New(assets.Assets)
	updateOpts := generator.UpdateOptions{
		Cwd: dir, // dir is valid here, references tempDir
	}
	err = gen.Update(updateOpts, configFallback)
	if err != nil {
		t.Fatalf("Update (Empty Config) failed: %v", err)
	}
	contentIgnored, _ := os.ReadFile(rulePath)
	if !strings.Contains(string(contentIgnored), "# Should Be Preserved") {
		t.Errorf("Expected content to be ignored (preserved), but was modified: %s", string(contentIgnored))
	}

	// 9. Verify Agent Switch (Enable Claude)
	// We need to update .kex.yaml to enable Claude (set strategies to overwrite)
	// Or run kex init again? Or kex update with config hack?
	// Let's hack .kex.yaml

	cfgData := `source: contents
update:
  documents:
    kex: all
  ai:
    targets: [claude]
    keywords: [documentation, kex]
`
	os.WriteFile(filepath.Join(dir, ".kex.yaml"), []byte(cfgData), 0644)

	// 8. Create Template for Claude (Simulating User Local Assets? No, binary uses embedded)
	// If binary uses embedded assets, writing files here does nothing unless Kex supports override.
	// Assuming Kex uses embedded assets:
	// We rely on the embedded manifest.json having "skills" configured.
	// We rely on the embedded templates having the skills template.
	// IF the embedded assets are used, we don't need to write templates here.
	// HOWEVER, for the sake of the test environment matching expectations if logic changes:
	// Let's rely on the binary which we just built (which includes our fixed manifest).

	runKex("update")

	claudeRulePath := filepath.Join(dir, ".claude/rules/kex/use-kex.md")
	if _, err := os.Stat(claudeRulePath); os.IsNotExist(err) {
		t.Errorf("Claude rules were not generated in .claude/rules/kex/use-kex.md")
	}

	// Verify Skill Generation
	// We expect "templates/contents/documentation/kex/write-concise-content.md" (which has "kex" keyword)
	// to be generated as ".claude/skills/kex/write-concise-content.md"
	// Wait, document ID for that file is "kex/write-concise-content".
	// The template is "{{.SkillName}}.md.template".
	// SkillName = "kex/write-concise-content".
	// Output Pattern = ".claude/skills/kex/{{.SkillName}}.md".
	// Result = ".claude/skills/kex/kex/write-concise-content.md"?
	// No, currently generator uses `filepath.Join("kex", SkillName)` in skills_generator?
	// No, we updated skills_generator to use the pattern completely.
	// The manifest says: ".claude/skills/kex/{{.SkillName}}.md.template"
	// So pattern is: ".claude/skills/kex/{{.SkillName}}.md"
	// If SkillName is "kex/write-concise-content", then path is ".claude/skills/kex/kex/write-concise-content.md".
	// This seems like double nesting.
	// Issue: SkillName (ID) includes the directory structure if `kex` is the corpus name.
	// If ID is `kex/write-concise-content`, duplicating `kex` in the pattern is bad.
	// The pattern in manifest is `.claude/skills/kex/{{.SkillName}}.md.template`.
	// If ID is `kex/foo`, output is `.claude/skills/kex/kex/foo.md`.
	// Maybe pattern should be `.claude/skills/{{.SkillName}}.md.template`?
	// Let's check what ID looks like. `domain.ParseDocument` assigns ID.
	// If my document is at `contents/documentation/kex/write-concise-content.md`.
	// ID is typically `documentation/kex/write-concise-content`?
	// Scope might be `documentation`, `kex`.
	// ID generation logic matches file path relative to root?
	// Let's assume ID is `documentation/kex/write-concise-content`.
	// Then output is `.claude/skills/kex/documentation/kex/write-concise-content.md`.
	// This is very nested.
	// But let's verify if file exists at all first.
	// We'll walk .claude/skills to see what was generated.

	foundSkill := false
	filepath.WalkDir(filepath.Join(dir, ".claude", "skills"), func(path string, d os.DirEntry, err error) error {
		// New structure: .claude/skills/kex.documentation.kex.write-concise-content/SKILL.md
		// We check if "write-concise-content" is in the path.
		if strings.Contains(path, "write-concise-content") {
			foundSkill = true
		}
		return nil
	})
	if !foundSkill {
		t.Errorf("Skill directory for write-concise-content was not generated in .claude/skills")
	}
}

func TestKexUpdate_CustomRoot(t *testing.T) {
	t.Run("it should respect configured root directory during update", func(t *testing.T) {
		tempDir := t.TempDir()

		// 1. Create .kex.yaml with custom root
		configContent := "source: custom_docs\nupdate:\n  documents:\n    kex: all\n"
		if err := os.WriteFile(filepath.Join(tempDir, ".kex.yaml"), []byte(configContent), 0644); err != nil {
			t.Fatal(err)
		}

		// 2. Run kex update
		cmd := exec.Command(kexBinary, "update")
		cmd.Dir = tempDir
		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("kex update failed: %v\nOutput: %s", err, output)
		}

		// 3. Verify content is in custom_docs/documentation/kex
		expectedPath := filepath.Join(tempDir, "custom_docs", "documentation", "kex", "write-concise-content.md")
		if _, err := os.Stat(expectedPath); os.IsNotExist(err) {
			t.Fatalf("Content expected at %s, but not found. (Bug Reproducible)", expectedPath)
		}

		// Verify content substitution
		content, err := os.ReadFile(expectedPath)
		if err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(string(content), "custom_docs/") {
			t.Errorf("Expected content to contain path context 'custom_docs/', but got: %s", string(content))
		}
	})
}
