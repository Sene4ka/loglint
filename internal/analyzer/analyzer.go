package analyzer

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	Name: "log-lint",
	Doc:  "reports wrong log function usage",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		ast.Inspect(file, func(n ast.Node) bool {
			call, ok := n.(*ast.CallExpr)
			if !ok {
				return true
			}

			if !isSupportedLoggerCall(call) {
				return true
			}

			if len(call.Args) < 1 {
				return true
			}

			msg, pos, ok := getStringLiteral(call.Args[0])
			if !ok {
				return true
			}

			checkLowercaseStart(msg, pos, pass)
			checkOnlyLatinLetters(msg, pos, pass)
			checkNoSpecialSymbols(msg, pos, pass)
			checkNoKeywords(msg, pos, pass)

			return true
		})
	}
	return nil, nil
}
