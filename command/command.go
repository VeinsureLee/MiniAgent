package command

import "strings"

type Command struct {
	Name        string
	Description string
	Handler     Handler
}

func IsCommand(cmd string) bool {
	return strings.HasPrefix(cmd, "/")
}
