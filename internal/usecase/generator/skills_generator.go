package generator

import (
	"bytes"
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/mew-ton/kex/internal/domain"
	"github.com/mew-ton/kex/internal/infrastructure/config"
	kexfs "github.com/mew-ton/kex/internal/infrastructure/fs"
)

type SkillsGenerator struct {
	config config.AiSkills
}

func NewSkillsGenerator(cfg config.AiSkills) *SkillsGenerator {
	return &SkillsGenerator{config: cfg}
}

// Generate finds documents matching keywords and returns a map of target path -> content
func (g *SkillsGenerator) Generate(rootDir, templateContent, outputPattern string) (map[string]string, error) {
	results := make(map[string]string)

	if len(g.config.Keywords) == 0 {
		return results, nil
	}

	tmpl, err := template.New("skill").Parse(templateContent)
	if err != nil {
		return nil, fmt.Errorf("failed to parse skill template: %w", err)
	}

	filenameTmpl, err := template.New("filename").Parse(outputPattern)
	if err != nil {
		return nil, fmt.Errorf("failed to parse filename template: %w", err)
	}

	err = filepath.WalkDir(rootDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if filepath.Ext(path) != ".md" {
			return nil
		}

		// Parse document to check keywords
		doc, err := kexfs.ParseDocument(path, rootDir)
		if err != nil {
			// Skip malformed documents
			return nil
		}

		if g.matchesKeywords(doc) {
			// Apply template
			var buf bytes.Buffer
			data := SkillTemplateData{
				SkillName:   doc.ID,
				Title:       doc.Title,
				Description: doc.Description,
				Body:        doc.Body,
			}
			if err := tmpl.Execute(&buf, data); err != nil {
				return fmt.Errorf("failed to execute template for %s: %w", path, err)
			}

			// Generate filename from pattern
			var filenameBuf bytes.Buffer
			if err := filenameTmpl.Execute(&filenameBuf, data); err != nil {
				return fmt.Errorf("failed to execute filename template for %s: %w", path, err)
			}
			filename := filenameBuf.String()

			results[filename] = buf.String()
		}

		return nil
	})

	return results, err
}

type SkillTemplateData struct {
	SkillName   string
	Title       string
	Description string
	Body        string
}

func (g *SkillsGenerator) matchesKeywords(doc *domain.Document) bool {
	// 1. Check document keywords
	for _, k := range doc.Keywords {
		for _, target := range g.config.Keywords {
			if strings.EqualFold(k, target) {
				return true
			}
		}
	}

	// 2. Check path/scopes (implicit keywords)
	for _, scope := range doc.Scopes {
		for _, target := range g.config.Keywords {
			if strings.EqualFold(scope, target) {
				return true
			}
		}
	}

	return false
}
