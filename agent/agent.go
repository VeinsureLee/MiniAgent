package agent

import (
	"context"
	"fmt"
	"os"

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
		AppConfig:   &appConfig,
	}, nil
}

func (chatAgent *Agent) SwitchModel(userInputChannel chan string) error {
	for {
		fmt.Println("Please enter your model alias: ")
		modelAlias := <-userInputChannel

		replacement, err := New(*chatAgent.AppConfig, modelAlias)
		if err != nil {
			return err
		}
		*chatAgent = *replacement
		return nil
	}
}

func (chatAgent *Agent) StreamReply(
	ctx context.Context,
	conversation responses.ResponseInputParam,
) (string, error) {
	fmt.Print("AI: ")
	responseStream := chatAgent.APIClient.Responses.NewStreaming(
		ctx,
		responses.ResponseNewParams{
			Model: chatAgent.ActiveModel.ModelID,
			Input: responses.ResponseNewParamsInputUnion{
				OfInputItemList: conversation,
			},
		},
	)
	defer func() { _ = responseStream.Close() }()

	var assistantReply string
	for responseStream.Next() {
		event := responseStream.Current()
		fmt.Print(event.Delta)
		if event.JSON.Text.Valid() {
			fmt.Println()
			fmt.Println("==================")
			assistantReply = event.Text
			break
		}
	}
	if err := responseStream.Err(); err != nil {
		return "", err
	}
	return assistantReply, nil
}

func (chatAgent *Agent) Run(userInputChannel chan string) {
	ctx := context.Background()

	conversation := responses.ResponseInputParam{
		responses.ResponseInputItemParamOfMessage(
			"You are a helpful assistant.",
			responses.EasyInputMessageRoleSystem,
		),
	}

	for {
		fmt.Print("You: ")

		userInput, ok := <-userInputChannel
		if !ok {
			fmt.Println("You are not input channel.")
		}

		if userInput == "/q" {
			break
		}
		if userInput == "/clear" {
			conversation = conversation[:1]
			continue
		}
		if userInput == "/model" {
			for err := chatAgent.SwitchModel(userInputChannel); err != nil; {
				fmt.Println("Error: ", err)
				err = chatAgent.SwitchModel(userInputChannel)
			}
			continue
		}
		if userInput == "" {
			continue
		}

		conversation = append(conversation,
			responses.ResponseInputItemParamOfMessage(
				userInput,
				responses.EasyInputMessageRoleUser,
			),
		)

		assistantReply, err := chatAgent.StreamReply(ctx, conversation)
		if err != nil {
			panic(err)
		}

		conversation = append(conversation,
			responses.ResponseInputItemParamOfMessage(
				assistantReply,
				responses.EasyInputMessageRoleAssistant,
			),
		)
	}
}
