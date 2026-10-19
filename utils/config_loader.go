package utils

import (
	"os"

	"go.yaml.in/yaml/v4"
)

type Model struct {
	Model     string `yaml:"model"`
	APIKeyEnv string `yaml:"api_key_env"`
	BaseURL   string `yaml:"base_url"`
	APIType   string `yaml:"api_type"`
}

type Config struct {
	DefaultModel string           `yaml:"default_model"`
	Models       map[string]Model `yaml:"models"`
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}
