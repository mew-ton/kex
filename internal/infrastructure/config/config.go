package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Source      string  `yaml:"source"` // Path to document source
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
	Strategies map[string]string `yaml:"strategies"`
}

type Agent struct {
	Type   string   `yaml:"type"`
	Scopes []string `yaml:"scopes"`
}

func Load(projectRoot string) (Config, error) {
	config := Config{
		Source:  "contents", // Default
		BaseURL: "",
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

	// Fallback if source empty
	if config.Source == "" {
		config.Source = "contents"
	}

	return config, nil
}
