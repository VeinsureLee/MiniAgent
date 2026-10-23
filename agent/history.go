package agent

type History struct {
	Messages []Message
}

func NewHistory(systemPrompt string) History {
	msg := Message{
		Role:    RoleSystem,
		Content: systemPrompt,
	}
	return History{Messages: []Message{msg}}
}

func (h *History) ClearMessages() {
	h.Messages = h.Messages[:1]
}

func (h *History) PushMessage(message string, role Role) {
	h.Messages = append(h.Messages, Message{
		Role:    role,
		Content: message,
	})
}
