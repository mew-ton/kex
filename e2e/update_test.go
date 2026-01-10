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
	runKex("init", "--agent-type=general")

	// 2. Modify a system document (should be overwritten)
	targetDoc := filepath.Join(dir, "contents/documentation/kex/choose-effective-keywords.md")
	err := os.WriteFile(targetDoc, []byte("Modified Content"), 0644)
	if err != nil {
		t.Fatalf("Failed to modify file: %v", err)
	}

	// 3. Setup user content in .agent/rules/kex-coding.md
	rulePath := filepath.Join(dir, ".agent/rules/kex-coding.md")
	// Note: The template places the dogfooding line APTER the marker block.
	// We want to test that content OUTSIDE the marker block is preserved.
	// But our template has content BEFORE and AFTER marker block?
	// Actually kex-coding.md has header (outside), markers, footer (outside).
	// Let's modify the header part.

	initialContent, _ := os.ReadFile(rulePath)
	// Inject a custom header line
	customHeader := "# My Custom Rules\n" + string(initialContent)
	err = os.WriteFile(rulePath, []byte(customHeader), 0644)
	if err != nil {
		t.Fatalf("Failed to write kex-coding.md: %v", err)
	}

	// 4. Run Update (general agent default)
	runKex("update")

	// 5. Verify system document is reverted (Overwritten)
	content, err := os.ReadFile(targetDoc)
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}
	if string(content) == "Modified Content" {
		t.Errorf("System document was NOT overwritten/updated")
	}

	// 6. Verify kex-coding.md is PRESERVED (New Default Behavior: Rules default to skip)
	contentPreservedRule, _ := os.ReadFile(rulePath)
	if !strings.Contains(string(contentPreservedRule), "# My Custom Rules") {
		t.Errorf("Rule file should have been preserved (skipped) by default update, but was overwritten")
	}

	// 7. Test Strategy Override via Config (Force Overwrite)
	// Modify content again to verify we can force overwrite it
	customContent2 := "# Another Custom Rule\n"
	err = os.WriteFile(rulePath, []byte(customContent2), 0644)
	if err != nil {
		t.Fatalf("Failed to modify kex-coding.md for overwrite test: %v", err)
	}

	// Configure kex-coding.md to overwrite (Explicit override)
	gen := generator.New(assets.Templates)
	agentConfig := &config.Agent{Type: "general", Scopes: []string{"coding", "documentation"}}
	strategies := config.Strategies{
		AgentKexCoding: "overwrite",
	}

	err = gen.Update(dir, "", generator.AgentTypeGeneral, strategies, agentConfig)
	if err != nil {
		t.Fatalf("Update (Overwrite Strategy) failed: %v", err)
	}

	contentOverwritten, _ := os.ReadFile(rulePath)
	if strings.Contains(string(contentOverwritten), "# Another Custom Rule") {
		t.Errorf("Expected content to be overwritten to template, but custom content persisted")
	}
	if !strings.Contains(string(contentOverwritten), "Coding Rules") {
		t.Errorf("Expected content to be overwritten to template (missing header)")
	}

	// 8. Test Default Fallback with Empty Config (Pass empty struct)
	// Resolution: Empty string in struct -> ResolveStrategy defaults to Overwrite.
	strategiesFallback := config.Strategies{}

	err = gen.Update(dir, "", generator.AgentTypeGeneral, strategiesFallback, agentConfig)
	if err != nil {
		t.Fatalf("Update (Empty Strategy) failed: %v", err)
	}
	// Content should still be overwritten (template)
	contentFallback, _ := os.ReadFile(rulePath)
	if !strings.Contains(string(contentFallback), "Coding Rules") {
		t.Errorf("Expected content to be standard template after empty strategy update")
	}

	// 7. Verify Agent Type Switch (General -> Claude)
	runKex("update", "--agent-type=claude")

	claudeRulePath := filepath.Join(dir, ".claude/rules/kex/follow-coding-rules.md")
	if _, err := os.Stat(claudeRulePath); os.IsNotExist(err) {
		t.Errorf("Claude rules were not generated in .claude/rules/kex/follow-coding-rules.md")
	}

	// 8. Verify Configural Update (Scope filtering)
	// Modify .kex.yaml to only have "documentation" scope
	configData := `root: contents
agent:
  type: general
  scopes: ["documentation"]
`
	if err := os.WriteFile(filepath.Join(dir, ".kex.yaml"), []byte(configData), 0644); err != nil {
		t.Fatal(err)
	}

	// Run Update again as general
	runKex("update", "--agent-type=general")

	// Check if coding rule is removed?
	// The generator does NOT delete files. It just doesn't generate/update valid ones.
	// But kex-documentation.md should be updated/generated.
	// Check if kex-documentation.md has correct templated path
	docRulePath := filepath.Join(dir, ".agent/rules/kex-documentation.md")
	docRuleContent, _ := os.ReadFile(docRulePath)
	if !strings.Contains(string(docRuleContent), "documents under `./contents`") {
		t.Errorf("Documentation rule content should contain default root path `./contents`")
	}
}

func TestKexUpdate_CustomRoot(t *testing.T) {
	t.Run("it should respect configured root directory during update", func(t *testing.T) {
		tempDir := t.TempDir()

		// 1. Create .kex.yaml with custom root
		configContent := "root: custom_docs\nagent:\n  type: general\n  scopes: []\n"
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
	})
}
