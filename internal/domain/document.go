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
	ID          string         `yaml:"id"`
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

	doc, err := parseDocumentContent(path, content)
	if err != nil {
		return nil, err
	}

	// Derive scopes from path
	// Example path: contents/coding/typescript/foo.md
	// Scopes: ["coding", "typescript"]

	// Normalize path separators
	cleanPath := filepath.ToSlash(path)

	// Find where "contents/" ends (or rely on relative path from root if passed that way)
	// Currently path seems to be absolute or relative to cwd.
	// Let's assume we get a consistent path. Ideally ParseDocument receives a relative path from the index root.
	// But indexer passes full path.

	// We need to extract segments between root and filename.
	// Since we don't know the root here easily without changing signature,
	// let's look for known domains or assume standard structure.

	// Better approach: Use the relative path logic in Indexer or pass root here.
	// However, to keep it simple and stateless:
	// "examples/contents/coding/typescript/foo.md" -> split

	dirs := strings.Split(filepath.Dir(cleanPath), "/")
	var scopes []string

	// Heuristic: Collect all segments that are likely scopes.
	// We can skip "examples", "contents".
	for _, d := range dirs {
		if d == "." || d == "examples" || d == "contents" {
			continue
		}
		scopes = append(scopes, d)
	}
	doc.Scopes = scopes

	return doc, nil
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
