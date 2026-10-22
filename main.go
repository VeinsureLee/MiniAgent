package main

import (
	"MiniAgent/agent"
	"MiniAgent/cli"
	"MiniAgent/config"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	appConfig, err := config.Load("config/models.yaml")
	if err != nil {
		panic(err)
	}
	chatAgent, err := agent.NewAgent(*appConfig, appConfig.DefaultModelAlias)
	if err != nil {
		panic(err)
	}

	userInputChannel := make(chan string)
	go cli.Input(userInputChannel)

	cli.Run(userInputChannel, chatAgent)
}
