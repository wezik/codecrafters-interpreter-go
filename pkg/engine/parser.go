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

type ExprUnary struct {
	LexToken LexToken
	Expr     Expr
}

func (e ExprUnary) String() string {
	return fmt.Sprintf("(%s %s)", e.LexToken.Lexeme, e.Expr.String())
}

type ExprGroup struct {
	Exprs  []Expr
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

type ParserState struct {
	errs []error
	globalGroup *ExprGroup
	currentGroup *ExprGroup
	readIndex int
	tokens []LexToken
}

// func Parse
//
// func Parse(lexTokens []LexToken) ([]Expr, []error) {
// 	var errs []error
// 	globalGroup := &ExprGroup{Exprs: []Expr{}, Parent: nil}
// 	currentGroup := globalGroup
//
// 	for lexI := 0; lexI < len(lexTokens); lexI++ {
// 		lexT := lexTokens[lexI]
// 		switch lexT.TokenType {
// 		case TOKEN_LEFT_PAREN:
// 			newGroup := ExprGroup{Exprs: []Expr{}, Parent: currentGroup}
// 			currentGroup.Exprs = append(currentGroup.Exprs, &newGroup)
// 			currentGroup = &newGroup
// 		case TOKEN_RIGHT_PAREN:
// 			if currentGroup.Parent == nil {
// 				errs = append(errs, fmt.Errorf("Error: Unmatched parentheses."))
// 				continue
// 			} else if len(currentGroup.Exprs) == 0 {
// 				currentGroup.Parent.Exprs = []Expr{}
// 				errs = append(errs, fmt.Errorf("Error: Empty group."))
// 			}
// 			currentGroup = currentGroup.Parent
// 		case TOKEN_NUMBER, TOKEN_STRING, TOKEN_TRUE, TOKEN_FALSE, TOKEN_NIL:
// 			currentGroup.Exprs = append(currentGroup.Exprs, ExprLiteral{lexT})
// 		case TOKEN_MINUS, TOKEN_BANG:
// 			if lexI < len(lexTokens)-1 {
// 				nextLexT := lexTokens[lexI+1]
// 				currentGroup.Exprs = append(currentGroup.Exprs, ExprUnary{lexT, ExprLiteral{nextLexT}})
// 				lexI++
// 				continue
// 			}
// 			errs = append(errs, fmt.Errorf("Error: Unexpected token: %s", lexT.Lexeme))
// 		}
// 	}
// 	if currentGroup.Parent != nil {
// 		errs = append(errs, fmt.Errorf("Error: Unmatched parentheses."))
// 		currentGroup.Parent.Exprs = []Expr{}
// 	}
// 	return globalGroup.Exprs, errs
// }
//
// var tokenExpressions = map[string]func(*ParserState) Expr{}
//
// func init() {
// 	tokenExpressions = map[string]func(*ParserState) Expr{
// 		TOKEN_NUMBER: parseLiteral,
// 		TOKEN_STRING: parseLiteral,
// 		TOKEN_TRUE: parseLiteral,
// 		TOKEN_FALSE: parseLiteral,
// 		TOKEN_NIL: parseLiteral,
// 		TOKEN_MINUS: parseUnary,
// 		TOKEN_BANG: parseUnary,
// 		TOKEN_LEFT_PAREN: parseNewGroup,
// 	}
// }

// func Parse(lexTokens []LexToken) ([]Expr, []error) {
// 	state := ParserState{
// 		errs: []error{},
// 		globalGroup: &ExprGroup{Exprs: []Expr{}, Parent: nil},
// 		readIndex: 0,
// 		tokens: lexTokens,
// 	}
// 	state.currentGroup = state.globalGroup
//
// 	for ; state.readIndex < len(state.tokens); state.readIndex++ {
// 		lexT := state.tokens[state.readIndex]
// 		exprs := &state.currentGroup.Exprs
// 		expr, err := processToken(&state, lexT)
// 		if err != nil {
// 			state.errs = append(state.errs, err)
// 			continue
// 		}
// 		if expr == nil {
// 			continue
// 		}
// 		*exprs = append(*exprs, expr)
// 	}
//
// 	return state.currentGroup.Exprs, state.errs
// }
//
// func processToken(state *ParserState, lexT LexToken) (Expr, error) {
// 	exprs := &state.currentGroup.Exprs
// 	switch lexT.TokenType {
// 	case TOKEN_EOF:
// 		break
// 	case TOKEN_NUMBER, TOKEN_STRING, TOKEN_TRUE, TOKEN_FALSE, TOKEN_NIL:
// 		*exprs = append(*exprs, ExprLiteral{lexT})
// 	case TOKEN_MINUS, TOKEN_BANG:
// 		parseUnary(state)
// 	case TOKEN_LEFT_PAREN:
// 		parseGroup(state), nil
// 	case TOKEN_RIGHT_PAREN:
// 		return 
// 	}
// 	return nil, fmt.Errorf("Error: Unexpected token: %s", lexT.Lexeme)
// }
//
// func parseGroup(state *ParserState) Expr {
// 	group := ExprGroup{Exprs: []Expr{}, Parent: state.currentGroup}
// 	state.currentGroup.Exprs = append(state.currentGroup.Exprs, &group)
// 	state.currentGroup = &group
// 	return &group
// }
//
// func parseGroup(state *ParserState) (*ExprGroup, error) {
// 	newGroup := ExprGroup{Exprs: []Expr{}, Parent: state.currentGroup}
// 	state.currentGroup = &newGroup
// 	for ; state.readIndex < len(state.tokens); state.readIndex++ {
// 		lexT := state.tokens[state.readIndex]
// 		expr, err := processToken(state, lexT)
// 		if err != nil {
// 			state.errs = append(state.errs, err)
// 			continue
// 		}
// 		state.currentGroup.Exprs = append(state.currentGroup.Exprs, expr)
// 	}
// 	state.currentGroup = state.currentGroup.Parent
// 	return &newGroup, nil
// }
//
//
// func parseUnary(state *ParserState) error {
// 	exprUnary := ExprUnary{state.tokens[state.readIndex], nil}
// 	state.readIndex++
// 	nextLexT := state.tokens[state.readIndex]
// 	expr, err := processToken(state, nextLexT)
// 	if err != nil {
// 		return nil, err
// 	}
// 	exprUnary.Expr = expr
//
// }

func Parse(lexTokens []LexToken) ([]Expr, []error) {
	var topExpr Expr
	var errs []error



	return []Expr{topExpr}, errs
}

