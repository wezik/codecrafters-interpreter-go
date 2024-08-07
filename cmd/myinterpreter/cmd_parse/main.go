package cmd_parse

import (
	"os"

	"github.com/codecrafters-io/interpreter-starter-go/pkg/engine"
	util "github.com/codecrafters-io/interpreter-starter-go/pkg/util/cmd_util"
)

func Run(fileContents []byte) {
	exitCode := 0
	lexTokens, errs := engine.Tokenize(fileContents)

	if len(errs) > 0 {
		util.PrintErrors(errs)
		os.Exit(65)
	}

	ast, errs := engine.Parse(lexTokens)

	if len(errs) > 0 {
		util.PrintErrors(errs)
		exitCode = 65
	}

	util.PrintAST(ast)
	os.Exit(exitCode)
}
