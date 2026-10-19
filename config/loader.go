package config

import (
	"fmt"
	"os"

	"go.yaml.in/yaml/v4"
)

type ModelConfig struct {
	ModelID      string `yaml:"model"`
	APIKeyEnvVar string `yaml:"api_key_env"`
	BaseURL      string `yaml:"base_url"`
	APIType      string `yaml:"api_type"`
}

type Config struct {
	DefaultModelAlias string                 `yaml:"default_model"`
	ModelsByAlias     map[string]ModelConfig `yaml:"models"`
}

func Load(filePath string) (*Config, error) {
	configYAML, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("read config %q: %w", filePath, err)
	}

	var appConfig Config
	if err := yaml.Unmarshal(configYAML, &appConfig); err != nil {
		return nil, fmt.Errorf("parse config %q: %w", filePath, err)
	}

	return &appConfig, nil
}
