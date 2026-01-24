package generator

import (
	"bytes"
	"embed"
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/mew-ton/kex/internal/domain"
	"github.com/mew-ton/kex/internal/infrastructure/config"
	kexfs "github.com/mew-ton/kex/internal/infrastructure/fs"
	"github.com/mew-ton/kex/internal/infrastructure/logger"
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

	// Initialize Infrastructure for Document Loading (Indexer)
	// We use stderr logger for generator output
	l := logger.NewStderrLogger()

	// Provider Factory (needs generic config)
	// We reconstruct a partial config for factory
	factoryCfg := config.Config{
		Source:     localSource,
		References: references,
	}
	factory := kexfs.NewProviderFactory(factoryCfg, l)

	var providers []kexfs.DocumentProvider

	// Add Local Source
	if localSource != "" {
		p, _, err := factory.CreateProvider(localSource, false, cwd)
		if err != nil {
			return fmt.Errorf("failed to create local provider for skills: %w", err)
		}
		providers = append(providers, p)
	}

	// Add References
	for _, ref := range references {
		p, _, err := factory.CreateProvider(ref, true, cwd)
		if err != nil {
			// Warn but continue for references?
			fmt.Printf("Warning: failed to load reference '%s': %v\n", ref, err)
			continue
		}
		providers = append(providers, p)
	}

	if len(providers) == 0 {
		return nil
	}

	compositeProvider := kexfs.NewCompositeProvider(providers)
	indexer := kexfs.New(compositeProvider, l)

	if err := indexer.Load(); err != nil {
		return fmt.Errorf("failed to load documents for skills generation: %w", err)
	}

	// Execute Search (Exact Match on Keywords as Scopes? Or just Keywords?)
	// Implementation Plan says: Execute search: `search.Execute(keywords, "", true)` (exactScopeMatch=true).
	// But we are calling Indexer directly.
	// indexer.Search(keywords, scopes, exactScopeMatch)
	// If we use exactScopeMatch=true, then `keywords` are treated as valid scopes.
	// We pass `nil` for explicit scopes because the keywords provided in config ARE the scopes we want.
	docs := indexer.Search(cfg.AiSkills.Keywords, nil, true)

	// Ensure content is loaded

	// Ensure content is loaded
	var loadedDocs []*domain.Document
	for _, d := range docs {
		fullDoc, ok := indexer.GetByID(d.ID)
		if ok {
			loadedDocs = append(loadedDocs, fullDoc)
		}
	}

	// Generate Skills for each Agent
	for _, t := range cfg.AiSkills.Targets {
		agentName := strings.TrimSpace(t)
		agentDef, ok := manifest.AiAgents[agentName]
		if !ok {
			continue
		}

		if err := g.processAgentSkills(cwd, agentDef, loadedDocs); err != nil {
			return err
		}
	}
	return nil
}

func (g *Generator) processAgentSkills(cwd string, agentDef AgentDef, docs []*domain.Document) error {
	for _, skillManifestPath := range agentDef.Files.Skills {
		// Load Template
		skillTemplatePath := filepath.Join("templates", skillManifestPath)
		skillTemplateData, err := fs.ReadFile(g.Assets, skillTemplatePath)
		if err != nil {
			return fmt.Errorf("failed to read skill template %s: %w", skillTemplatePath, err)
		}

		targetPattern := strings.TrimSuffix(skillManifestPath, ".template")
		tmplContent := string(skillTemplateData)

		// Generate files for each document
		for _, doc := range docs {
			filename, content, err := g.generateSkillFile(doc, tmplContent, targetPattern)
			if err != nil {
				return fmt.Errorf("failed to generate skill for %s: %w", doc.ID, err)
			}

			outPath := filepath.Join(cwd, filename)
			if err := kexfs.WriteFile(outPath, []byte(content)); err != nil {
				return err
			}
		}
	}
	return nil
}

type SkillTemplateData struct {
	SkillName   string
	Title       string
	Description string
	Body        string
}

func (g *Generator) generateSkillFile(doc *domain.Document, templateContent, outputPattern string) (string, string, error) {
	tmpl, err := template.New("skill").Parse(templateContent)
	if err != nil {
		return "", "", fmt.Errorf("parse template: %w", err)
	}

	filenameTmpl, err := template.New("filename").Parse(outputPattern)
	if err != nil {
		return "", "", fmt.Errorf("parse filename template: %w", err)
	}

	data := SkillTemplateData{
		SkillName:   doc.ID,
		Title:       doc.Title,
		Description: doc.Description,
		Body:        doc.Body,
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", "", fmt.Errorf("execute template: %w", err)
	}

	var filenameBuf bytes.Buffer
	if err := filenameTmpl.Execute(&filenameBuf, data); err != nil {
		return "", "", fmt.Errorf("execute filename template: %w", err)
	}

	return filenameBuf.String(), buf.String(), nil
}
