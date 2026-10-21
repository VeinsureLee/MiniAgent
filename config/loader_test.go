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

connections:
  deepseek:
    base_url: https://api.deepseek.com
    api_key_env: DEEPSEEK_API_KEY
    api_type: chat_completions

models:
  deepseek:
    id: deepseek-chat
    connection: deepseek
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

		Connections: map[string]ConnectionConfig{
			"deepseek": {
				BaseURL:      "https://api.deepseek.com",
				APIKeyEnvVar: "DEEPSEEK_API_KEY",
				APIType:      "chat_completions",
			},
		},

		ModelsByAlias: map[string]ModelConfig{
			"deepseek": {
				ModelID:    "deepseek-chat",
				Connection: "deepseek",
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

func TestResolve(t *testing.T) {
	cfg := Config{
		Connections: map[string]ConnectionConfig{
			"openai": {
				BaseURL:      "https://api.openai.com/v1",
				APIKeyEnvVar: "OPENAI_API_KEY",
				APIType:      "responses",
			},
		},
		ModelsByAlias: map[string]ModelConfig{
			"gpt": {
				ModelID:    "gpt-4.1",
				Connection: "openai",
			},
			"gpt-chat": {
				ModelID:    "gpt-4.1",
				Connection: "openai",
				APIType:    "chat_completions",
			},
			"broken": {
				ModelID:    "gpt-4.1",
				Connection: "missing",
			},
		},
	}

	tests := []struct {
		name    string
		alias   string
		apiType string
		wantErr bool
	}{
		{
			name:    "default protocol",
			alias:   "gpt",
			apiType: "responses",
		},
		{
			name:    "override protocol",
			alias:   "gpt-chat",
			apiType: "chat_completions",
		},
		{
			name:    "unknown model",
			alias:   "unknown",
			wantErr: true,
		},
		{
			name:    "unknown connection",
			alias:   "broken",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := cfg.Resolve(tt.alias)

			if (err != nil) != tt.wantErr {
				t.Fatalf("Resolve(%q) error = %v", tt.alias, err)
			}
			if tt.wantErr {
				return
			}

			if got.APIType != tt.apiType {
				t.Errorf("APIType = %q, want %q", got.APIType, tt.apiType)
			}
			if got.ModelID != "gpt-4.1" {
				t.Errorf("ModelID = %q, want gpt-4.1", got.ModelID)
			}
			if got.BaseURL != "https://api.openai.com/v1" {
				t.Errorf("BaseURL = %q, want OpenAI URL", got.BaseURL)
			}
			if got.APIKeyEnvVar != "OPENAI_API_KEY" {
				t.Errorf("APIKeyEnvVar = %q, want OPENAI_API_KEY", got.APIKeyEnvVar)
			}
		})
	}
}
