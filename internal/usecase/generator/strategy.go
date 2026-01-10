package generator

import (
	"os"
	"path/filepath"

	"github.com/mew-ton/kex/internal/infrastructure/config"
)

type UpdateContext struct {
	TargetPath   string
	TemplateData []byte
	AgentConfig  *config.Agent
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

// AppendStrategy appends content to the file
type AppendStrategy struct{}

func (s *AppendStrategy) Apply(ctx UpdateContext) error {
	f, err := os.OpenFile(ctx.TargetPath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	if _, err := f.Write(ctx.TemplateData); err != nil {
		return err
	}
	return nil
}

// CreateStrategy creates the file if it doesn't exist, otherwise skips
type CreateStrategy struct{}

func (s *CreateStrategy) Apply(ctx UpdateContext) error {
	if _, err := os.Stat(ctx.TargetPath); os.IsNotExist(err) {
		if err := os.MkdirAll(filepath.Dir(ctx.TargetPath), 0755); err != nil {
			return err
		}
		return os.WriteFile(ctx.TargetPath, ctx.TemplateData, 0644)
	}
	return nil
}

// StrategyResolver determines the strategy
func ResolveStrategy(path string, strategies config.Strategies) UpdateStrategy {
	// Look up the strategy name using the config accessor
	// Path MUST be the canonical path (e.g. .agent/rules/..., contents/...)
	strategyName := strategies.StrategyFor(path)

	// If result is empty, fallback to overwrite (standard behavior)
	if strategyName == "" {
		strategyName = "overwrite"
	}

	switch strategyName {
	case "overwrite":
		return &OverwriteStrategy{}
	case "skip":
		// "skip" maps to CreateStrategy (Write if missing, Skip if exists)
		return &CreateStrategy{}
	default:
		// Unknown to Overwrite
		return &OverwriteStrategy{}
	}
}
