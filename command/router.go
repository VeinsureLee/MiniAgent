package command

import (
	"MiniAgent/agent"
	"fmt"
)

const (
	CLEAR       = "/clear"
	QUIT        = "/q"
	ModelSwitch = "/model"
)

type Handler func(
	history *agent.History,
	chatAgent *agent.Agent,
	input chan string,
) (quit bool, err error)

func Router(cmd string) Handler {
	switch cmd {
	case CLEAR:
		return func(h *agent.History, a *agent.Agent, input chan string) (bool, error) {
			h.ClearMessages()
			return false, nil
		}

	case QUIT:
		return func(h *agent.History, a *agent.Agent, input chan string) (bool, error) {
			return true, nil
		}

	case ModelSwitch:
		return func(h *agent.History, a *agent.Agent, input chan string) (bool, error) {
			err := a.SwitchModel(input)
			return false, err
		}

	default:
		return func(h *agent.History, a *agent.Agent, input chan string) (bool, error) {
			fmt.Println("Unsupported command:", cmd)
			return false, nil
		}
	}
}
