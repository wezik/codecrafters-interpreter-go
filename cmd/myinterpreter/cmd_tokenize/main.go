package cmd_tokenize

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
		exitCode = 65
	}

	util.PrintTokens(lexTokens)
	os.Exit(exitCode)
}
