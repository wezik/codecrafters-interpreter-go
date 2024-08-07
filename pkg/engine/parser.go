package engine

import "slices"

type LexAST struct {
	//TODO
	value string
}

func (a LexAST) String() string {
	//TODO
	return a.value
}

var tokenToLexeme = []string{
	TOKEN_TRUE,
	TOKEN_FALSE,
	TOKEN_NIL,
}

var tokenToLiteral = []string{
	TOKEN_NUMBER,
	TOKEN_STRING,
}

func Parse(lexTokens []LexToken) ([]LexAST, []error) {
	var ast []LexAST
	for _, t := range lexTokens {
		if slices.Contains(tokenToLexeme, t.TokenType) {
			ast = append(ast, LexAST{t.Lexeme})
		} else if slices.Contains(tokenToLiteral, t.TokenType) {
			ast = append(ast, LexAST{t.Literal})
		}
	}
	return ast, nil
}
