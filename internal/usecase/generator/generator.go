package generator

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/mew-ton/kex/internal/infrastructure/config"
)

type AgentType string

const (
	AgentTypeGeneral AgentType = "general"
	AgentTypeClaude  AgentType = "claude"
)

type Generator struct {
	Templates embed.FS
}

func New(templates embed.FS) *Generator {
	return &Generator{Templates: templates}
}

func (g *Generator) Generate(cwd string, agentType AgentType, agentConfig *config.Agent) error {
	// Extract templates mirroring the structure in assets/templates
	err := fs.WalkDir(g.Templates, "templates", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel("templates", path)
		if err != nil {
			return err
		}

		if relPath == "." {
			return nil
		}

		// Filter based on AgentType
		// General -> Skip .claude
		// Claude -> Skip .agent
		if agentType == AgentTypeGeneral && strings.HasPrefix(relPath, ".claude") {
			return nil
		}
		if agentType == AgentTypeClaude && strings.HasPrefix(relPath, ".agent") {
			return nil
		}

		// Skip directories themselves (we create them as needed)
		if d.IsDir() {
			return nil
		}

		// Filter rules based on scope
		// Check for rules in any agent directory
		if strings.Contains(relPath, "/rules/") {
			if agentConfig != nil && len(agentConfig.Scopes) > 0 {
				filename := filepath.Base(relPath)
				// Check if the rule file corresponds to an enabled scope
				// Convention: kex-<scope>.md
				// e.g. kex-coding.md -> scope: coding
				scopeName := strings.TrimPrefix(strings.TrimSuffix(filename, ".md"), "kex-")
				// Map imperative filenames to scopes
				if filename == "follow-coding-rules.md" {
					scopeName = "coding"
				} else if filename == "follow-documentation-rules.md" {
					scopeName = "documentation"
				}

				scopeEnabled := false
				for _, s := range agentConfig.Scopes {
					if s == scopeName {
						scopeEnabled = true
						break
					}
				}

				if !scopeEnabled {
					return nil // Skip this rule file
				}
			}
		}

		targetPath := filepath.Join(cwd, relPath)

		// Ensure parent dir exists
		if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
			return err
		}

		data, err := fs.ReadFile(g.Templates, path)
		if err != nil {
			return err
		}

		// Don't overwrite if exists
		if _, err := os.Stat(targetPath); err == nil {
			return nil
		}

		// Substitute template variables
		// Currently only .Root is supported
		contentStr := string(data)

		// For Generate (Init), we assume default root "contents"
		effectiveRoot := "contents"
		if !strings.HasPrefix(effectiveRoot, "./") && !strings.HasPrefix(effectiveRoot, "/") {
			effectiveRoot = "./" + effectiveRoot
		}

		contentStr = strings.ReplaceAll(contentStr, "{{.Root}}", effectiveRoot)
		data = []byte(contentStr)

		return os.WriteFile(targetPath, data, 0644)
	})

	if err != nil {
		return fmt.Errorf("failed to extract templates: %w", err)
	}

	return nil
}

// Update updates the kex repository files based on configuration
func (g *Generator) Update(cwd, rootDir string, agentType AgentType, strategies config.Strategies, agentConfig *config.Agent) error {
	return fs.WalkDir(g.Templates, "templates", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel("templates", path)
		if err != nil {
			return err
		}

		if relPath == "." {
			return nil
		}

		// Filter based on AgentType
		if agentType == AgentTypeGeneral && strings.HasPrefix(relPath, ".claude") {
			return nil
		}
		if agentType == AgentTypeClaude && strings.HasPrefix(relPath, ".agent") {
			return nil
		}

		if d.IsDir() {
			return nil
		}

		// Filter rules based on scope
		if strings.Contains(relPath, "/rules/") {
			if agentConfig != nil && len(agentConfig.Scopes) > 0 {
				filename := filepath.Base(relPath)
				scopeName := strings.TrimPrefix(strings.TrimSuffix(filename, ".md"), "kex-")

				scopeEnabled := false
				for _, s := range agentConfig.Scopes {
					if s == scopeName {
						scopeEnabled = true
						break
					}
				}

				if !scopeEnabled {
					return nil // Skip
				}
			}
		}

		mappedPath := relPath
		// Use custom root for contents
		if rootDir != "" && strings.HasPrefix(mappedPath, "contents") {
			mappedPath = filepath.Join(rootDir, strings.TrimPrefix(mappedPath, "contents"))
		}

		targetPath := filepath.Join(cwd, mappedPath)

		data, err := fs.ReadFile(g.Templates, path)
		if err != nil {
			return err
		}

		// Substitute template variables
		// Currently only .Root is supported
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

		// RELPATH (from template structure) is used for strategy lookup to match .agent/... or .claude/...
		// But wait, strategy lookup expects the CANONICAL path.
		// If relPath is .claude/rules/..., that is the canonical path.
		// If relPath is contents/..., that is the canonical path.
		// So passing relPath to ResolveStrategy is correct.
		strategy := ResolveStrategy(relPath, strategies)
		ctx := UpdateContext{
			TargetPath:   targetPath,
			TemplateData: data,
			AgentConfig:  agentConfig,
			Generator:    g,
		}

		return strategy.Apply(ctx)
	})
}
