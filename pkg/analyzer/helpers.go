package analyzer

import (
	"go/ast"
	"go/token"
	"strconv"
	"strings"

	"golang.org/x/tools/go/analysis"
)

var loggers = map[string][]string{
	"log":  {"Print", "Printf", "Println", "Fatal", "Fatalf", "Panic", "Panicf"},
	"slog": {"Debug", "Info", "Warn", "Error", "Log"},
	"zap":  {"Debug", "Info", "Warn", "Error", "DPanic", "Panic", "Fatal"},
}

/*
ExprPart is a type for elements of logger call expression first argument
ex. in log.Print("password: " + password) there is 2 ExprPart objects
"password: " of PartConst and password of PartVar types
*/
type ExprPart struct {
	Value string
	Pos   token.Pos
	End   token.Pos
	Type  PartType
}

/*
PartType is a type of path, can be PartConst for constant strings, PartVar for variables and PartOther for
other things like other function calls or non-add expressions
*/
type PartType int

const (
	PartConst PartType = iota
	PartVar
	PartOther
)

func isSupportedLoggerCall(call *ast.CallExpr, pass *analysis.Pass) bool {
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return false
	}

	ident, ok := sel.X.(*ast.Ident)
	if !ok {
		return false
	}
	selName := ident.Name

	// zap logger is of zap.Logger type instead of package call
	// so we need to type check because selector name might be anything like "logger"
	if pass.TypesInfo != nil {
		tp := pass.TypesInfo.TypeOf(sel.X)
		if tp != nil {
			// simple type check by name because we need mock package in testdata
			// do actual type check instead if needed, or full string equality type check like "go.uber.org/zap.Logger"
			// but it will crash analyzer tests because they can't use external dependencies
			typeStr := tp.String()
			if strings.Contains(typeStr, "zap.Logger") {
				selName = "zap"
			} else if strings.Contains(typeStr, "slog.Logger") {
				selName = "slog"
			}
		}
	}

	methodName := sel.Sel.Name

	methods, exists := loggers[selName]
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

// recursively parse AST node to format it into ExprPart array in appearance order and append them in output array
func getExpressionParts(expr ast.Expr, output []ExprPart) []ExprPart {
	switch e := expr.(type) {
	case *ast.BasicLit:
		if e.Kind == token.STRING {
			val, _ := strconv.Unquote(e.Value)
			return append(output, ExprPart{
				Value: val,
				Pos:   e.Pos() + 1,
				End:   e.End() - 1,
				Type:  PartConst,
			})
		}
	case *ast.BinaryExpr:
		if e.Op == token.ADD {
			output = getExpressionParts(e.X, output)
			return getExpressionParts(e.Y, output)
		}
	case *ast.Ident:
		return append(output, ExprPart{
			Value: e.Name,
			Pos:   e.Pos(),
			End:   e.End(),
			Type:  PartVar,
		})
	case *ast.CallExpr:
		/*
			for now ignore function selector and method, works well with fmt like Sprint/Sprintf,
			but can false-positive on other functions which (depending on their implementation)
		*/
		for _, arg := range e.Args {
			output = getExpressionParts(arg, output)
		}
		return output
	}

	// return PartOther to naturally separate sequence of PartConst strings preventing their merging later
	return append(output, ExprPart{
		Value: "",
		Pos:   token.NoPos,
		End:   token.NoPos,
		Type:  PartOther,
	})
}

/*
merge continuous sequences of PartConst strings so we can detect harder sensitive keyword patterns later
like log.Print("a" + "u" + "t" + "h" + "_token: " + auth_token) will in the end merge into
{PartConst{Value: "auth_token"}, PartVar{Value: "auth_token"}}
*/
func foldConstantStrings(parts []ExprPart) []ExprPart {
	estimatedLen := 0
	constLen := 0
	valueLen := 0
	mxLen := 0
	for _, part := range parts {
		if part.Type == PartConst && constLen == 0 {
			constLen++
			valueLen += len(part.Value)
		} else if part.Type != PartConst {
			estimatedLen = constLen + 1
			mxLen = max(mxLen, valueLen)
			constLen = 0
			valueLen = 0
		}
	}

	result := make([]ExprPart, 0, estimatedLen)
	builder := strings.Builder{}
	builder.Grow(mxLen)
	processed := 0
	for i, part := range parts {
		if part.Type == PartConst {
			processed++
			builder.WriteString(part.Value)
		} else {
			if processed > 0 {
				result = append(result, ExprPart{
					Value: builder.String(),
					Pos:   parts[i-processed].Pos,
					End:   parts[i-1].End,
					Type:  PartConst,
				})
				processed = 0
				builder.Reset()
			}
			result = append(result, part)
		}
	}
	if processed > 0 {
		result = append(result, ExprPart{
			Value: builder.String(),
			Pos:   parts[len(parts)-processed].Pos,
			End:   parts[len(parts)-1].End,
			Type:  PartConst,
		})
	}
	return result
}
