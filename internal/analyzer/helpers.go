package analyzer

import (
	"go/ast"
	"go/token"
	"strconv"
)

var loggers = map[string][]string{
	"log":  {"Print", "Printf", "Println", "Fatal", "Fatalf", "Panic", "Panicf"},
	"slog": {"Debug", "Info", "Warn", "Error", "Log"},
	"zap":  {"Debug", "Info", "Warn", "Error", "DPanic", "Panic", "Fatal"},
}

func isSupportedLoggerCall(call *ast.CallExpr) bool {
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return false
	}

	pkgName, ok := sel.X.(*ast.Ident)
	if !ok {
		return false
	}

	methodName := sel.Sel.Name

	methods, exists := loggers[pkgName.Name]
	if !exists {
		return false
	}

	for _, m := range methods {
		if methodName == m {
			return true
		}
	}
	return false
}

func getStringLiteral(expr ast.Expr) (string, token.Pos, bool) {
	if lit, ok := expr.(*ast.BasicLit); ok && lit.Kind == token.STRING {
		value, err := strconv.Unquote(lit.Value)
		if err == nil {
			return value, lit.Pos() + 1, true
		}
	}
	return "", token.NoPos, false
}
