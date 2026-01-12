package fs

import (
	"fmt"
	"strconv"
	"strings"
)

// CompositeProvider aggregates multiple DocumentProviders into a single view.
// It manages path routing by prefixing document paths with the provider index.
type CompositeProvider struct {
	Providers []DocumentProvider
}

// NewCompositeProvider creates a new CompositeProvider
func NewCompositeProvider(providers []DocumentProvider) *CompositeProvider {
	return &CompositeProvider{
		Providers: providers,
	}
}

// Load aggregates schemas from all providers.
// Paths are modified to include the provider index (e.g. "0:path/to/doc.md")
// to ensure uniqueness and allow correct routing in FetchContent.
func (c *CompositeProvider) Load() (*IndexSchema, []error) {
	combinedSchema := &IndexSchema{
		Documents: []*DocumentSchema{},
	}
	var allErrors []error

	for i, p := range c.Providers {
		schema, errs := p.Load()
		if len(errs) > 0 {
			allErrors = append(allErrors, errs...)
		}
		// If schema is nil (fatal load error), we skip this provider but continue with others if possible?
		// Or should we fail hard? Standard behavior: return what we have, append errors.
		if schema == nil {
			continue
		}

		for _, doc := range schema.Documents {
			// Prefix path with provider index
			doc.Path = fmt.Sprintf("%d:%s", i, doc.Path)
			combinedSchema.Documents = append(combinedSchema.Documents, doc)
		}
	}

	return combinedSchema, allErrors
}

// FetchContent routes the request to the correct provider based on the path prefix.
func (c *CompositeProvider) FetchContent(path string) (string, error) {
	parts := strings.SplitN(path, ":", 2)
	if len(parts) != 2 {
		return "", fmt.Errorf("invalid composite path format: %s", path)
	}

	indexStr := parts[0]
	actualPath := parts[1]

	index, err := strconv.Atoi(indexStr)
	if err != nil {
		return "", fmt.Errorf("invalid provider index in path: %s", path)
	}

	if index < 0 || index >= len(c.Providers) {
		return "", fmt.Errorf("provider index out of range: %d", index)
	}

	return c.Providers[index].FetchContent(actualPath)
}
