package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

type ChatResponse struct {
	Choices []struct {
		Message Message `json:"message"`
	} `json:"choices"`
}

func askAI(
	client *http.Client,
	apiURL, apiKey, model string,
	messages []Message,
) (string, error) {
	reqBody := ChatRequest{
		Model:    model,
		Messages: messages,
	}

	data, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest(
		http.MethodPost,
		apiURL,
		bytes.NewReader(data),
	)
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(io.LimitReader(resp.Body, 10<<20))
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf(
			"API error (%d): %s",
			resp.StatusCode,
			string(body),
		)
	}

	var result ChatResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return "", err
	}

	if len(result.Choices) == 0 {
		return "", fmt.Errorf("AI returned no choices")
	}

	return result.Choices[0].Message.Content, nil
}

func main() {
	// 加载 .env 中的环境变量
	if err := godotenv.Load(); err != nil {
		fmt.Println("Failed to load .env:", err)
		return
	}

	apiKey := os.Getenv("AI_API_KEY")
	apiURL := os.Getenv("AI_API_URL")
	model := os.Getenv("AI_MODEL")

	if apiKey == "" || apiURL == "" || model == "" {
		fmt.Println("Please set AI_API_KEY, AI_API_URL and AI_MODEL")
		return
	}

	client := &http.Client{
		Timeout: 120 * time.Second,
	}

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Buffer(make([]byte, 4096), 1024*1024)

	messages := []Message{
		{
			Role:    "system",
			Content: "You are a helpful coding assistant.",
		},
	}

	fmt.Println("AI CLI started. Type /exit to quit, /clear to reset.")

	for {
		fmt.Print("\nYou > ")

		if !scanner.Scan() {
			break
		}

		input := strings.TrimSpace(scanner.Text())

		switch input {
		case "":
			continue
		case "/exit":
			fmt.Println("Goodbye!")
			return
		case "/clear":
			messages = messages[:1]
			fmt.Println("Conversation cleared.")
			continue
		}

		messages = append(messages, Message{
			Role:    "user",
			Content: input,
		})

		answer, err := askAI(
			client, apiURL, apiKey, model, messages,
		)
		if err != nil {
			fmt.Println("Error:", err)
			messages = messages[:len(messages)-1]
			continue
		}

		fmt.Println("\nAI >", answer)

		messages = append(messages, Message{
			Role:    "assistant",
			Content: answer,
		})
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Input error:", err)
	}
}
