package domain

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// DocumentStatus represents the status of a document
type DocumentStatus string

const (
	StatusDraft   DocumentStatus = "draft"
	StatusAdopted DocumentStatus = "adopted"
)

// Document represents a single guideline document
type Document struct {
	ID          string         `yaml:"-"`
	Title       string         `yaml:"title"`
	Description string         `yaml:"description"`
	Keywords    []string       `yaml:"keywords"`
	Status      DocumentStatus `yaml:"status"`
	Sources     []struct {
		Name string `yaml:"name"`
		URL  string `yaml:"url"`
	} `yaml:"sources"`

	// Body content (markdown)
	Body string `yaml:"-"`

	// Metadata derived from file path
	Path     string   `yaml:"-"`
	Language string   `yaml:"-"` // Deprecated: derived from Scopes logic if needed
	Scopes   []string `yaml:"-"` // Derived from directory structure
}

// ParseDocument reads a markdown file and parses its frontmatter
func ParseDocument(path, root string) (*Document, error) {
	// Read file content
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer f.Close()

	content, err := io.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	doc, err := parseDocumentContent(path, content)
	if err != nil {
		return nil, err
	}

	doc.Scopes, err = deriveScopes(path, root)
	if err != nil {
		// Fallback or just ignore error? deriveScopes handles fallback internally if needed
		// But here deriveScopes returns valid slice or empty.
		// Let's make deriveScopes return pure logic.
	}

	doc.ID = generateID(path, doc.Scopes)

	return doc, nil
}

func deriveScopes(path, root string) ([]string, error) {
	rel, err := filepath.Rel(root, path)
	if err != nil {
		return nil, err
	}

	dirs := strings.Split(filepath.Dir(rel), string(filepath.Separator))
	var scopes []string

	for _, d := range dirs {
		if d == "." || d == "" {
			continue
		}
		scopes = append(scopes, d)
	}
	return scopes, nil
}

func generateID(path string, scopes []string) string {
	filename := filepath.Base(path)
	ext := filepath.Ext(filename)
	basename := strings.TrimSuffix(filename, ext)

	if len(scopes) > 0 {
		return strings.Join(scopes, ".") + "." + basename
	}
	return basename
}

func parseDocumentContent(path string, content []byte) (*Document, error) {
	// Split frontmatter and body
	// We assume frontmatter is at the top, enclosed by ---
	sContent := string(content)
	if !strings.HasPrefix(sContent, "---\n") {
		return nil, fmt.Errorf("missing frontmatter")
	}

	parts := strings.SplitN(sContent, "\n---\n", 2)
	if len(parts) < 2 {
		return nil, fmt.Errorf("invalid frontmatter format")
	}

	frontmatter := parts[0][4:] // remove leading ---\n
	body := parts[1]

	doc := &Document{}

	if err := yaml.Unmarshal([]byte(frontmatter), doc); err != nil {
		return nil, fmt.Errorf("failed to parse frontmatter: %w", err)
	}

	doc.Body = body
	doc.Path = path

	// Basic validation
	if doc.Title == "" {
		return nil, fmt.Errorf("title is required")
	}

	return doc, nil
}
