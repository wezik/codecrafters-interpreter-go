package main

import (
	"fmt"
	"os"
	"slices"

	"github.com/codecrafters-io/interpreter-starter-go/cmd/myinterpreter/cmd_parse"
	"github.com/codecrafters-io/interpreter-starter-go/cmd/myinterpreter/cmd_tokenize"
)

func main() {
	fileContents, cmd := readArgs(os.Args)

	if cmd == "tokenize" {
		cmd_tokenize.Run(fileContents)
	} else if cmd == "parse" {
		cmd_parse.Run(fileContents)
	}

	os.Exit(0)
}

var availableCommands = []string{
	"tokenize",
	"parse",
}

func readArgs(args []string) ([]byte, string) {
	if len(args) < 3 {
		fmt.Fprintln(os.Stderr, "Usage: ./your_program.sh tokenize/parse <filename>")
		os.Exit(1)
	}

	command := args[1]

	if !slices.Contains(availableCommands, command) {
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
		os.Exit(1)
	}

	filename := os.Args[2]
	fileContents, err := os.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}

	return fileContents, command
}
