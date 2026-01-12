package fs

import (
	"fmt"
	"testing"
)

// MockProvider is a simple mock for DocumentProvider
type MockProvider struct {
	Documents []*DocumentSchema
	Content   map[string]string
	Err       error
}

func (m *MockProvider) Load() (*IndexSchema, []error) {
	if m.Err != nil {
		return nil, []error{m.Err}
	}
	return &IndexSchema{Documents: m.Documents}, nil
}

func (m *MockProvider) FetchContent(path string) (string, error) {
	if content, ok := m.Content[path]; ok {
		return content, nil
	}
	return "", fmt.Errorf("not found")
}

func TestCompositeProvider_Load(t *testing.T) {
	p1 := &MockProvider{
		Documents: []*DocumentSchema{
			{ID: "doc1", Path: "doc1.md"},
		},
	}
	p2 := &MockProvider{
		Documents: []*DocumentSchema{
			{ID: "doc2", Path: "sub/doc2.md"},
		},
	}

	c := NewCompositeProvider([]DocumentProvider{p1, p2})
	schema, errs := c.Load()

	if len(errs) != 0 {
		t.Errorf("expected no errors, got %v", errs)
	}

	if len(schema.Documents) != 2 {
		t.Errorf("expected 2 documents, got %d", len(schema.Documents))
	}

	// Verify Path Prefixing
	foundDoc1 := false
	foundDoc2 := false
	for _, doc := range schema.Documents {
		if doc.ID == "doc1" {
			if doc.Path != "0:doc1.md" {
				t.Errorf("expected doc1 path '0:doc1.md', got '%s'", doc.Path)
			}
			foundDoc1 = true
		}
		if doc.ID == "doc2" {
			if doc.Path != "1:sub/doc2.md" {
				t.Errorf("expected doc2 path '1:sub/doc2.md', got '%s'", doc.Path)
			}
			foundDoc2 = true
		}
	}

	if !foundDoc1 || !foundDoc2 {
		t.Error("missing documents in composite schema")
	}
}

func TestCompositeProvider_FetchContent(t *testing.T) {
	p1 := &MockProvider{
		Content: map[string]string{"doc1.md": "content1"},
	}
	p2 := &MockProvider{
		Content: map[string]string{"sub/doc2.md": "content2"},
	}

	c := NewCompositeProvider([]DocumentProvider{p1, p2})

	tests := []struct {
		name    string
		path    string
		want    string
		wantErr bool
	}{
		{"Provider 0", "0:doc1.md", "content1", false},
		{"Provider 1", "1:sub/doc2.md", "content2", false},
		{"Invalid Format", "doc1.md", "", true},
		{"Invalid Index Type", "a:doc1.md", "", true},
		{"Index Out of Range", "99:doc1.md", "", true},
		{"Not Found", "0:missing.md", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := c.FetchContent(tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("CompositeProvider.FetchContent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("CompositeProvider.FetchContent() = %v, want %v", got, tt.want)
			}
		})
	}
}
