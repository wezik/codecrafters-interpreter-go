package cmd_evaluate

import (
	"fmt"
	"os"

	"github.com/codecrafters-io/interpreter-starter-go/pkg/engine"
	util "github.com/codecrafters-io/interpreter-starter-go/pkg/util/cmd_util"
)

func Run(fileContents []byte) {
	lexTokens, errs := engine.Tokenize(fileContents)

	if len(errs) > 0 {
		util.PrintErrors(errs)
		os.Exit(65)
	}

	ast, errs := engine.Parse(lexTokens)

	if len(errs) > 0 {
		util.PrintErrors(errs)
		os.Exit(65)
	}

	evaluatedAST, errs := engine.Evaluate(ast)

	if len(errs) > 0 {
		util.PrintErrors(errs)
		os.Exit(65)
	}

	fmt.Println(evaluatedAST)

	os.Exit(0)
}
