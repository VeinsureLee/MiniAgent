package agent

import "github.com/openai/openai-go/v3/responses"

type History struct {
	Messages responses.ResponseInputParam
}

func NewHistory(systemPrompt string) History {
	conversation := responses.ResponseInputParam{
		responses.ResponseInputItemParamOfMessage(
			systemPrompt,
			responses.EasyInputMessageRoleSystem,
		),
	}
	return History{Messages: conversation}
}

func (h *History) ClearMessages() {
	h.Messages = h.Messages[:1]
}

func (h *History) PushMessage(message string, role responses.EasyInputMessageRole) {
	h.Messages = append(h.Messages,
		responses.ResponseInputItemParamOfMessage(
			message,
			role,
		),
	)
}
