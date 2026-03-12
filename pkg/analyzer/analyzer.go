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

			parts := make([]ExprPart, 0, 2)

			parts = getExpressionParts(call.Args[0], parts)

			parts = foldConstantStrings(parts)

			if len(call.Args) > 1 {
				for _, arg := range call.Args[1:] {
					parts = getExpressionParts(arg, parts)
				}
			}

			checkShouldStartWithLowercase(parts, pass)
			checkShouldContainOnlyEnglish(parts, pass)
			checkShouldNotContainSpecialSymbols(parts, pass)
			checkShouldNotContainSensitiveInformation(parts, pass)

			return true
		})
	}
	return nil, nil
}
