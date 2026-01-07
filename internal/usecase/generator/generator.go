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

// Strategy for file mapping
type FileMapper func(relPath string) (string, bool)

func GeneralMapper(relPath string) (string, bool) {
	return relPath, true
}

func ClaudeMapper(relPath string) (string, bool) {
	if filepath.Base(relPath) == "AGENTS.md" {
		return filepath.Join(filepath.Dir(relPath), "CLAUDE.md"), true
	}
	return relPath, true
}

func (g *Generator) Generate(cwd string, agentType AgentType, agentConfig *config.Agent) error {
	var mapper FileMapper
	switch agentType {
	case AgentTypeClaude:
		mapper = ClaudeMapper
	default:
		mapper = GeneralMapper
	}

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

		// Skip partials
		if strings.Contains(relPath, "partials") {
			return nil
		}

		// Apply strategy
		mappedPath, ok := mapper(relPath)
		if !ok {
			return nil
		}

		targetPath := filepath.Join(cwd, mappedPath)

		if d.IsDir() {
			return os.MkdirAll(targetPath, 0755)
		}

		data, err := fs.ReadFile(g.Templates, path)
		if err != nil {
			return err
		}

		// Special handling for AGENTS.md (or mapped CLAUDE.md) based on original filename
		if filepath.Base(relPath) == "AGENTS.md" {
			// Generate dynamic content if Agent Config is provided
			if agentConfig != nil {
				dynamicData, err := g.generateAgentContent(agentConfig)
				if err != nil {
					return err
				}
				return g.handleAgentMD(targetPath, dynamicData)
			}
			return g.handleAgentMD(targetPath, data)
		}

		// Don't overwrite if exists
		if _, err := os.Stat(targetPath); err == nil {
			return nil
		}

		return os.WriteFile(targetPath, data, 0644)
	})

	if err != nil {
		return fmt.Errorf("failed to extract templates: %w", err)
	}

	return nil
}

func (g *Generator) handleAgentMD(targetPath string, templateData []byte) error {
	// If file doesn't exist, create it
	if _, err := os.Stat(targetPath); os.IsNotExist(err) {
		return os.WriteFile(targetPath, templateData, 0644)
	}

	// Read existing file
	content, err := os.ReadFile(targetPath)
	if err != nil {
		return err
	}

	return g.updateWithMarkers(targetPath, content, templateData)
}

func (g *Generator) updateWithMarkers(targetPath string, currentContent, templateData []byte) error {
	const startMarker = "<!-- kex: auto-update start -->"
	const endMarker = "<!-- kex: auto-update end -->"

	contentStr := string(currentContent)
	tmplStr := string(templateData)

	// Helper to find range
	findRange := func(s, startM, endM string) (int, int) {
		start := strings.Index(s, startM)
		end := strings.Index(s, endM)
		if start != -1 && end != -1 && start < end {
			return start, end + len(endM)
		}
		return -1, -1
	}

	// Find markers in Current Content
	cStart, cEnd := findRange(contentStr, startMarker, endMarker)

	// Find markers in Template Content
	tStart, tEnd := findRange(tmplStr, startMarker, endMarker)

	// If template doesn't have markers, fallback to full replacement logic or error?
	// We assume template HAS markers (we just added them).
	if tStart == -1 || tEnd == -1 {
		// Template is missing markers. Fallback to append logic?
		// For now, if template has no markers, we assume it's not a marker-updatable file?
		// But handleAgentMD calls this specifically for AGENTS.md which we updated.
		return nil
	}

	newSegment := tmplStr[tStart:tEnd]
	var finalContent string

	if cStart != -1 && cEnd != -1 {
		// Replace existsing block
		before := contentStr[:cStart]
		after := contentStr[cEnd:]
		finalContent = before + newSegment + after
	} else {
		// Append to end
		if len(contentStr) > 0 && !strings.HasSuffix(contentStr, "\n") {
			finalContent = contentStr + "\n\n" + newSegment
		} else {
			finalContent = contentStr + "\n" + newSegment
		}
	}

	return os.WriteFile(targetPath, []byte(finalContent), 0644)
}

// Update updates the kex repository files based on configuration
func (g *Generator) Update(cwd string, agentType AgentType, config map[string]string, agentConfig *config.Agent) error {
	var mapper FileMapper
	switch agentType {
	case AgentTypeClaude:
		mapper = ClaudeMapper
	default:
		mapper = GeneralMapper
	}

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

		// Skip partials
		if strings.Contains(relPath, "partials") {
			return nil
		}

		mappedPath, ok := mapper(relPath)
		if !ok {
			return nil
		}

		targetPath := filepath.Join(cwd, mappedPath)

		// Execute Strategy
		if d.IsDir() {
			return nil // Dirs are created implicitly by WriteFile or MkdirAll
		}

		data, err := fs.ReadFile(g.Templates, path)
		if err != nil {
			return err
		}

		strategy := ResolveStrategy(mappedPath, config)
		ctx := UpdateContext{
			TargetPath:   targetPath,
			TemplateData: data,
			AgentConfig:  agentConfig,
			Generator:    g,
		}

		return strategy.Apply(ctx)
	})
}

// generateAgentContent assembles AGENTS.md content from partials
func (g *Generator) generateAgentContent(cfg *config.Agent) ([]byte, error) {
	var sb strings.Builder

	// Markers Start
	sb.WriteString("<!-- kex: auto-update start -->\n")

	// Header
	header, err := fs.ReadFile(g.Templates, "templates/partials/header.md")
	if err != nil {
		return nil, err
	}
	sb.Write(header)
	sb.WriteString("\n\n")

	// Scopes
	for _, scope := range cfg.Scopes {
		var partialFile string
		switch scope {
		case "coding":
			partialFile = "templates/partials/coding.md"
		case "documentation":
			partialFile = "templates/partials/documentation.md"
		}

		if partialFile != "" {
			content, err := fs.ReadFile(g.Templates, partialFile)
			if err != nil {
				return nil, err
			}
			sb.Write(content)
			sb.WriteString("\n\n")
		}
	}

	// Footer
	footer, err := fs.ReadFile(g.Templates, "templates/partials/footer.md")
	if err != nil {
		return nil, err
	}
	sb.Write(footer)
	sb.WriteString("\n")

	// Markers End
	sb.WriteString("<!-- kex: auto-update end -->\n")

	return []byte(sb.String()), nil
}
