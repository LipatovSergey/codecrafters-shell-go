package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func isExecutable(command string) (string, bool) {
	path := os.Getenv("PATH")
	dirs := filepath.SplitList(path)
	for _, dir := range dirs {
		fullpath := filepath.Join(dir, command)
		info, err := os.Stat(fullpath)
		if err == nil && info.Mode().Perm()&0o111 != 0 {
			return fullpath, true
		}
	}
	return "", false
}

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
		if input == "pwd" {
			fmt.Println(os.Getwd())
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
			fullpath, found := isExecutable(command)
			if found {
				fmt.Println(command, "is", fullpath)
				continue
			}
			fmt.Println(command + ": not found")
			continue
		}
		inputFields := strings.Fields(input)
		programName := inputFields[0]
		_, found := isExecutable(programName)
		if found {
			programArgs := inputFields[1:]
			cmd := exec.Command(programName, programArgs...)
			cmd.Stdin = os.Stdin
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			cmd.Run()
			continue
		}
		fmt.Println(input + ": command not found")
	}
}
