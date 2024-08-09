package engine

import (
	"strconv"
)

func Evaluate(ast []Expr) (string, []error) {
	stringAST := ""
	var errs []error
	for _, expr := range ast {
		if expr, ok := expr.(*ExprLiteral); ok {
			if expr.LexToken.TokenType == TOKEN_NUMBER {
				stringAST += trimEmptyDecimal(expr.String())
			} else {
				stringAST += expr.String()
			}
		}
	}
	return stringAST, errs
}

func trimEmptyDecimal(s string) string {
	f, _ := strconv.ParseFloat(s, 64)
	return strconv.FormatFloat(f, 'f', -1, 64)
}
