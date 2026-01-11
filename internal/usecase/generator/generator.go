package generator

import (
	"embed"
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/mew-ton/kex/internal/infrastructure/config"
)

const ()

type Generator struct {
	Assets embed.FS
}

func New(assets embed.FS) *Generator {
	return &Generator{Assets: assets}
}

// Update updates the kex repository files based on configuration strategies
func (g *Generator) Update(cwd, rootDir string, strategies config.Strategies) error {
	manifest, err := LoadManifest(g.Assets)
	if err != nil {
		return err
	}

	// 1. Determine which files to update
	// path -> strategyName (e.g. "overwrite", "skip")
	filesToUpdate := make(map[string]string)

	// Kex Files
	// Check "kex" config. Default to "skip" (create only) if missing? or "all"?
	// If "kex" key is present with "all", generate everything with "overwrite".
	// If "kex" key is missing, defaults?
	// Given earlier requirements, "kex init" sets "kex: all" (to be implemented).
	// So we can be strict: only generate if configured.
	// But what about updates? If user has no "kex" key (legacy or manual), we might want to default to something.
	// For now, let's respect the map. If "kex": "all" -> overwrite.
	kexMode := strategies["kex"]
	if kexMode == "all" {
		for _, f := range manifest.Kex {
			filesToUpdate[f] = "overwrite" // Enforce overwrite for kex system files
		}
	}

	// Agent Files
	for agentKey, agentDef := range manifest.AiAgents {
		mode := strategies[agentKey]
		if mode == "" || mode == "none" {
			continue
		}

		targetStrategy := "overwrite" // Default to overwrite (enforce kex management)
		// "skip" mode was removed. If user wants to stop updates, they should remove the key or set to "none".

		scopes := ResolveFileScopes(mode)
		for _, scope := range scopes {
			files := scope.SelectFiles(agentDef)
			// Helper to add files
			for _, f := range files {
				filesToUpdate[f] = targetStrategy
			}
		}
	}

	// 2. Process Files
	for relPath, strategyName := range filesToUpdate {
		strategy := ResolveStrategy(strategyName)
		if strategy == nil {
			continue
		}

		// Map contents/... paths to respect rootDir
		mappedPath := relPath
		if rootDir != "" && strings.HasPrefix(mappedPath, "contents") {
			mappedPath = filepath.Join(rootDir, strings.TrimPrefix(mappedPath, "contents"))
		}

		targetPath := filepath.Join(cwd, mappedPath)
		templatePath := filepath.Join("templates", relPath)

		data, err := fs.ReadFile(g.Assets, templatePath)
		if err != nil {
			// If template is missing but present in manifest, that's an issue.
			return fmt.Errorf("failed to read template %s: %w", templatePath, err)
		}

		// Substitute template variables
		contentStr := string(data)

		// Determine root path for display/template
		effectiveRoot := rootDir
		if effectiveRoot == "" {
			effectiveRoot = "contents"
		}
		if !strings.HasPrefix(effectiveRoot, "./") && !strings.HasPrefix(effectiveRoot, "/") {
			effectiveRoot = "./" + effectiveRoot
		}

		contentStr = strings.ReplaceAll(contentStr, "{{.Root}}", effectiveRoot)
		data = []byte(contentStr)

		ctx := UpdateContext{
			TargetPath:   targetPath,
			TemplateData: data,
			Strategy:     strategyName, // Just for context/logging if needed
			Generator:    g,
		}

		if err := strategy.Apply(ctx); err != nil {
			return err
		}
	}

	return nil
}
