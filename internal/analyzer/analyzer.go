package analyzer

import (
	"go/ast"
	"go/token"
	"strconv"
	"unicode"
	"unicode/utf8"

	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	Name: "log-lint",
	Doc:  "reports wrong log function usage",
	Run:  run,
}

var loggers = map[string][]string{
	"log":  {"Print", "Printf", "Println", "Fatal", "Fatalf", "Panic", "Panicf"},
	"slog": {"Debug", "Info", "Warn", "Error", "Log"},
	"zap":  {"Debug", "Info", "Warn", "Error", "DPanic", "Panic", "Fatal"},
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

			return true
		})
	}
	return nil, nil
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

func checkLowercaseStart(msg string, pos token.Pos, pass *analysis.Pass) {
	if len(msg) == 0 {
		return
	}

	r, size := utf8.DecodeRuneInString(msg)
	if r == utf8.RuneError && size == 1 {
		return
	}

	if unicode.IsUpper(r) {
		diag := analysis.Diagnostic{
			Pos:      pos,
			Category: "Warning",
			Message:  "log message should start with lowercase letter",
		}

		fixed := string(unicode.ToLower(r)) + msg[size:]
		diag.SuggestedFixes = []analysis.SuggestedFix{
			{
				Message: "Convert first letter to lowercase",
				TextEdits: []analysis.TextEdit{
					{
						Pos:     pos,
						End:     pos + token.Pos(size),
						NewText: []byte(fixed),
					},
				},
			},
		}

		pass.Report(diag)
	}
}
