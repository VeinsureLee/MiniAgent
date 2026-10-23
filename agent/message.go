package agent

import "github.com/openai/openai-go/v3/responses"

type Role string

const (
	RoleSystem    Role = "system"
	RoleUser      Role = "user"
	RoleAssistant Role = "assistant"
)

type Message struct {
	Role    Role
	Content string
}

func ToOpenAI(m Message) responses.ResponseInputItemUnionParam {
	return responses.ResponseInputItemParamOfMessage(
		m.Content,
		responses.EasyInputMessageRole(m.Role),
	)
}
