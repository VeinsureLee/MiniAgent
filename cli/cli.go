package cli

import (
	"MiniAgent/agent"
	"MiniAgent/command"
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/openai/openai-go/v3"
)

func Input(c chan<- string) {
	defer close(c)

	inputScanner := bufio.NewScanner(os.Stdin)

	for inputScanner.Scan() {
		userInput := strings.TrimSpace(inputScanner.Text())
		c <- userInput
	}
	if err := inputScanner.Err(); err != nil {
		fmt.Println("输入错误:", err)
	}
}

func Run(userInputChannel chan string, chatAgent *agent.Agent) {
	ctx := context.Background()

	history := agent.NewHistory("You are a helpful assistant.")
	for {
		fmt.Print("You: ")

		userInput, ok := <-userInputChannel
		if !ok {
			fmt.Println()
			return
		}
		userInput = strings.TrimSpace(userInput)

		if command.IsCommand(userInput) {
			handler := command.Router(userInput)
			quit, err := handler(
				&history,
				chatAgent,
				userInputChannel,
			)

			if err != nil {
				fmt.Println("Error:", err)
			}

			if quit {
				return
			}
			continue
		}

		messageCount := len(history.Messages)
		history.PushMessage(userInput, agent.RoleUser)

		assistantReply, err := chatAgent.StreamReply(ctx, history)
		if err != nil {
			history.Messages = history.Messages[:messageCount]
			fmt.Println("Error:", err)
			var apiErr *openai.Error
			if errors.As(err, &apiErr) && apiErr.Message != "" {
				fmt.Println("Details:", apiErr.Message)
			}
			continue
		}

		history.PushMessage(assistantReply, agent.RoleAssistant)
		fmt.Println("==================")
	}
}
