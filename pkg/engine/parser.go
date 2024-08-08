package engine

import (
	"fmt"
	"strings"
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
	case TOKEN_TRUE, TOKEN_FALSE:
		return e.LexToken.Lexeme
	case TOKEN_NIL:
		return e.LexToken.Lexeme
	}
	return "nil"
}

type ExprGroup struct {
	Exprs []Expr
	Parent *ExprGroup
}

func (e *ExprGroup) String() string {
	var exprs []string
	for _, expr := range e.Exprs {
		exprs = append(exprs, expr.String())
	}
	if e.Parent == nil {
		return strings.Join(exprs, " ")
	}
	return fmt.Sprintf("(group %s)", strings.Join(exprs, " "))
}

func Parse(lexTokens []LexToken) ([]Expr, []error) {
	var errs []error
	globalGroup := &ExprGroup{Exprs: []Expr{}, Parent: nil}
	currentGroup := globalGroup

	for _, lexT := range lexTokens {
		switch lexT.TokenType {
		case TOKEN_LEFT_PAREN:
			newGroup := ExprGroup{Exprs: []Expr{}, Parent: currentGroup}
			currentGroup.Exprs = append(currentGroup.Exprs, &newGroup)
			currentGroup = &newGroup
		case TOKEN_RIGHT_PAREN:
			if currentGroup.Parent == nil {
				errs = append(errs, fmt.Errorf("Error: Unmatched parentheses."))
				continue
			} else if len(currentGroup.Exprs) == 0 {
				currentGroup.Parent.Exprs = []Expr{}
				errs = append(errs, fmt.Errorf("Error: Empty group."))
			}
			currentGroup = currentGroup.Parent
		case TOKEN_NUMBER, TOKEN_STRING, TOKEN_TRUE, TOKEN_FALSE, TOKEN_NIL:
			currentGroup.Exprs = append(currentGroup.Exprs, ExprLiteral{lexT})
		}
	}
	if currentGroup.Parent != nil {
		errs = append(errs, fmt.Errorf("Error: Unmatched parentheses."))
		currentGroup.Parent.Exprs = []Expr{}
	}
	return globalGroup.Exprs, errs
}
