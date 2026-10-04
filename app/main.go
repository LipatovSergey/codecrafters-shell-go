package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"unicode"
)

type builtinFunc func([]string, io.Writer)

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

		command, args, outputFile := parseInput(input)

		if command == "" {
			continue
		}

		if command == "exit" {
			return
		}

		if handler, ok := s.builtins[command]; ok {
			if outputFile != "" {
				file, err := os.Create(outputFile)
				if err != nil {
					fmt.Fprintln(os.Stderr, err)
				}
				handler(args, file)
				file.Close()
			} else {
				handler(args, os.Stdout)
			}
			continue
		}

		s.runExternal(command, args)
	}
}

func parseInput(input string) (string, []string, string) {
	result := []string{}
	currentArg := ""
	argStarted := false
	inSingleQuotes := false
	inDoubleQuotes := false
	nextToBackSlash := false
	outupRedirection := false
	outputFile := ""
	for i, r := range input {
		switch {
		case r == '\'' && !inDoubleQuotes && !nextToBackSlash:
			inSingleQuotes = !inSingleQuotes

		case r == '"' && !inSingleQuotes && !nextToBackSlash:
			inDoubleQuotes = !inDoubleQuotes

		case r == '>' || (r == '1' && input[i+1] == '>') && !inSingleQuotes && !inDoubleQuotes && !nextToBackSlash:
			outupRedirection = true

		case !inSingleQuotes && !inDoubleQuotes && !nextToBackSlash && unicode.IsSpace(r):
			if argStarted {
				if outupRedirection {
					outputFile = currentArg
					fmt.Println("output file:", outputFile)
					outupRedirection = false
				} else {
					result = append(result, currentArg)
				}
				currentArg = ""
				argStarted = false
			}

		case r == '\\' && !inSingleQuotes && !nextToBackSlash:
			nextToBackSlash = true

		case nextToBackSlash:
			currentArg += string(r)
			nextToBackSlash = false

		default:
			argStarted = true
			currentArg += string(r)
		}
	}
	return result[0], result[1:], outputFile
}

func (s *Shell) echoCommand(args []string, output io.Writer) {
	fmt.Fprintln(output, strings.Join(args, " "))
}

func (s *Shell) cdCommand(args []string, _ io.Writer) {
	path := args[0]
	if path == "~" {
		home := os.Getenv("HOME")
		os.Chdir(home)
		return
	}

	err := os.Chdir(path)
	if err != nil {
		fmt.Fprintln(os.Stderr, "cd:", args[0]+": No such file or directory")
		return
	}
}

func (s *Shell) pwdCommand(_ []string, output io.Writer) {
	path, err := os.Getwd()
	if err != nil {
		fmt.Fprintln(os.Stderr, "pwd: failed to get current directory")
		return
	}

	fmt.Fprintln(output, path)
}

func (s *Shell) typeCommand(args []string, output io.Writer) {
	if len(args) == 0 {
		return
	}

	command := args[0]

	if command == "exit" {
		fmt.Fprintln(output, command+" is a shell builtin")
		return
	}

	if _, ok := s.builtins[command]; ok {
		fmt.Fprintln(output, command+" is a shell builtin")
		return
	}

	fullPath, found := findExecutable(command)
	if found {
		fmt.Fprintln(output, command, "is", fullPath)
		return
	}

	fmt.Fprintln(os.Stderr, command+": not found")
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
