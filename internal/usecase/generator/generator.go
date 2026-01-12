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

const ()

type Generator struct {
	Assets embed.FS
}

func New(assets embed.FS) *Generator {
	return &Generator{Assets: assets}
}

// Update updates the kex repository files based on configuration strategies
func (g *Generator) Update(cwd, rootDir string, cfg config.UpdateConfig) error {
	manifest, err := LoadManifest(g.Assets)
	if err != nil {
		return err
	}

	// 1. Determine which files to update
	// path -> strategyName (e.g. "overwrite", "skip")
	filesToUpdate := make(map[string]string)

	// A. System Documents (legacy "kex" key support via Documents map)
	// If "kex": "all" is present in Documents map
	if cfg.Documents["kex"] == "all" {
		for _, f := range manifest.Kex {
			filesToUpdate[f] = "overwrite"
		}
	}

	// B. MCP Rules (AiMcpRules)
	// Targets: comma separated string, e.g. "antigravity, claude"
	targets := strings.Split(cfg.AiMcpRules.Targets, ",")
	scopesList := strings.Split(cfg.AiMcpRules.Scopes, ",")
	// clean up whitespace
	for i := range scopesList {
		scopesList[i] = strings.TrimSpace(scopesList[i])
	}

	scopeStrategies := ResolveFileScopes(scopesList)

	for _, t := range targets {
		agentName := strings.TrimSpace(t)
		if agentName == "" {
			continue
		}

		agentDef, ok := manifest.AiAgents[agentName]
		if !ok {
			// Unknown agent, maybe warn? or ignore?
			// For now ignore to avoid noise if user mistypes
			continue
		}

		// Apply scopes
		for _, scopeStr := range scopeStrategies {
			files := scopeStr.SelectFiles(agentDef)
			for _, f := range files {
				filesToUpdate[f] = "overwrite"
			}
		}
	}

	// 2. Process Template Files (MCP Rules & System Docs)
	for relPath, strategyName := range filesToUpdate {
		strategy := ResolveStrategy(strategyName)
		if strategy == nil {
			continue
		}

		// Calculate Source and Target Paths
		// 1. Source: explicit relPath from manifest (e.g., "templates/foo.md.template")
		templatePath := filepath.Join("templates", relPath)

		// 2. Target: strip .template if present, and map to rootDir if needed
		targetRelPath := strings.TrimSuffix(relPath, ".template")

		// Map contents/... paths to respect rootDir
		mappedPath := targetRelPath
		if rootDir != "" && strings.HasPrefix(mappedPath, "contents") {
			mappedPath = filepath.Join(rootDir, strings.TrimPrefix(mappedPath, "contents"))
		}

		targetPath := filepath.Join(cwd, mappedPath)

		data, err := fs.ReadFile(g.Assets, templatePath)
		if err != nil {
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
			Strategy:     strategyName,
			Generator:    g,
		}

		if err := strategy.Apply(ctx); err != nil {
			return err
		}
	}

	// 3. Process AI Skills (Dynamic Content)
	if cfg.AiSkills.Targets != "" && len(cfg.AiSkills.Keywords) > 0 {
		skillsGen := NewSkillsGenerator(cfg.AiSkills)

		// Determine source directory for documents (root/contents or just root?)
		// Config.Root points to where contents are.
		// If rootDir is "contents", we assume documents are in cwd/contents.
		contentSourceDir := filepath.Join(cwd, rootDir)

		// Determine targets
		skillTargets := strings.Split(cfg.AiSkills.Targets, ",")
		for _, t := range skillTargets {
			agentName := strings.TrimSpace(t)
			agentDef, ok := manifest.AiAgents[agentName]
			if !ok {
				continue
			}

			// Iterate over configured skills templates in manifest
			for _, skillManifestPath := range agentDef.Files.Skills {
				// Source: templates/ + manifest path
				skillTemplatePath := filepath.Join("templates", skillManifestPath)
				skillTemplateData, err := fs.ReadFile(g.Assets, skillTemplatePath)
				if err != nil {
					return fmt.Errorf("failed to read skill template %s: %w", skillTemplatePath, err)
				}

				// Target Pattern: strip .template
				// e.g. .claude/skills/kex/{{.SkillName}}.md.template -> .claude/skills/kex/{{.SkillName}}.md
				targetPattern := strings.TrimSuffix(skillManifestPath, ".template")

				// Generate with explicit pattern
				skills, err := skillsGen.Generate(contentSourceDir, string(skillTemplateData), targetPattern)
				if err != nil {
					return fmt.Errorf("failed to generate skills for %s: %w", skillManifestPath, err)
				}

				// Write generated skills
				for filename, content := range skills {
					outPath := filepath.Join(cwd, filename)
					// Strategy: overwrite? Use simple overwrite logic
					if err := EnsureDir(filepath.Dir(outPath)); err != nil {
						return err
					}
					if err := WriteFile(outPath, []byte(content)); err != nil {
						return err
					}
				}
			}
		}
	}

	return nil
}

// Simple helpers - duplicate from internal/infrastructure/fs or similar if strictly needed there,
// but for Generator package utils are fine.

func EnsureDir(path string) error {
	return os.MkdirAll(path, 0755)
}

func WriteFile(path string, data []byte) error {
	return os.WriteFile(path, data, 0644)
}
