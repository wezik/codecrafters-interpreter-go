package engine

import (
	"slices"
)

type Expr interface {
	String() string
}

type ExprLiteral struct {
	lexToken LexToken
}

func (e ExprLiteral) String() string {
	tokenType := e.lexToken.TokenType
	switch tokenType {
	case TOKEN_NUMBER, TOKEN_STRING:
		return e.lexToken.Literal
	case TOKEN_TRUE, TOKEN_FALSE:
		return e.lexToken.Lexeme
	case TOKEN_NIL:
		return e.lexToken.Lexeme
	}
	return "nil"
}

var literalTokens = []string{
	TOKEN_NUMBER,
	TOKEN_STRING,
	TOKEN_TRUE,
	TOKEN_FALSE,
	TOKEN_NIL,
}

func Parse(lexTokens []LexToken) ([]Expr, []error) {
	var ast []Expr
	for i = 0; i < len(lexTokens); i++ {
		t := lexTokens[i]
		if slices.Contains(literalTokens, t.TokenType) {
			ast = append(ast, ExprLiteral{t})
		}
	}
	return ast, nil
}
