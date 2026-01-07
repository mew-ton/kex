package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Root        string `yaml:"root"`
	BaseURL     string `yaml:"baseURL"`
	RemoteToken string `yaml:"remoteToken"`
	Update      Update `yaml:"update"`
	Agent       Agent  `yaml:"agent"`
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
		Root:    "contents", // Default
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

	return config, nil
}
