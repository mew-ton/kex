package domain

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
