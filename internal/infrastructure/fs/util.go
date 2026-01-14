package fs

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// EnsureDir ensures that the directory for the given path exists.
func EnsureDir(path string) error {
	return os.MkdirAll(path, 0755)
}

// WriteFile writes data to a file, creating it if necessary.
func WriteFile(path string, data []byte) error {
	if err := EnsureDir(filepath.Dir(path)); err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

// CopyFile copies a file from src to dst.
func CopyFile(src, dst string) error {
	if err := EnsureDir(filepath.Dir(dst)); err != nil {
		return err
	}

	s, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("failed to open source: %w", err)
	}
	defer s.Close()

	d, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("failed to create destination: %w", err)
	}
	defer d.Close()

	_, err = io.Copy(d, s)
	if err != nil {
		return fmt.Errorf("failed to copy content: %w", err)
	}
	return nil
}
