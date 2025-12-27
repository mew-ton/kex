package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Root string `yaml:"root"`
}

func Load() (Config, error) {
	config := Config{
		Root: "contents", // Default
	}

	data, err := os.ReadFile(".kex.yaml")
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
