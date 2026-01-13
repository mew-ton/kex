package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Source      string       `yaml:"source"`
	References  []string     `yaml:"references,omitempty"`
	BaseURL     string       `yaml:"baseURL,omitempty"`
	RemoteToken string       `yaml:"remoteToken,omitempty"`
	Update      UpdateConfig `yaml:"update"`
	Logging     Logging      `yaml:"logging,omitempty"`
}

type Logging struct {
	Level string `yaml:"level,omitempty"`
	File  string `yaml:"file,omitempty"`
}

type UpdateConfig struct {
	Documents  map[string]string `yaml:"documents"`
	AiMcpRules AiMcpRules        `yaml:"ai-mcp-rules"`
	AiSkills   AiSkills          `yaml:"ai-skills"`
}

// AiMcpRules configuration
type AiMcpRules struct {
	Targets []string `yaml:"targets"` // List of agents: ["antigravity", "claude"]
	Scopes  []string `yaml:"scopes"`  // List of scopes: ["coding", "documentation"]
}

// AiSkills configuration
type AiSkills struct {
	Targets  []string `yaml:"targets"`  // List of agents: ["claude"]
	Keywords []string `yaml:"keywords"` // List of keywords to include
}

func Load(projectRoot string) (Config, error) {
	// 1. Set Defaults
	config := Config{
		BaseURL: "",
		Update: UpdateConfig{
			Documents: make(map[string]string),
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
		return config, err
	}

	if err := yaml.Unmarshal(data, &config); err != nil {
		return config, err
	}

	return config, nil
}
