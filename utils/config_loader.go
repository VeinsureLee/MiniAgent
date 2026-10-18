package utils

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
