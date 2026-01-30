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
	Documents map[string]string `yaml:"documents"`
	Ai        AiConfig          `yaml:"ai-skills"`
}

// AiConfig configuration
type AiConfig struct {
	Targets  []string `yaml:"targets"` // List of agents: ["antigravity", "cursor", "claude"]
	Keywords []string `yaml:"keywords,omitempty"`
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

// FromCLI creates a Config object solely from CLI arguments.
// This is used when the user provides references directly via the command line.
func FromCLI(references []string, token string) Config {
	return Config{
		References:  references,
		RemoteToken: token,
		// Update strategy defaults are not populated here as this is an ephemeral runtime config
		Update: UpdateConfig{
			Documents: make(map[string]string),
		},
	}
}
