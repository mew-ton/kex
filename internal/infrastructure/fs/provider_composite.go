package fs

import (
	"fmt"
)

// CompositeProvider aggregates multiple DocumentProviders
type CompositeProvider struct {
	providers []DocumentProvider
	// Map from path to provider index for efficient fetching
	pathMap map[string]int
}

// NewCompositeProvider creates a new CompositeProvider
func NewCompositeProvider(providers []DocumentProvider) *CompositeProvider {
	return &CompositeProvider{
		providers: providers,
		pathMap:   make(map[string]int),
	}
}

// Load loads documents from all providers and merges them
func (c *CompositeProvider) Load() (*IndexSchema, []error) {
	unifiedSchema := &IndexSchema{
		Documents: []*DocumentSchema{},
	}
	var allErrors []error

	// Track IDs for collision detection
	// ID -> Path (or Source)
	seenIDs := make(map[string]string)

	for i, provider := range c.providers {
		schema, errs := provider.Load()
		if len(errs) > 0 {
			allErrors = append(allErrors, errs...)
		}
		if schema == nil {
			continue
		}

		for _, doc := range schema.Documents {
			// Collision Detection (Policy: Strict Uniqueness)
			if source, exists := seenIDs[doc.ID]; exists {
				allErrors = append(allErrors, fmt.Errorf("duplicate document ID '%s' found in '%s' and '%s'", doc.ID, source, doc.Path))
				continue
			}
			seenIDs[doc.ID] = doc.Path

			unifiedSchema.Documents = append(unifiedSchema.Documents, doc)
			// Map path to provider for retrieval
			c.pathMap[doc.Path] = i
		}
	}

	return unifiedSchema, allErrors
}

// FetchContent retrieves content from the appropriate provider
func (c *CompositeProvider) FetchContent(path string) (string, error) {
	indices, ok := c.pathMap[path]
	if !ok {
		// Fallback: try all providers if path not in map (should not happen if Load was successful)
		// But FetchContent might be called directly? Unlikely in current Indexer flow.
		// Let's iterate.
		for _, p := range c.providers {
			content, err := p.FetchContent(path)
			if err == nil {
				return content, nil
			}
		}
		return "", fmt.Errorf("document not found: %s", path)
	}

	return c.providers[indices].FetchContent(path)
}
