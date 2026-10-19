package agent

import (
	"MiniAgent/utils"
	"fmt"
	"os"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

type Agent struct {
	Model  utils.Model
	Client *openai.Client
}

func NewAgent(
	config utils.Config,
	modelName string,
) (*Agent, error) {
	model, ok := config.Models[modelName]
	if !ok {
		return nil, fmt.Errorf("model %q not found", modelName)
	}

	apiKey := os.Getenv(model.APIKeyEnv)
	if apiKey == "" {
		return nil, fmt.Errorf(
			"environment variable %q is not set",
			model.APIKeyEnv,
		)
	}

	client := openai.NewClient(
		option.WithAPIKey(apiKey),
		option.WithBaseURL(model.BaseURL),
	)

	return &Agent{
		Model:  model,
		Client: &client,
	}, nil
}

func (agent *Agent) SwitchModel(config utils.Config, newModelName string) error {
	next, err := NewAgent(config, newModelName)
	if err != nil {
		return err
	}

	*agent = *next
	return nil
}
