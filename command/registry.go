package command

import (
	"MiniAgent/agent"
	"fmt"
)

type Registry map[string]Command

func NewRegistry() Registry {
	r := Registry{
		"/clear": Command{
			Name:        "clear",
			Description: "clear all history",
			Handler:     handleClear,
		},
		"/quit": Command{
			Name:        "quit",
			Description: "quit the agent",
			Handler:     handleQuit,
		},
		"/model": Command{
			Name:        "model",
			Description: "switch a model",
			Handler:     handleSwitch,
		},
	}
	r["/help"] = Command{
		Name:        "help",
		Description: "show help",
		Handler:     newHelpHandler(r),
	}
	return r
}

func (r Registry) Router(cmd string) Handler {
	if command, ok := r[cmd]; ok {
		return command.Handler
	}

	return func(_ *agent.History, _ *agent.Agent, _ chan string) (bool, error) {
		fmt.Println("Unsupported command:", cmd)
		return false, nil
	}
}
