package indexer

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"

	"kex/internal/domain"
)

// Indexer manages the document collection and search index
type Indexer struct {
	Root      string
	Documents map[string]*domain.Document   // ID -> Document
	Index     map[string][]*domain.Document // Keyword -> Documents
}

// New creates a new Indexer
func New(root string) *Indexer {
	return &Indexer{
		Root:      root,
		Documents: make(map[string]*domain.Document),
		Index:     make(map[string][]*domain.Document),
	}
}

// Load scans the root directory and populates the index
// Dogfooding: Top-down decomposition
func (i *Indexer) Load() error {
	// 1. Collect all markdown file paths
	paths, err := i.collectMarkdownFiles()
	if err != nil {
		return fmt.Errorf("failed to collect files: %w", err)
	}

	// 2. Parse valid documents
	docs, errs := i.parseDocuments(paths)
	if len(errs) > 0 {
		// In a real CLI we might want to log warnings,
		// but for now we just proceed with valid ones or return error?
		// Design says "Check" command validates. "Start" just loads valid ones.
		// We'll proceed with valid docs.
	}

	// 3. Register documents to internal maps
	for _, doc := range docs {
		i.addDocument(doc)
	}

	return nil
}

func (i *Indexer) collectMarkdownFiles() ([]string, error) {
	var paths []string
	err := filepath.WalkDir(i.Root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			// Skip hidden directories like .git
			if strings.HasPrefix(d.Name(), ".") && d.Name() != "." {
				return filepath.SkipDir
			}
			return nil
		}
		if filepath.Ext(path) == ".md" {
			paths = append(paths, path)
		}
		return nil
	})
	return paths, err
}

func (i *Indexer) parseDocuments(paths []string) ([]*domain.Document, []error) {
	var docs []*domain.Document
	var errs []error

	for _, path := range paths {
		doc, err := domain.ParseDocument(path)
		if err != nil {
			errs = append(errs, fmt.Errorf("%s: %w", path, err))
			continue
		}
		// Skip drafts
		if doc.Status == domain.StatusDraft {
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

// Search returns documents matching the keywords
func (i *Indexer) Search(keywords []string) []*domain.Document {
	// Simple OR search for now
	// To improve: scoring, AND search, etc.
	seen := make(map[string]struct{})
	var results []*domain.Document

	for _, keyword := range keywords {
		k := strings.ToLower(keyword)
		docs, ok := i.Index[k]
		if !ok {
			continue
		}
		for _, doc := range docs {
			if _, exists := seen[doc.ID]; !exists {
				seen[doc.ID] = struct{}{}
				results = append(results, doc)
			}
		}
	}
	return results
}
