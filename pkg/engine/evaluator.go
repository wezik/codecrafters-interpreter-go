package engine

import "fmt"

func Evaluate(ast []Expr) (string, []error) {
	stringAST := ""
	var errs []error
	return stringAST, append(errs, fmt.Errorf("Error: unimplemented"))
}
