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

	// 3. Setup user content in .antigravity/rules/kex-coding.md
	rulePath := filepath.Join(dir, ".antigravity/rules/kex-coding.md")

	initialContent, _ := os.ReadFile(rulePath)
	// Inject a custom header line
	customHeader := "# My Custom Rules\n" + string(initialContent)
	err = os.WriteFile(rulePath, []byte(customHeader), 0644)
	if err != nil {
		t.Fatalf("Failed to write kex-coding.md: %v", err)
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

	// 6. Verify kex-coding.md is OVERWRITTEN (Because kex init sets it to overwrite)
	// ISSUE: the test expected "preserved" before because defaults were "skip".
	// Now defaults are "ignore", but init sets to "overwrite".
	// If I want to test "skip" (Create Only), I need to modify config.

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

	strategiesFallback := config.Strategies{} // All ignore
	// If logic is "Ignore", then Update() does nothing. File preserves state.

	gen := generator.New(assets.Assets)
	err = gen.Update(dir, "", strategiesFallback)
	if err != nil {
		t.Fatalf("Update (Empty Strategy) failed: %v", err)
	}
	contentIgnored, _ := os.ReadFile(rulePath)
	if !strings.Contains(string(contentIgnored), "# Should Be Preserved") {
		t.Errorf("Expected content to be ignored (preserved), but was modified: %s", string(contentIgnored))
	}

	// 9. Verify Agent Switch (Enable Claude)
	// We need to update .kex.yaml to enable Claude (set strategies to overwrite)
	// Or run kex init again? Or kex update with config hack?
	// Let's hack .kex.yaml

	cfgData := `root: contents
update:
  strategies:
    claude: all
`
	os.WriteFile(filepath.Join(dir, ".kex.yaml"), []byte(cfgData), 0644)

	runKex("update")

	claudeRulePath := filepath.Join(dir, ".claude/rules/kex/follow-coding-rules.md")
	if _, err := os.Stat(claudeRulePath); os.IsNotExist(err) {
		t.Errorf("Claude rules were not generated in .claude/rules/kex/follow-coding-rules.md")
	}
}

func TestKexUpdate_CustomRoot(t *testing.T) {
	t.Run("it should respect configured root directory during update", func(t *testing.T) {
		tempDir := t.TempDir()

		// 1. Create .kex.yaml with custom root
		configContent := "root: custom_docs\nupdate:\n  strategies:\n    kex: all\n"
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
