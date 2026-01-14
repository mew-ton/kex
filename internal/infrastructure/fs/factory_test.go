package fs_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/mew-ton/kex/internal/infrastructure/config"
	"github.com/mew-ton/kex/internal/infrastructure/fs"
	"github.com/mew-ton/kex/internal/infrastructure/logger"
)

func TestProviderFactory_CreateProvider(t *testing.T) {
	// Setup temporary directory for local tests
	tmpDir, err := os.MkdirTemp("", "kex-factory-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create a dummy file
	dummyFile := filepath.Join(tmpDir, "dummy")
	if err := os.WriteFile(dummyFile, []byte("content"), 0644); err != nil {
		t.Fatal(err)
	}

	logger := &logger.NoOpLogger{}
	cfg := config.Config{}
	factory := fs.NewProviderFactory(cfg, logger)

	tests := []struct {
		name        string
		pathOrURL   string
		isReference bool
		cwd         string
		wantErr     bool
		wantType    string // "Local" or "Remote"
	}{
		{
			name:      "Remote URL",
			pathOrURL: "https://example.com/kex.json",
			wantErr:   false,
			wantType:  "Remote",
		},
		{
			name:      "Local Absolute File",
			pathOrURL: dummyFile,
			cwd:       "",
			wantErr:   false,
			wantType:  "Local",
		},
		{
			name:      "Local Relative File",
			pathOrURL: "dummy",
			cwd:       tmpDir,
			wantErr:   false,
			wantType:  "Local",
		},
		{
			name:        "Non-existent Local Source",
			pathOrURL:   "non-existent",
			cwd:         tmpDir,
			isReference: false,
			wantErr:     true,
		},
		{
			name:        "Non-existent Local Reference",
			pathOrURL:   "non-existent-ref",
			cwd:         tmpDir,
			isReference: true,
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p, resolved, err := factory.CreateProvider(tt.pathOrURL, tt.isReference, tt.cwd)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateProvider() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if p == nil {
					t.Error("CreateProvider() returned nil provider")
				}
				if tt.wantType == "Remote" {
					if _, ok := p.(*fs.RemoteProvider); !ok {
						t.Errorf("CreateProvider() returned type %T, want RemoteProvider", p)
					}
					if resolved != tt.pathOrURL {
						t.Errorf("resolved path = %s, want %s", resolved, tt.pathOrURL)
					}
				} else {
					if _, ok := p.(*fs.LocalProvider); !ok {
						t.Errorf("CreateProvider() returned type %T, want LocalProvider", p)
					}
					// For local, verify resolved path is absolute
					if !filepath.IsAbs(resolved) {
						t.Errorf("resolved path %s is not absolute", resolved)
					}
				}
			}
		})
	}
}
