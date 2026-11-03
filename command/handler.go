package command

import (
	"MiniAgent/agent"
	"fmt"
	"sort"
)

type Handler func(
	history *agent.History,
	chatAgent *agent.Agent,
	input chan string,
) (quit bool, err error)

func handleClear(h *agent.History, a *agent.Agent, input chan string) (bool, error) {
	h.ClearMessages()
	return false, nil
}

func handleQuit(h *agent.History, a *agent.Agent, input chan string) (bool, error) {
	return true, nil
}

func handleSwitch(h *agent.History, a *agent.Agent, input chan string) (bool, error) {
	err := a.SwitchModel(input)
	return false, err
}

func newHelpHandler(r Registry) Handler {
	return func(_ *agent.History, _ *agent.Agent, _ chan string) (bool, error) {
		keys := make([]string, 0, len(r))
		for key := range r {
			keys = append(keys, key)
		}
		sort.Strings(keys)
		for _, key := range keys {
			cmd := r[key]
			fmt.Printf("%s\t%s\n", key, cmd.Description)
		}
		return false, nil
	}
}
