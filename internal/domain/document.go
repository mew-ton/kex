package domain

// DocumentStatus represents the status of a document
type DocumentStatus string

const (
	StatusDraft   DocumentStatus = "draft"
	StatusAdopted DocumentStatus = "adopted"
)

// DocumentType represents the type of document (indicator or constraint)
type DocumentType string

const (
	TypeConstraint DocumentType = "constraint"
	TypeIndicator  DocumentType = "indicator"
)

// Document represents a single guideline document
type Document struct {
	ID          string         `yaml:"-"`
	Title       string         `yaml:"title"`
	Description string         `yaml:"description"`
	Type        DocumentType   `yaml:"type"`
	Keywords    []string       `yaml:"keywords"`
	Extensions  []string       `yaml:"extensions,omitempty"`
	Status      DocumentStatus `yaml:"status"`
	Sources     []struct {
		Name string `yaml:"name"`
		URL  string `yaml:"url"`
	} `yaml:"sources"`

	// Body content (markdown)
	Body string `yaml:"-"`

	// Metadata derived from file path
	Path string `yaml:"-"`

	Scopes []string `yaml:"-"` // Derived from directory structure
}
