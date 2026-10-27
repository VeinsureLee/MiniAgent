package config

import (
	"fmt"
	"os"
	"sync"

	"go.yaml.in/yaml/v4"
)

var (
	appConfig     *Config
	appConfigErr  error
	appConfigOnce sync.Once
)

func GetAppConfig() (*Config, error) {
	appConfigOnce.Do(func() {
		appConfig, appConfigErr = Load("config/models.yaml")
	})

	return appConfig, appConfigErr
}

// ConnectionConfig API 连接配置
type ConnectionConfig struct {
	BaseURL      string `yaml:"base_url"`
	APIKeyEnvVar string `yaml:"api_key_env"`
	APIType      string `yaml:"api_type"`
}

// ModelConfig 模型配置
type ModelConfig struct {
	ModelID    string `yaml:"id"`
	Connection string `yaml:"connection"`
	APIType    string `yaml:"api_type"` // 可选，覆盖连接默认值
}

// Config 完整配置
type Config struct {
	DefaultModelAlias string                      `yaml:"default_model"`
	Connections       map[string]ConnectionConfig `yaml:"connections"`
	ModelsByAlias     map[string]ModelConfig      `yaml:"models"`
}

// ResolvedModelConfig 解析后的模型配置
type ResolvedModelConfig struct {
	ModelID      string
	BaseURL      string
	APIKeyEnvVar string
	APIType      string
}

// Load 读取 YAML
func Load(filePath string) (*Config, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("read config %q: %w", filePath, err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parse config %q: %w", filePath, err)
	}

	return &cfg, nil
}

// Resolve 将模型别名解析为完整配置
func (c *Config) Resolve(alias string) (ResolvedModelConfig, error) {
	model, ok := c.ModelsByAlias[alias]
	if !ok {
		return ResolvedModelConfig{}, fmt.Errorf("unknown model: %s", alias)
	}

	conn, ok := c.Connections[model.Connection]
	if !ok {
		return ResolvedModelConfig{}, fmt.Errorf(
			"model %s references unknown connection: %s",
			alias, model.Connection,
		)
	}

	apiType := conn.APIType
	if model.APIType != "" {
		apiType = model.APIType
	}

	return ResolvedModelConfig{
		ModelID:      model.ModelID,
		BaseURL:      conn.BaseURL,
		APIKeyEnvVar: conn.APIKeyEnvVar,
		APIType:      apiType,
	}, nil
}
