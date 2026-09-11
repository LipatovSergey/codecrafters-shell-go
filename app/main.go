package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	builtInCommands := map[string]struct{}{
		"echo": {},
		"exit": {},
		"type": {},
	}
	for {
		fmt.Print("$ ")
		input, _ := reader.ReadString('\n')
		input = input[:len(input)-1]
		if input == "exit" {
			break
		}
		if strings.HasPrefix(input, "echo ") {
			fmt.Println(input[5:])
			continue
		}
		if strings.HasPrefix(input, "type ") {
			command := input[5:]
			if _, ok := builtInCommands[command]; ok {
				fmt.Println(command + " is a shell builtin")
			} else {
				fmt.Println(command + ": not found")
			}
			continue
		}
		fmt.Println(input + ": command not found")
	}
}
