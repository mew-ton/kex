package fs

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/mew-ton/kex/internal/domain"
	"github.com/mew-ton/kex/internal/infrastructure/logger"
)

type LocalProvider struct {
	Root    string
	Logger  logger.Logger
	BaseURL string // Optional, for generating absolute URLs if needed (though mostly for generate)
	// Actually LocalProvider for 'start' just needs Root.
	// schema generation logic uses BaseURL but that might be separate?
	// Let's keep it simple.
}

func NewLocalProvider(root string, logger logger.Logger) *LocalProvider {
	return &LocalProvider{Root: root, Logger: logger}
}

func (l *LocalProvider) Load() (*IndexSchema, []error) {
	schema := &IndexSchema{
		Documents: []*DocumentSchema{},
	}
	var errs []error

	// Logic extracted from original Indexer.Load
	paths, err := l.collectMarkdownFiles()
	if err != nil {
		return nil, []error{err}
	}

	for _, path := range paths {
		doc, err := domain.ParseDocument(path, l.Root)
		if err != nil {
			// We might want to handle errors differently (collect them?),
			// but Provider interface signature simply returns error for Load.
			// The original Indexer collected errors.
			// Let's modify behavior: Skip invalid files but log?
			// Or just return minimal schema and let Indexer validate?
			// ParseDocument returns error if frontmatter is invalid.
			// We should probably skip invalid docs here or return error.
			// Original logic: "i.Errors = append..."
			// Let's skip and log formatted error to stderr for now?
			// Better: return what we can.
			errs = append(errs, err)
			continue
		}

		// Default to Adopted if status is missing in local files
		if doc.Status == "" {
			doc.Status = domain.StatusAdopted
		}

		relPath, _ := filepath.Rel(l.Root, path)

		schema.Documents = append(schema.Documents, &DocumentSchema{
			ID:          doc.ID,
			Title:       doc.Title,
			Description: doc.Description,
			Keywords:    doc.Keywords,
			Scopes:      doc.Scopes,
			Status:      string(doc.Status),
			Path:        relPath, // Relative to Root
		})
	}

	return schema, errs
}

func (l *LocalProvider) FetchContent(path string) (string, error) {
	fullPath := filepath.Join(l.Root, path)

	if l.Logger != nil {
		l.Logger.Info("[Filesystem] Read: %s", fullPath)
	}

	content, err := os.ReadFile(fullPath)
	if err != nil {
		return "", err
	}
	// Parse body from content (remove frontmatter)
	sContent := string(content)
	parts := strings.SplitN(sContent, "\n---\n", 2)
	if len(parts) >= 2 {
		return parts[1], nil
	}
	return sContent, nil
}

func (l *LocalProvider) collectMarkdownFiles() ([]string, error) {
	var paths []string
	err := filepath.WalkDir(l.Root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if filepath.Ext(path) == ".md" {
			paths = append(paths, path)
		}
		return nil
	})
	return paths, err
}
