package engine

import (
	"fmt"
)

type Expr interface {
	String() string
}

type ExprLiteral struct {
	LexToken LexToken
}

func (e ExprLiteral) String() string {
	tokenType := e.LexToken.TokenType
	switch tokenType {
	case TOKEN_NUMBER, TOKEN_STRING:
		return e.LexToken.Literal
	default:
		return e.LexToken.Lexeme
	}
}

type ExprGroup struct {
	expr Expr
}

func (e *ExprGroup) String() string {
	return fmt.Sprintf("(group %s)", e.expr)
}

type ExprUnary struct {
	LexToken LexToken
	Expr     *Expr
}

func (e *ExprUnary) String() string {
	return fmt.Sprintf("(%s %s)", e.LexToken.Lexeme, (*e.Expr).String())
}

type ExprBinary struct {
	Left     *Expr
	Operator LexToken
	Right    *Expr
}

type Parser struct {
	tokens []LexToken
	tokensIndex int
}

func (p *Parser) current() LexToken {
	return p.tokens[p.tokensIndex]
}

func (p *Parser) previous() LexToken {
	return p.tokens[p.tokensIndex-1]
}

func (p *Parser) advance() LexToken {
	if !p.isAtEnd() {
		p.tokensIndex++
	}
	return p.previous()
}

func (p *Parser) check(tokenType string) bool {
	return p.current().TokenType == tokenType
}

func (p *Parser) match(tokenTypes ...string) bool {
	for _, tokenType := range tokenTypes {
		if p.check(tokenType) {
			return true
		}
	}
	return false
}

func (p *Parser) isAtEnd() bool {
	return p.peek() == TOKEN_EOF
}

func (p *Parser) peek() string {
	return p.tokens[p.tokensIndex].TokenType
}

func Parse(lexTokens []LexToken) ([]Expr, []error) {
	parser := Parser{tokens: lexTokens, tokensIndex: 0}
	var exprs []Expr
	var errs []error

	for !parser.isAtEnd() {
		expr, err := parseExpression(&parser)
		if err != nil {
			errs = append(errs, err)
		}
		if expr != nil {
			exprs = append(exprs, expr)
		}
	}
	return exprs, errs
}

func parseExpression(parser *Parser) (Expr, error) {
	return parseEquality(parser)
}

func parseEquality(parser *Parser) (Expr, error) {
	if false {
	}

	return parseComparison(parser)
}

func parseComparison(parser *Parser) (Expr, error) {
	if false {
	}

	return parseTerm(parser)
}

func parseTerm(parser *Parser) (Expr, error) {
	if false {
	}

	return parseFactor(parser)
}

func parseFactor(parser *Parser) (Expr, error) {
	if false {
	}

	return parseUnary(parser)
}

func parseUnary(parser *Parser) (Expr, error) {
	if parser.match(TOKEN_BANG, TOKEN_MINUS) {
		parser.advance()
		operator := parser.previous()
		right, err := parseUnary(parser)
		if err != nil {
			return nil, err
		}
		return &ExprUnary{operator, &right}, nil
	}
	return parsePrimary(parser)
}

func parsePrimary(parser *Parser) (Expr, error) {
	if parser.match(TOKEN_NUMBER, TOKEN_STRING, TOKEN_TRUE, TOKEN_FALSE, TOKEN_NIL) {
		parser.advance()
		return &ExprLiteral{parser.previous()}, nil
	} else if parser.check(TOKEN_LEFT_PAREN) {
		parser.advance()
		expr, err := parseExpression(parser)
		if err != nil {
			parser.advance()
			return nil, err
		}
		if expr == nil {
			parser.advance()
			return nil, nil
		}
		if parser.check(TOKEN_RIGHT_PAREN) {
			parser.advance()
			return &ExprGroup{expr}, nil
		}
	}
	return nil, fmt.Errorf("Error: unmatched parenthesis")
}
