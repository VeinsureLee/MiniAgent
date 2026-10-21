package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func Input(c chan<- string) {
	defer close(c)

	inputScanner := bufio.NewScanner(os.Stdin)

	for inputScanner.Scan() {
		if err := inputScanner.Err(); err != nil {
			fmt.Println("输入错误:", err)
		}
		userInput := strings.TrimSpace(inputScanner.Text())
		c <- userInput
	}
}
