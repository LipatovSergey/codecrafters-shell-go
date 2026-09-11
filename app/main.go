package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("$ ")
		command, _ := reader.ReadString('\n')
		command = strings.TrimSpace(command)
		if command == "exit" {
			break
		}
		fmt.Print(command, ": command not found\n")
	}
}
