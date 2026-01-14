package fs

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/mew-ton/kex/internal/infrastructure/config"
	"github.com/mew-ton/kex/internal/infrastructure/logger"
)

// ProviderFactory handles the creation of DocumentProviders
type ProviderFactory struct {
	logger logger.Logger
	cfg    config.Config
}

func NewProviderFactory(cfg config.Config, l logger.Logger) *ProviderFactory {
	return &ProviderFactory{
		cfg:    cfg,
		logger: l,
	}
}

// CreateProvider creates a DocumentProvider for the given path or URL.
// It handles local paths and remote URLs, including token resolution for remote sources.
func (f *ProviderFactory) CreateProvider(pathOrURL string, isReference bool, cwd string) (DocumentProvider, string, error) {
	if isURL(pathOrURL) {
		return f.createRemoteProvider(pathOrURL)
	}
	return f.createLocalProvider(pathOrURL, isReference, cwd)
}

func (f *ProviderFactory) createRemoteProvider(url string) (DocumentProvider, string, error) {
	token := os.Getenv("KEX_REMOTE_TOKEN")
	if token == "" && f.cfg.RemoteToken != "" {
		token = f.cfg.RemoteToken
	}

	// Logging is handled by the caller or provider itself usually, but check.go/start.go did some stdout logging.
	// We'll leave UI logging to the CLI layer, this factory just returns the provider.
	return NewRemoteProvider(url, token, f.logger), url, nil
}

func (f *ProviderFactory) createLocalProvider(path string, isReference bool, cwd string) (DocumentProvider, string, error) {
	var fullPath string
	if filepath.IsAbs(path) {
		fullPath = path
	} else {
		fullPath = filepath.Join(cwd, path)
	}

	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		if isReference {
			return nil, "", fmt.Errorf("reference '%s' not found", path)
		}
		return nil, "", fmt.Errorf("source '%s' not found", path)
	}

	return NewLocalProvider(fullPath, f.logger), fullPath, nil
}

func isURL(s string) bool {
	return strings.HasPrefix(s, "http://") || strings.HasPrefix(s, "https://")
}
