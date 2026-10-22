package agent

import (
	"context"
	"fmt"
	"os"
	"strings"

	"MiniAgent/config"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/responses"
)

type Agent struct {
	ActiveModel config.ResolvedModelConfig
	APIClient   *openai.Client
	AppConfig   *config.Config
}

func NewAgent(appConfig config.Config, modelAlias string) (*Agent, error) {
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
		AppConfig:   &appConfig,
	}, nil
}

func (chatAgent *Agent) SwitchModel(userInputChannel chan string) error {
	fmt.Println("Please enter your model alias: ")
	modelAlias, ok := <-userInputChannel
	if !ok {
		return fmt.Errorf("input closed while selecting a model")
	}

	replacement, err := NewAgent(*chatAgent.AppConfig, modelAlias)
	if err != nil {
		return err
	}
	*chatAgent = *replacement
	return nil
}

func (chatAgent *Agent) StreamReply(
	ctx context.Context,
	conversation History,
) (string, error) {
	fmt.Print("AI: ")
	defer fmt.Println()

	responseStream := chatAgent.APIClient.Responses.NewStreaming(
		ctx,
		responses.ResponseNewParams{
			Model: chatAgent.ActiveModel.ModelID,
			Input: responses.ResponseNewParamsInputUnion{
				OfInputItemList: conversation.Messages,
			},
		},
	)
	defer func() { _ = responseStream.Close() }()

	var assistantReply strings.Builder
	completed := false
	for responseStream.Next() {
		event := responseStream.Current()
		switch event.Type {
		case "response.output_text.delta", "response.refusal.delta":
			fmt.Print(event.Delta)
			assistantReply.WriteString(event.Delta)
		case "response.completed":
			completed = true
		case "response.failed":
			return "", fmt.Errorf("response failed (%s): %s", event.Response.Error.Code, event.Response.Error.Message)
		case "response.incomplete":
			return "", fmt.Errorf("response incomplete: %s", event.Response.IncompleteDetails.Reason)
		case "error":
			return "", fmt.Errorf("response error (%s): %s", event.Code, event.Message)
		}
	}
	if err := responseStream.Err(); err != nil {
		return "", fmt.Errorf("stream response: %w", err)
	}
	if !completed {
		return "", fmt.Errorf("response stream ended before completion")
	}
	if strings.TrimSpace(assistantReply.String()) == "" {
		return "", fmt.Errorf("model returned an empty reply")
	}
	return assistantReply.String(), nil
}
