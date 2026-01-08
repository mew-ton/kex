package fs

import (
	"path/filepath"
	"strings"

	"github.com/mew-ton/kex/internal/infrastructure/config"
	"github.com/mew-ton/kex/internal/infrastructure/logger"
)

// Factory creates a repository (Indexer) from configuration
type Factory struct {
	Logger logger.Logger
}

func NewFactory(l logger.Logger) *Factory {
	return &Factory{Logger: l}
}

// CreateRepository creates an Indexer based on the provided configuration and runtime source overrides
func (f *Factory) CreateRepository(cfg config.Config, cliSources []string, projectRoot string) (*Indexer, error) {
	// 1. Determine effective sources
	// CLI overrides config if provided
	var effectiveSources []string

	if len(cliSources) > 0 {
		effectiveSources = cliSources
	} else {
		// Use Source from config
		effectiveSources = []string{cfg.Source}
	}

	// 2. Create Providers for each source
	var providers []DocumentProvider

	for _, source := range effectiveSources {
		// Differentiate Local vs Remote
		if isRemote(source) {
			// Remote Provider
			token := cfg.RemoteToken // Use token from config (Environment variable handling should happen before calling this or inside provider)
			// Ideally environment variable is handled by the caller or passed in config?
			// The original code checked env var in createRepository.
			// We should assume cfg has the right token or pass it explicitly?
			// Let's pass the raw token from cfg for now. The CLI layer usually resolves Env Vs Config.

			p := NewRemoteProvider(source, token, f.Logger) // RemoteProvider signature might need check
			providers = append(providers, p)
		} else {
			// Local Provider
			var localPath string
			if filepath.IsAbs(source) {
				localPath = source
			} else {
				localPath = filepath.Join(projectRoot, source)
			}

			p := NewLocalProvider(localPath, f.Logger)
			providers = append(providers, p)
		}
	}

	// 3. Create Composite Provider
	// Even if there is only 1 provider, we use CompositeProvider to unify logic (ID collision check, etc)
	// Although CompositeProvider overhead is minimal.
	composite := NewCompositeProvider(providers)

	// 4. Create Indexer
	repo := New(composite, f.Logger)

	// 5. Load (Should this happen here? Original code called Load explicitly)
	// Usually Factory just returns the instance. Caller calls Load.
	return repo, nil
}

func isRemote(path string) bool {
	return strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://")
}
