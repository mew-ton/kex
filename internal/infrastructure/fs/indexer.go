package fs

import (
	"fmt"
	"strings"

	"github.com/mew-ton/kex/internal/domain"
)

// Indexer manages the document collection and search index
type Indexer struct {
	Provider DocumentProvider
	BaseURL  string // Should this be removed? Provider has it.
	// We might need it if we're doing path resolution here, but Provider handles Paths in Schema.
	// Actually GetByID delegates to Provider.FetchContent.
	// So we don't need BaseURL here.
	IncludeDrafts bool                          // If true, draft documents are indexed
	Documents     map[string]*domain.Document   // ID -> Document
	Index         map[string][]*domain.Document // Keyword -> Documents
	Errors        []error                       // Validation errors found during load
	Schema        *IndexSchema                  // Unified Schema
}

// New creates a new Indexer
func New(provider DocumentProvider) *Indexer {
	return &Indexer{
		Provider:  provider,
		Documents: make(map[string]*domain.Document),
		Index:     make(map[string][]*domain.Document),
		Errors:    []error{},
	}
}

// Load scans the root directory and populates the index
// Dogfooding: Top-down decomposition
func (i *Indexer) Load() error {
	// 1. Load Schema from Provider
	schema, errs := i.Provider.Load()
	if schema == nil {
		if len(errs) > 0 {
			return fmt.Errorf("failed to load documents: %w", errs[0])
		}
		return fmt.Errorf("failed to load documents: unknown error")
	}

	// Store non-fatal errors
	if len(errs) > 0 {
		i.Errors = append(i.Errors, errs...)
	}

	i.Schema = schema

	// 2. Convert Schema to Domain Documents
	for _, sd := range schema.Documents {
		doc := &domain.Document{
			ID:          sd.ID,
			Title:       sd.Title,
			Description: sd.Description,
			Keywords:    sd.Keywords,
			Scopes:      sd.Scopes,
			Status:      domain.DocumentStatus(sd.Status),
			Path:        sd.Path,
		}

		// Map empty status to Adopted if missing?
		// If sd.Status is empty (from remote kex.json), imply Adopted.
		if doc.Status == "" {
			doc.Status = domain.StatusAdopted
		}

		i.addDocument(doc)
	}

	return nil
}

// Export generates the IndexSchema from current valid documents (Adopted Only)
func (i *Indexer) Export() (*IndexSchema, error) {
	schema := &IndexSchema{
		// GeneratedAt logic?
		Documents: []*DocumentSchema{},
	}

	for _, doc := range i.Documents {
		if doc.Status != domain.StatusAdopted {
			continue
		}

		// Use original path from Schema if available?
		// i.Documents path is relative/absolute depending on load?
		// LocalProvider: Path is relative to Root.
		// RemoteProvider: Path is relative or absolute.
		// For Export (Generate), we usually assume Local load.
		// We want paths relative to Root (so they look like 'contents/foo.md').

		schema.Documents = append(schema.Documents, &DocumentSchema{
			ID:          doc.ID,
			Title:       doc.Title,
			Description: doc.Description,
			Keywords:    doc.Keywords,
			Scopes:      doc.Scopes,
			// Status is implicitly adopted in kex.json output?
			// Or we can include it.
			// Task said "Status field removed (implict adopted)".
			// But I added it back to Schema.
			// If I added it back, I should populate it?
			// Or populate it as "adopted" explicitly.
			// Or omit it if omitempty?
			// If we filter only adopted, we can probably omit it if it matches default?
			// But explicitness is fine.
			Status: string(doc.Status),
			Path:   doc.Path,
		})
	}
	return schema, nil
}

func (i *Indexer) parseDocuments(paths []string) ([]*domain.Document, []error) {
	var docs []*domain.Document
	var errs []error

	for _, path := range paths {
		doc, err := domain.ParseDocument(path, "") // Root is no longer needed here
		if err != nil {
			errs = append(errs, fmt.Errorf("%s: %w", path, err))
			continue
		}
		// Skip drafts unless configured
		if !i.IncludeDrafts && doc.Status == domain.StatusDraft {
			continue
		}
		docs = append(docs, doc)
	}
	return docs, errs
}

func (i *Indexer) addDocument(doc *domain.Document) {
	// Add to ID map
	i.Documents[doc.ID] = doc

	// Add to Keyword index
	for _, keyword := range doc.Keywords {
		k := strings.ToLower(keyword)
		i.Index[k] = append(i.Index[k], doc)
	}
}

// Search returns documents matching the keywords and scopes
func (i *Indexer) Search(keywords []string, scopes []string) []*domain.Document {
	// 1. Find matches by keywords (OR logic)
	seen := make(map[string]struct{})
	var candidates []*domain.Document

	for _, keyword := range keywords {
		k := strings.ToLower(keyword)
		docs, ok := i.Index[k]
		if !ok {
			continue
		}
		for _, doc := range docs {
			if _, exists := seen[doc.ID]; !exists {
				seen[doc.ID] = struct{}{}
				candidates = append(candidates, doc)
			}
		}
	}

	// 2. Filter by scopes (Intersection logic)
	// If requestsScopes is empty, return all candidates (no filter).
	if len(scopes) == 0 {
		return candidates
	}

	var results []*domain.Document
	for _, doc := range candidates {
		if hasIntersection(doc.Scopes, scopes) {
			results = append(results, doc)
		}
	}

	return results
}

func (i *Indexer) GetAll() []*domain.Document {
	docs := make([]*domain.Document, 0, len(i.Documents))
	for _, doc := range i.Documents {
		docs = append(docs, doc)
	}
	return docs
}

func (i *Indexer) GetErrors() []error {
	return i.Errors
}

func (i *Indexer) GetByID(id string) (*domain.Document, bool) {
	doc, ok := i.Documents[id]
	if !ok {
		return nil, false
	}

	// Lazy Loading
	if doc.Body == "" {
		content, err := i.Provider.FetchContent(doc.Path)
		if err == nil {
			doc.Body = content
		} else {
			// Log error?
			fmt.Printf("Failed to fetch content for %s: %v\n", id, err)
		}
	}

	return doc, ok
}

func hasIntersection(a, b []string) bool {
	// Optimize for small slices
	for _, vA := range a {
		for _, vB := range b {
			if vA == vB {
				return true
			}
		}
	}
	return false
}
