package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/responses"
)

func main() {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	client := openai.NewClient(
		option.WithAPIKey(os.Getenv("DEEPSEEK_API_KEY")),
		option.WithBaseURL("https://api.deepseek.com"),
	)

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
		stream := client.Responses.NewStreaming(
			ctx,
			responses.ResponseNewParams{
				Model: "deepseek-flash",
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
