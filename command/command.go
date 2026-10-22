package command

import "strings"

func IsCommand(cmd string) bool {
	return strings.HasPrefix(cmd, "/")
}
