package main

import (
	"MiniAgent/agent"
	"MiniAgent/utils"
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/openai/openai-go/v3/responses"
)

func main() {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	cfg, _ := utils.LoadConfig("config/models.yaml")
	MyAgent, err := agent.NewAgent(*cfg, cfg.DefaultModel)
	if err != nil {
		fmt.Printf("DefaultModel: %q\n", cfg.DefaultModel)
		fmt.Printf("Models: %#v\n", cfg.Models)
		panic(err)
	}

	ctx := context.Background()
	scanner := bufio.NewScanner(os.Stdin)

	// 使用 SDK 官方类型保存历史消息
	history := responses.ResponseInputParam{
		responses.ResponseInputItemParamOfMessage(
			"You are a helpful assistant.",
			responses.EasyInputMessageRoleSystem,
		),
	}

	for {
		fmt.Print("You: ")
		if !scanner.Scan() {
			break
		}

		question := strings.TrimSpace(scanner.Text())
		if question == "/q" {
			break
		}
		if question == "/clear" {
			history = history[:1]
			continue
		}
		if question == "/model" {
			for {
				fmt.Println("Enter model name:")
				scanner.Scan()
				modelName := scanner.Text()
				err := MyAgent.SwitchModel(*cfg, modelName)
				if err != nil {
					fmt.Printf("Model switch error: %s, Please try again.\n", err)
				} else {
					break
				}
			}
			continue
		}
		if question == "" {
			continue
		}

		// 追加用户消息
		history = append(history,
			responses.ResponseInputItemParamOfMessage(
				question,
				responses.EasyInputMessageRoleUser,
			),
		)

		// 发送完整历史
		fmt.Print("AI: ")
		stream := MyAgent.Client.Responses.NewStreaming(
			ctx,
			responses.ResponseNewParams{
				Model: MyAgent.Model.Model,
				Input: responses.ResponseNewParamsInputUnion{
					OfInputItemList: history,
				},
			},
		)

		defer func() { _ = stream.Close() }()

		var completeText string

		for stream.Next() {
			data := stream.Current()
			fmt.Print(data.Delta)
			if data.JSON.Text.Valid() {
				fmt.Println()
				fmt.Println("==================")
				completeText = data.Text
				break
			}
		}

		if stream.Err() != nil {
			panic(stream.Err())
		}

		// 3. 获取并输出回复
		answer := completeText

		// 4. 保存 assistant 消息
		history = append(history,
			responses.ResponseInputItemParamOfMessage(
				answer,
				responses.EasyInputMessageRoleAssistant,
			),
		)
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("输入错误:", err)
	}
}
