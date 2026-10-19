package config

import (
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestLoad(t *testing.T) {
	filePath := filepath.Join(t.TempDir(), "models.yaml")
	configYAML := `default_model: deepseek
models:
  deepseek:
    model: deepseek-chat
    api_key_env: DEEPSEEK_API_KEY
    base_url: https://api.deepseek.com
    api_type: chat_completions
`
	if err := os.WriteFile(filePath, []byte(configYAML), 0600); err != nil {
		t.Fatal(err)
	}

	appConfig, err := Load(filePath)
	if err != nil {
		t.Fatal(err)
	}

	want := Config{
		DefaultModelAlias: "deepseek",
		ModelsByAlias: map[string]ModelConfig{
			"deepseek": {
				ModelID:      "deepseek-chat",
				APIKeyEnvVar: "DEEPSEEK_API_KEY",
				BaseURL:      "https://api.deepseek.com",
				APIType:      "chat_completions",
			},
		},
	}
	if !reflect.DeepEqual(*appConfig, want) {
		t.Errorf("Load() = %#v, want %#v", *appConfig, want)
	}
}

func TestLoadMissingFile(t *testing.T) {
	filePath := filepath.Join(t.TempDir(), "missing.yaml")
	appConfig, err := Load(filePath)
	if !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("Load() error = %v, want file-not-found error", err)
	}
	if appConfig != nil {
		t.Errorf("Load() config = %#v, want nil", appConfig)
	}
}

func TestLoadInvalidYAML(t *testing.T) {
	filePath := filepath.Join(t.TempDir(), "invalid.yaml")
	if err := os.WriteFile(filePath, []byte("models: ["), 0600); err != nil {
		t.Fatal(err)
	}
	appConfig, err := Load(filePath)
	if err == nil {
		t.Fatal("Load() succeeded for invalid YAML")
	}
	if appConfig != nil {
		t.Errorf("Load() config = %#v, want nil", appConfig)
	}
}
