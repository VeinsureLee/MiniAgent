package agent

import (
	"fmt"
	"os"

	"MiniAgent/config"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

type Agent struct {
	ActiveModel config.ResolvedModelConfig
	APIClient   *openai.Client
}

func New(appConfig config.Config, modelAlias string) (*Agent, error) {
	resolvedModelConfig, err := appConfig.Resolve(modelAlias)
	if err != nil {
		return nil, fmt.Errorf("model alias %q not found", modelAlias)
	}

	apiKey := os.Getenv(resolvedModelConfig.APIKeyEnvVar)
	if apiKey == "" {
		return nil, fmt.Errorf(
			"environment variable %q is not set",
			resolvedModelConfig.APIKeyEnvVar,
		)
	}

	apiClient := openai.NewClient(
		option.WithAPIKey(apiKey),
		option.WithBaseURL(resolvedModelConfig.BaseURL),
	)

	return &Agent{
		ActiveModel: resolvedModelConfig,
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
