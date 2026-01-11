package generator

import (
	"os"
	"path/filepath"
)

type UpdateContext struct {
	TargetPath   string
	TemplateData []byte
	Strategy     string
	Generator    *Generator
}

type UpdateStrategy interface {
	Apply(ctx UpdateContext) error
}

// OverwriteStrategy replaces the file
type OverwriteStrategy struct{}

func (s *OverwriteStrategy) Apply(ctx UpdateContext) error {
	if err := os.MkdirAll(filepath.Dir(ctx.TargetPath), 0755); err != nil {
		return err
	}
	return os.WriteFile(ctx.TargetPath, ctx.TemplateData, 0644)
}

// ResolveStrategy determines the strategy
func ResolveStrategy(strategyName string) UpdateStrategy {
	switch strategyName {
	case "overwrite":
		return &OverwriteStrategy{}
	case "ignore":
		return nil
	default:
		// Unknown -> Ignore
		return nil
	}
}
