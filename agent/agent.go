package agent

import (
	"fmt"
	"os"

	"MiniAgent/config"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

type Agent struct {
	ActiveModel config.ModelConfig
	APIClient   *openai.Client
}

func New(appConfig config.Config, modelAlias string) (*Agent, error) {
	modelConfig, ok := appConfig.ModelsByAlias[modelAlias]
	if !ok {
		return nil, fmt.Errorf("model alias %q not found", modelAlias)
	}

	apiKey := os.Getenv(modelConfig.APIKeyEnvVar)
	if apiKey == "" {
		return nil, fmt.Errorf(
			"environment variable %q is not set",
			modelConfig.APIKeyEnvVar,
		)
	}

	apiClient := openai.NewClient(
		option.WithAPIKey(apiKey),
		option.WithBaseURL(modelConfig.BaseURL),
	)

	return &Agent{
		ActiveModel: modelConfig,
		APIClient:   &apiClient,
	}, nil
}

func (a *Agent) SwitchModel(appConfig config.Config, modelAlias string) error {
	replacement, err := New(appConfig, modelAlias)
	if err != nil {
		return err
	}

	*a = *replacement
	return nil
}
