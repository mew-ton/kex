package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Root        string  `yaml:"root"`
	BaseURL     string  `yaml:"baseURL"`
	RemoteToken string  `yaml:"remoteToken"`
	Update      Update  `yaml:"update"`
	Agent       Agent   `yaml:"agent"`
	Logging     Logging `yaml:"logging"`
}

type Logging struct {
	Level string `yaml:"level"`
	File  string `yaml:"file"`
}

type Update struct {
	Strategies Strategies `yaml:"strategies"`
}

type Strategies struct {
	// Rules (General)
	AgentKexCoding        string `yaml:".agent/rules/kex-coding.md"`
	AgentKexDocumentation string `yaml:".agent/rules/kex-documentation.md"`

	// Rules (Claude)
	ClaudeKexCoding        string `yaml:".claude/rules/kex/follow-coding-rules.md"`
	ClaudeKexDocumentation string `yaml:".claude/rules/kex/follow-documentation-rules.md"`

	// System Documentation (Canonical paths starting with contents/)
	ChooseEffectiveKeywords   string `yaml:"contents/documentation/kex/choose-effective-keywords.md"`
	FollowDirectoryHierarchy  string `yaml:"contents/documentation/kex/follow-directory-hierarchy.md"`
	IncludeValidFrontmatter   string `yaml:"contents/documentation/kex/include-valid-frontmatter.md"`
	PrioritizeEnglishLanguage string `yaml:"contents/documentation/kex/prioritize-english-language.md"`
	UseImperativeFilenames    string `yaml:"contents/documentation/kex/use-imperative-filenames.md"`
	WriteConciseContent       string `yaml:"contents/documentation/kex/write-concise-content.md"`
}

type Agent struct {
	Type   string   `yaml:"type"`
	Scopes []string `yaml:"scopes"`
}

func Load(projectRoot string) (Config, error) {
	// 1. Set Defaults
	config := Config{
		Root:    "contents",
		BaseURL: "",
		Update: Update{
			Strategies: Strategies{
				// Default Strategies
				// Rules: Skip (Create Only) to preserve user edits
				AgentKexCoding:         "skip",
				AgentKexDocumentation:  "skip",
				ClaudeKexCoding:        "skip",
				ClaudeKexDocumentation: "skip",

				// System Docs: Overwrite to keep up to date
				ChooseEffectiveKeywords:   "overwrite",
				FollowDirectoryHierarchy:  "overwrite",
				IncludeValidFrontmatter:   "overwrite",
				PrioritizeEnglishLanguage: "overwrite",
				UseImperativeFilenames:    "overwrite",
				WriteConciseContent:       "overwrite",
			},
		},
	}

	configPath := filepath.Join(projectRoot, ".kex.yaml")
	data, err := os.ReadFile(configPath)
	if os.IsNotExist(err) {
		return config, nil
	}
	if err != nil {
		return config, err
	}

	if err := yaml.Unmarshal(data, &config); err != nil {
		// Use partial config even if error? No, fail.
		return config, err
	}

	// Re-apply defaults if empty (in case Unmarshal zeroes them out for missing keys)
	s := &config.Update.Strategies
	if s.AgentKexCoding == "" {
		s.AgentKexCoding = "skip"
	}
	if s.AgentKexDocumentation == "" {
		s.AgentKexDocumentation = "skip"
	}

	if s.ClaudeKexCoding == "" {
		s.ClaudeKexCoding = "skip"
	}
	if s.ClaudeKexDocumentation == "" {
		s.ClaudeKexDocumentation = "skip"
	}

	if s.ChooseEffectiveKeywords == "" {
		s.ChooseEffectiveKeywords = "overwrite"
	}
	if s.FollowDirectoryHierarchy == "" {
		s.FollowDirectoryHierarchy = "overwrite"
	}
	if s.IncludeValidFrontmatter == "" {
		s.IncludeValidFrontmatter = "overwrite"
	}
	if s.PrioritizeEnglishLanguage == "" {
		s.PrioritizeEnglishLanguage = "overwrite"
	}
	if s.UseImperativeFilenames == "" {
		s.UseImperativeFilenames = "overwrite"
	}
	if s.WriteConciseContent == "" {
		s.WriteConciseContent = "overwrite"
	}

	return config, nil
}

// StrategyFor returns the strategy for a given canonical path.
// This encapsulates the mapping between paths and struct fields.
func (s Strategies) StrategyFor(path string) string {
	switch path {
	// Rules (General)
	case ".agent/rules/kex-coding.md":
		return s.AgentKexCoding
	case ".agent/rules/kex-documentation.md":
		return s.AgentKexDocumentation

	// Rules (Claude)
	case ".claude/rules/kex/follow-coding-rules.md":
		return s.ClaudeKexCoding
	case ".claude/rules/kex/follow-documentation-rules.md":
		return s.ClaudeKexDocumentation

	// System Documentation
	case "contents/documentation/kex/choose-effective-keywords.md":
		return s.ChooseEffectiveKeywords
	case "contents/documentation/kex/follow-directory-hierarchy.md":
		return s.FollowDirectoryHierarchy
	case "contents/documentation/kex/include-valid-frontmatter.md":
		return s.IncludeValidFrontmatter
	case "contents/documentation/kex/prioritize-english-language.md":
		return s.PrioritizeEnglishLanguage
	case "contents/documentation/kex/use-imperative-filenames.md":
		return s.UseImperativeFilenames
	case "contents/documentation/kex/write-concise-content.md":
		return s.WriteConciseContent
	}
	return ""
}
