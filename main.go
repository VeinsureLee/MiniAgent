package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"MiniAgent/agent"
	"MiniAgent/config"

	"github.com/joho/godotenv"
	"github.com/openai/openai-go/v3/responses"
)

func main() {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	appConfig, err := config.Load("config/models.yaml")
	if err != nil {
		panic(err)
	}
	chatAgent, err := agent.New(*appConfig, appConfig.DefaultModelAlias)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	inputScanner := bufio.NewScanner(os.Stdin)

	conversation := responses.ResponseInputParam{
		responses.ResponseInputItemParamOfMessage(
			"You are a helpful assistant.",
			responses.EasyInputMessageRoleSystem,
		),
	}

	for {
		fmt.Print("You: ")
		if !inputScanner.Scan() {
			break
		}

		userInput := strings.TrimSpace(inputScanner.Text())
		if userInput == "/q" {
			break
		}
		if userInput == "/clear" {
			conversation = conversation[:1]
			continue
		}
		if userInput == "/model" {
			promptModelSwitch(inputScanner, *appConfig, chatAgent)
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

		assistantReply, err := streamReply(ctx, chatAgent, conversation)
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

	if err := inputScanner.Err(); err != nil {
		fmt.Println("输入错误:", err)
	}
}

func promptModelSwitch(inputScanner *bufio.Scanner, appConfig config.Config, chatAgent *agent.Agent) {
	for {
		fmt.Println("Enter model alias:")
		if !inputScanner.Scan() {
			return
		}

		modelAlias := strings.TrimSpace(inputScanner.Text())
		if modelAlias == "" {
			continue
		}
		if err := chatAgent.SwitchModel(appConfig, modelAlias); err != nil {
			fmt.Printf("Model switch error: %s, Please try again.\n", err)
			continue
		}
		return
	}
}

func streamReply(ctx context.Context, chatAgent *agent.Agent, conversation responses.ResponseInputParam) (string, error) {
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
