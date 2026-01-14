package generator

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/mew-ton/kex/internal/infrastructure/config"
	kexfs "github.com/mew-ton/kex/internal/infrastructure/fs"
)

const ()

type Generator struct {
	Assets embed.FS
}

func New(assets embed.FS) *Generator {
	return &Generator{Assets: assets}
}

// Update updates the kex repository files based on configuration strategies
func (g *Generator) Update(cwd, localSource string, references []string, cfg config.UpdateConfig) error {
	manifest, err := LoadManifest(g.Assets)
	if err != nil {
		return err
	}

	// 1. Determine which files to update
	filesToUpdate := g.determineFilesToUpdate(manifest, cfg)

	// 2. Process Template Files (MCP Rules & System Docs)
	// Template files are usually written to the local project structure.
	// localSource determines where content-related templates go.
	if err := g.processTemplateFiles(cwd, localSource, filesToUpdate); err != nil {
		return err
	}

	// 3. Process AI Skills (Dynamic Content)
	// Skills can draw validation/context from both local source and references.
	if err := g.processAiSkills(cwd, localSource, references, cfg, manifest); err != nil {
		return err
	}

	return nil
}

func (g *Generator) determineFilesToUpdate(manifest *Manifest, cfg config.UpdateConfig) map[string]string {
	// path -> strategyName (e.g. "overwrite", "skip")
	filesToUpdate := make(map[string]string)

	// A. System Documents (legacy "kex" key support via Documents map)
	if cfg.Documents["kex"] == "all" {
		for _, f := range manifest.Kex {
			filesToUpdate[f] = "overwrite"
		}
	}

	// B. MCP Rules (AiMcpRules)
	g.resolveMcpFiles(manifest, cfg, filesToUpdate)

	return filesToUpdate
}

func (g *Generator) resolveMcpFiles(manifest *Manifest, cfg config.UpdateConfig, filesToUpdate map[string]string) {
	targets := cfg.AiMcpRules.Targets
	scopesList := cfg.AiMcpRules.Scopes
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
}

func (g *Generator) processTemplateFiles(cwd, rootDir string, filesToUpdate map[string]string) error {
	for relPath, strategyName := range filesToUpdate {
		strategy := ResolveStrategy(strategyName)
		if strategy == nil {
			continue
		}

		// Calculate Source and Target Paths
		templatePath := filepath.Join("templates", relPath)
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
	return nil
}

func (g *Generator) processAiSkills(cwd, localSource string, references []string, cfg config.UpdateConfig, manifest *Manifest) error {
	if len(cfg.AiSkills.Targets) == 0 || len(cfg.AiSkills.Keywords) == 0 {
		return nil
	}

	skillsGen := NewSkillsGenerator(cfg.AiSkills)
	searchRoots := g.resolveSearchRoots(cwd, localSource, references)
	skillTargets := cfg.AiSkills.Targets

	for _, t := range skillTargets {
		agentName := strings.TrimSpace(t)
		agentDef, ok := manifest.AiAgents[agentName]
		if !ok {
			continue
		}

		if err := g.processAgentSkills(cwd, agentDef, searchRoots, skillsGen); err != nil {
			return err
		}
	}
	return nil
}

func (g *Generator) resolveSearchRoots(cwd, localSource string, references []string) []string {
	var searchRoots []string
	if localSource != "" {
		searchRoots = append(searchRoots, filepath.Join(cwd, localSource))
	}
	for _, ref := range references {
		if strings.HasPrefix(ref, "http://") || strings.HasPrefix(ref, "https://") {
			// TODO: Support remote references in skills generation
			continue
		}

		if filepath.IsAbs(ref) {
			searchRoots = append(searchRoots, ref)
		} else {
			searchRoots = append(searchRoots, filepath.Join(cwd, ref))
		}
	}
	return searchRoots
}

func (g *Generator) processAgentSkills(cwd string, agentDef AgentDef, searchRoots []string, skillsGen *SkillsGenerator) error {
	for _, skillManifestPath := range agentDef.Files.Skills {
		skillTemplatePath := filepath.Join("templates", skillManifestPath)
		skillTemplateData, err := fs.ReadFile(g.Assets, skillTemplatePath)
		if err != nil {
			return fmt.Errorf("failed to read skill template %s: %w", skillTemplatePath, err)
		}

		targetPattern := strings.TrimSuffix(skillManifestPath, ".template")

		// Iterate all roots and merge results
		allSkills := make(map[string]string)

		for _, root := range searchRoots {
			if _, err := os.Stat(root); os.IsNotExist(err) {
				continue
			}

			skills, err := skillsGen.Generate(root, string(skillTemplateData), targetPattern)
			if err != nil {
				return fmt.Errorf("failed to generate skills for %s from root %s: %w", skillManifestPath, root, err)
			}
			for k, v := range skills {
				allSkills[k] = v
			}
		}

		for filename, content := range allSkills {
			outPath := filepath.Join(cwd, filename)
			if err := kexfs.WriteFile(outPath, []byte(content)); err != nil {
				return err
			}
		}
	}
	return nil
}
