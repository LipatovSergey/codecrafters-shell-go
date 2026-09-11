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
		fmt.Print(command, ": command not found\n")
	}
}
