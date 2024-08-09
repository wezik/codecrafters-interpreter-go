package engine

func Evaluate(ast []Expr) (string, []error) {
	stringAST := ""
	var errs []error
	for _, expr := range ast {
		if expr, ok := expr.(*ExprLiteral); ok {
			stringAST += expr.String()
		}
	}
	return stringAST, errs
}
