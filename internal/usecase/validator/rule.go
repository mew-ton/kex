package validator

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/mew-ton/kex/internal/domain"
)

// ValidationRule logic for checking document validty
type ValidationRule interface {
	Validate(doc *domain.Document) error
}

// IDRequiredRule ensures ID is present
type IDRequiredRule struct{}

func (r *IDRequiredRule) Validate(doc *domain.Document) error {
	if doc.ID == "" {
		return fmt.Errorf("id is required")
	}
	return nil
}

// TitleRequiredRule ensures Title is present
type TitleRequiredRule struct{}

func (r *TitleRequiredRule) Validate(doc *domain.Document) error {
	if doc.Title == "" {
		return fmt.Errorf("title is required")
	}
	return nil
}

// FilenameMatchRule ensures filename matches ID
type FilenameMatchRule struct{}

func (r *FilenameMatchRule) Validate(doc *domain.Document) error {
	filename := filepath.Base(doc.Path)
	ext := filepath.Ext(filename)
	basename := strings.TrimSuffix(filename, ext)

	if basename != doc.ID {
		return fmt.Errorf("filename must match id (filename: %s, id: %s)", filename, doc.ID)
	}
	return nil
}
