package main

import (
	"fmt"
	"os"
	"slices"

	"github.com/codecrafters-io/interpreter-starter-go/cmd/myinterpreter/cmd_evaluate"
	"github.com/codecrafters-io/interpreter-starter-go/cmd/myinterpreter/cmd_parse"
	"github.com/codecrafters-io/interpreter-starter-go/cmd/myinterpreter/cmd_tokenize"
)

func main() {
	fileContents, cmd := readArgs(os.Args)

	if cmd == CMD_TOKENIZE || cmd == CMD_T {
		cmd_tokenize.Run(fileContents)
	} else if cmd == CMD_PARSE || cmd == CMD_P {
		cmd_parse.Run(fileContents)
	} else if cmd == CMD_EVALUATE || cmd == CMD_EVAL || cmd == CMD_E {
		cmd_evaluate.Run(fileContents)
	}

	os.Exit(0)
}

const (
	CMD_TOKENIZE = "tokenize"
	CMD_T        = "t"
	CMD_PARSE    = "parse"
	CMD_P        = "p"
	CMD_EVALUATE = "evaluate"
	CMD_EVAL     = "eval"
	CMD_E        = "e"
)

var availableCommands = []string{
	CMD_TOKENIZE,
	CMD_T,
	CMD_PARSE,
	CMD_P,
	CMD_EVALUATE,
	CMD_EVAL,
	CMD_E,
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
