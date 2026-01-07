package generator

import (
	"os"
	"path/filepath"
	"strings"

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

// SkipStrategy does nothing
type SkipStrategy struct{}

func (s *SkipStrategy) Apply(ctx UpdateContext) error {
	return nil
}

// MarkerUpdateStrategy updates content between markers
type MarkerUpdateStrategy struct{}

func (s *MarkerUpdateStrategy) Apply(ctx UpdateContext) error {
	data := ctx.TemplateData

	// Generate dynamic content if Agent Config is provided
	if ctx.AgentConfig != nil {
		dynamicData, err := ctx.Generator.generateAgentContent(ctx.AgentConfig)
		if err != nil {
			return err
		}
		data = dynamicData
	}

	return ctx.Generator.handleAgentMD(ctx.TargetPath, data)
}

// StrategyResolver determines the strategy
func ResolveStrategy(path string, strategies map[string]string) UpdateStrategy {
	strategyName := "skip"

	// Default Strategies
	if strings.Contains(path, "documentation/kex") {
		strategyName = "overwrite"
	}
	if filepath.Base(path) == "AGENTS.md" || filepath.Base(path) == "CLAUDE.md" {
		strategyName = "marker-update"
	}

	// Config Override
	for pattern, action := range strategies {
		matched, _ := filepath.Match(pattern, path)
		if matched {
			strategyName = action
		}
	}

	switch strategyName {
	case "overwrite":
		return &OverwriteStrategy{}
	case "marker-update":
		return &MarkerUpdateStrategy{}
	case "append":
		return &AppendStrategy{}
	default:
		return &SkipStrategy{}
	}
}
