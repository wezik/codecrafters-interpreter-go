package main

import (
	"fmt"
	"os"
)

type LexToken struct {
	tokenType string
	lexeme    string
	literal   string
}

func newToken(tokenType string, lexeme string, literal string) LexToken {
	return LexToken{tokenType, lexeme, literal}
}

func newTokenNoLit(tokenType string, lexeme string) LexToken {
	return LexToken{tokenType, lexeme, "null"}
}

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Fprintln(os.Stderr, "Logs from your program will appear here!")

	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "Usage: ./your_program.sh tokenize <filename>")
		os.Exit(1)
	}

	command := os.Args[1]

	if command != "tokenize" {
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
		os.Exit(1)
	}

	filename := os.Args[2]
	fileContents, err := os.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}

	tokens := []LexToken{}

	if len(fileContents) > 0 {
		for _, b := range fileContents {
			if b == '{' {
				tokens = append(tokens, newTokenNoLit("LEFT_BRACE", "{"))
			} else if b == '}' {
				tokens = append(tokens, newTokenNoLit("RIGHT_BRACE", "}"))
			} else if b == '(' {
				tokens = append(tokens, newTokenNoLit("LEFT_PAREN", "("))
			} else if b == ')' {
				tokens = append(tokens, newTokenNoLit("RIGHT_PAREN", ")"))
			}
		}
		tokens = append(tokens, newTokenNoLit("EOF", ""))
	} else {
		tokens = append(tokens, newTokenNoLit("EOF", ""))
	}

	for _, t := range tokens {
		fmt.Printf("%s %s %s\n", t.tokenType, t.lexeme, t.literal)
	}
}
