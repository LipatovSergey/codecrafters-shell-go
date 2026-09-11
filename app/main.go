package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	path := os.Getenv("PATH")
	dirs := filepath.SplitList(path)
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
				continue
			}
			for _, dir := range dirs {
				fullpath := filepath.Join(dir, command)
				info, err := os.Stat(fullpath)
				if err != nil {
					continue
				}
				if info.Mode().Perm()&0o111 == 0 {
					continue
				}
				fmt.Println(command, "is", fullpath)
				break
			}
			fmt.Println(command + ": not found")
			continue
		}
		fmt.Println(input + ": command not found")
	}
}
