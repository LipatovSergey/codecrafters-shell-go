package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type builtinFunc func([]string)

type Shell struct {
	reader   *bufio.Reader
	builtins map[string]builtinFunc
}

func NewShell() *Shell {
	shell := &Shell{
		reader: bufio.NewReader(os.Stdin),
	}

	shell.builtins = map[string]builtinFunc{
		"echo": shell.echoCommand,
		"pwd":  shell.pwdCommand,
		"type": shell.typeCommand,
		"cd":   shell.cdCommand,
	}

	return shell
}

func (s *Shell) Run() {
	for {
		fmt.Print("$ ")

		input, err := s.reader.ReadString('\n')
		if err != nil {
			return
		}

		input = strings.TrimSuffix(input, "\n")

		command, args := parseInput(input)

		if command == "" {
			continue
		}

		if command == "exit" {
			return
		}

		if handler, ok := s.builtins[command]; ok {
			handler(args)
			continue
		}

		s.runExternal(command, args)
	}
}

func parseInput(input string) (string, []string) {
	fields := strings.Fields(input)

	if len(fields) == 0 {
		return "", nil
	}

	command := fields[0]
	args := fields[1:]

	return command, args
}

func (s *Shell) echoCommand(args []string) {
	fmt.Println(strings.Join(args, " "))
}

func (s *Shell) cdCommand(args []string) {
	err := os.Chdir(args[0])
	if err != nil {
		fmt.Println("cd:", args[0]+": No such file or directory")
		return
	}
}

func (s *Shell) pwdCommand(_ []string) {
	path, err := os.Getwd()
	if err != nil {
		fmt.Println("pwd: failed to get current directory")
		return
	}

	fmt.Println(path)
}

func (s *Shell) typeCommand(args []string) {
	if len(args) == 0 {
		return
	}

	command := args[0]

	if command == "exit" {
		fmt.Println(command + " is a shell builtin")
		return
	}

	if _, ok := s.builtins[command]; ok {
		fmt.Println(command + " is a shell builtin")
		return
	}

	fullPath, found := findExecutable(command)
	if found {
		fmt.Println(command, "is", fullPath)
		return
	}

	fmt.Println(command + ": not found")
}

func findExecutable(command string) (string, bool) {
	path := os.Getenv("PATH")
	dirs := filepath.SplitList(path)

	for _, dir := range dirs {
		fullPath := filepath.Join(dir, command)

		info, err := os.Stat(fullPath)
		if err != nil {
			continue
		}

		if info.IsDir() {
			continue
		}

		if info.Mode().Perm()&0o111 != 0 {
			return fullPath, true
		}
	}

	return "", false
}

func (s *Shell) runExternal(command string, args []string) {
	_, found := findExecutable(command)
	if !found {
		fmt.Println(command + ": command not found")
		return
	}

	cmd := exec.Command(command, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	_ = cmd.Run()
}

func main() {
	shell := NewShell()
	shell.Run()
}
