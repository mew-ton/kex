package fs

import "time"

// IndexSchema represents the structure of kex.json
type IndexSchema struct {
	GeneratedAt time.Time         `json:"generated_at"`
	Documents   []*DocumentSchema `json:"documents"`
}

// DocumentSchema represents a lightweight document entry in kex.json
type DocumentSchema struct {
	ID          string   `json:"id"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Keywords    []string `json:"keywords"`
	Extensions  []string `json:"extensions,omitempty"`
	Scopes      []string `json:"scopes"`
	Type        string   `json:"type,omitempty"`
	Status      string   `json:"status,omitempty"`
	Path        string   `json:"path"` // Relative path to markdown file
}
