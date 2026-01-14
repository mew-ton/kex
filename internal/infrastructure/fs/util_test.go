package fs_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/mew-ton/kex/internal/infrastructure/fs"
)

func TestEnsureDir(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "kex-util-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	target := filepath.Join(tmpDir, "subdir", "deep")
	if err := fs.EnsureDir(target); err != nil {
		t.Errorf("EnsureDir() error = %v", err)
	}

	if _, err := os.Stat(target); os.IsNotExist(err) {
		t.Errorf("EnsureDir() did not create directory %s", target)
	}
}

func TestCopyFile(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "kex-util-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	src := filepath.Join(tmpDir, "src.txt")
	dst := filepath.Join(tmpDir, "dst", "dst.txt")
	content := []byte("hello world")

	if err := os.WriteFile(src, content, 0644); err != nil {
		t.Fatal(err)
	}

	if err := fs.CopyFile(src, dst); err != nil {
		t.Errorf("CopyFile() error = %v", err)
	}

	got, err := os.ReadFile(dst)
	if err != nil {
		t.Fatal(err)
	}

	if string(got) != string(content) {
		t.Errorf("CopyFile() content = %s, want %s", got, content)
	}
}

func TestWriteFile(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "kex-util-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	target := filepath.Join(tmpDir, "subdir", "file.txt")
	content := []byte("test content")

	if err := fs.WriteFile(target, content); err != nil {
		t.Errorf("WriteFile() error = %v", err)
	}

	got, err := os.ReadFile(target)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != string(content) {
		t.Errorf("WriteFile() content = %s, want %s", got, content)
	}
}
