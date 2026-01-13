package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

func Save(projectRoot string, cfg Config) error {
	data, err := yaml.Marshal(&cfg)
	if err != nil {
		return err
	}

	configPath := filepath.Join(projectRoot, ".kex.yaml")
	return os.WriteFile(configPath, data, 0644)
}
