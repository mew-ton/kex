package domain

import (
	"fmt"
	"io"
	"os"
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
	ID          string         `yaml:"id"`
	Title       string         `yaml:"title"`
	Description string         `yaml:"description"`
	Keywords    []string       `yaml:"keywords"`
	Status      DocumentStatus `yaml:"status"`
	Source      struct {
		URL string `yaml:"url"`
	} `yaml:"source"`

	// Body content (markdown)
	Body string `yaml:"-"`

	// Metadata derived from file path
	Path       string `yaml:"-"`
	Language   string `yaml:"-"` // e.g. "typescript", "go" (from path)
	IsUniversal bool   `yaml:"-"` // true if under universal/
}

// ParseDocument reads a markdown file and parses its frontmatter
func ParseDocument(path string) (*Document, error) {
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

	return parseDocumentContent(path, content)
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

	doc := &Document{
		Status: StatusAdopted, // Default status
	}

	if err := yaml.Unmarshal([]byte(frontmatter), doc); err != nil {
		return nil, fmt.Errorf("failed to parse frontmatter: %w", err)
	}

	doc.Body = body
	doc.Path = path
	
	// Basic validation
	if doc.ID == "" {
		return nil, fmt.Errorf("id is required")
	}
	if doc.Title == "" {
		return nil, fmt.Errorf("title is required")
	}

	return doc, nil
}
