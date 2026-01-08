package e2e

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestKexCheck_Composite(t *testing.T) {
	t.Run("it should valid documents from multiple local roots passed via flags", func(t *testing.T) {
		tempDir := t.TempDir()
		root1 := filepath.Join(tempDir, "root1")
		root2 := filepath.Join(tempDir, "root2")
		os.Mkdir(root1, 0755)
		os.Mkdir(root2, 0755)

		// Doc 1 in Root 1
		doc1 := `---
id: doc-1
title: Doc 1
keywords: [one]
---
Content 1`
		os.WriteFile(filepath.Join(root1, "doc-1.md"), []byte(doc1), 0644)

		// Doc 2 in Root 2
		doc2 := `---
id: doc-2
title: Doc 2
keywords: [two]
---
Content 2`
		os.WriteFile(filepath.Join(root2, "doc-2.md"), []byte(doc2), 0644)

		cmd := exec.Command(kexBinary, "check", root1, root2)
		cmd.Dir = tempDir
		output, err := cmd.CombinedOutput()

		if err != nil {
			t.Fatalf("Check failed expected success: %v\nOutput: %s", err, string(output))
		}

		if !strings.Contains(string(output), "All checks passed") {
			t.Errorf("Expected success message, got: %s", string(output))
		}
	})

	t.Run("it should validate documents from multiple local roots passed via positional args", func(t *testing.T) {
		tempDir := t.TempDir()
		root1 := filepath.Join(tempDir, "pos1")
		root2 := filepath.Join(tempDir, "pos2")
		os.Mkdir(root1, 0755)
		os.Mkdir(root2, 0755)

		os.WriteFile(filepath.Join(root1, "pos-1.md"), []byte(`---
id: pos-1
title: Pos 1
---
C1`), 0644)

		os.WriteFile(filepath.Join(root2, "pos-2.md"), []byte(`---
id: pos-2
title: Pos 2
---
C2`), 0644)

		cmd := exec.Command(kexBinary, "check", root1, root2)
		cmd.Dir = tempDir
		output, err := cmd.CombinedOutput()

		if err != nil {
			t.Fatalf("Check failed expected success: %v\nOutput: %s", err, string(output))
		}
	})

	t.Run("it should fail when duplicate IDs exist across roots", func(t *testing.T) {
		tempDir := t.TempDir()
		root1 := filepath.Join(tempDir, "dup1")
		root2 := filepath.Join(tempDir, "dup2")
		os.Mkdir(root1, 0755)
		os.Mkdir(root2, 0755)

		docContent := `---
id: DUPLICATE_ID
title: Duplicate
keywords: [dup]
---
Content`
		os.WriteFile(filepath.Join(root1, "doc.md"), []byte(docContent), 0644)
		os.WriteFile(filepath.Join(root2, "doc_copy.md"), []byte(docContent), 0644)

		cmd := exec.Command(kexBinary, "check", root1, root2)
		cmd.Dir = tempDir
		output, err := cmd.CombinedOutput()

		// Expect failure
		if err == nil {
			t.Error("Expected check failure due to duplicate ID, but got success")
		}

		outStr := string(output)
		if !strings.Contains(outStr, "duplicate document ID") && !strings.Contains(outStr, "DUPLICATE_ID") {
			t.Errorf("Expected error message about duplicate ID, got: %s", outStr)
		}
	})
}
