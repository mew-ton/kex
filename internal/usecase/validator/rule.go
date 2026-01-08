package validator

import (
	"fmt"

	"github.com/mew-ton/kex/internal/domain"
)

// ValidationRule logic for checking document validty
type ValidationRule interface {
	Validate(doc *domain.Document) error
}

// TitleRequiredRule ensures Title is present
type TitleRequiredRule struct{}

func (r *TitleRequiredRule) Validate(doc *domain.Document) error {
	if doc.Title == "" {
		return fmt.Errorf("title is required")
	}
	return nil
}
