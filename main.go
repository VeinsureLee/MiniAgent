package main

import (
	"context"
	"fmt"
	"os"

	"github.com/joho/godotenv"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/responses"
)

func main() {
	ctx := context.Background()
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	client := openai.NewClient(
		option.WithAPIKey(os.Getenv("DEEPSEEK_API_KEY")), // defaults to os.LookupEnv("OPENAI_API_KEY")
		option.WithBaseURL("https://api.deepseek.com"),
	)

	question := "什么是原神"

	resp, err := client.Responses.New(
		ctx, responses.ResponseNewParams{
			Input: responses.ResponseNewParamsInputUnion{OfString: openai.String(question)},
			Model: "deepseek-flash",
		},
	)

	if err != nil {
		panic(err)
	}

	fmt.Println(resp.OutputText())
}
