package util

import (
	"fmt"
	"os"

	"github.com/codecrafters-io/interpreter-starter-go/pkg/engine"
)

func PrintTokens(tokens []engine.LexToken) {
	for _, t := range tokens {
		fmt.Printf("%s %s %s\n", t.TokenType, t.Lexeme, t.Literal)
	}
}

func PrintAST(ast []engine.Expr) {
	for _, a := range ast {
		fmt.Println(a)
	}
}

func PrintErrors(errors []error) {
	for _, e := range errors {
		fmt.Fprintln(os.Stderr, e.Error())
	}
}
