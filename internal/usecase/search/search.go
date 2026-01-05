package search

import (
	"github.com/mew-ton/kex/internal/domain"
	"path/filepath"
)

type UseCase struct {
	Repo domain.DocumentRepository
}

func New(repo domain.DocumentRepository) *UseCase {
	return &UseCase{Repo: repo}
}

type Result struct {
	Documents []*domain.Document
	Message   string
}

func (uc *UseCase) Execute(keywords []string, filePath string) Result {
	scopes := deriveScopes(filePath)
	docs := uc.Repo.Search(keywords, scopes)

	return Result{
		Documents: docs,
	}
}

func deriveScopes(path string) []string {
	if path == "" {
		return nil
	}
	ext := filepath.Ext(path)
	switch ext {
	// TypeScript / JavaScript
	case ".ts", ".tsx", ".js", ".jsx":
		// Implicitly include 'vcs' because code changes often imply version control ops
		return []string{"coding", "typescript", "javascript", "frontend", "vcs"}

	// Go
	case ".go":
		return []string{"coding", "go", "backend", "vcs"}

	// Markdown / Documentation
	case ".md", ".txt":
		return []string{"documentation", "vcs"} // text files also version controlled

	default:
		// Default for unknown code files
		return []string{"coding", "vcs"}
	}
}
