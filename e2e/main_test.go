package e2e

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
)

var (
	kexBinary string
)

func TestMain(m *testing.M) {
	// Setup: Build the kex binary
	tempDir, err := os.MkdirTemp("", "kex-e2e-")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create temp dir: %v\n", err)
		os.Exit(1)
	}
	defer os.RemoveAll(tempDir)

	if runtime.GOOS == "windows" {
		kexBinary = filepath.Join(tempDir, "kex.exe")
	} else {
		kexBinary = filepath.Join(tempDir, "kex")
	}

	cmd := exec.Command("go", "build", "-o", kexBinary, "../cmd/kex")
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to build kex: %v\n", err)
		os.Exit(1)
	}

	// Run tests
	exitCode := m.Run()

	// Cleanup
	os.Exit(exitCode)
}
