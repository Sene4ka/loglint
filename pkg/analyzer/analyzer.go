package analyzer

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	Name: "loglint",
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

			if !isSupportedLoggerCall(call, pass) {
				return true
			}

			if len(call.Args) < 1 {
				return true
			}

			parts := getStringLiterals(call.Args[0])

			parts = foldConstantStrings(parts)

			checkLowercaseStart(parts, pass)
			checkOnlyLatinLetters(parts, pass)
			checkNoSpecialSymbols(parts, pass)
			checkNoSensitiveKeywordsAndVariables(parts, pass)

			return true
		})
	}
	return nil, nil
}
