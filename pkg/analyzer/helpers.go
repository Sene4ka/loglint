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
VarPart is a type for elements of logger call expression first argument
ex. in log.Print("password: " + password) there is 2 VarPart objects
"password: " of PartConst and password of PartVar types
*/
type VarPart struct {
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

// recursively parse AST node to format it into VarPart array in appearance order
func getStringLiterals(expr ast.Expr) []VarPart {
	switch e := expr.(type) {
	case *ast.BasicLit:
		if e.Kind == token.STRING {
			val, _ := strconv.Unquote(e.Value)
			return []VarPart{{
				Value: val,
				Pos:   e.Pos() + 1,
				End:   e.End() - 1,
				Type:  PartConst,
			}}
		}
	case *ast.BinaryExpr:
		if e.Op == token.ADD {
			left := getStringLiterals(e.X)
			right := getStringLiterals(e.Y)
			return mergeArrays(left, right)
		}
	case *ast.Ident:
		return []VarPart{{
			Value: e.Name,
			Pos:   e.Pos(),
			End:   e.End(),
			Type:  PartVar,
		}}
	}

	// return PartOther to naturally separate sequence of PartConst strings preventing their merging later
	return []VarPart{{
		Value: "",
		Pos:   token.NoPos,
		End:   token.NoPos,
		Type:  PartOther,
	}}
}

func mergeArrays[T any](a, b []T) []T {
	result := make([]T, len(a)+len(b))
	copy(result, a)
	copy(result[len(a):], b)
	return result
}

/*
merge continuous sequences of PartConst strings so we can detect harder sensitive keyword patterns later
like log.Print("a" + "u" + "t" + "h" + "_token: " + auth_token) will in the end merge into
{PartConst{Value: "auth_token"}, PartVar{Value: "auth_token"}}
*/
func foldConstantStrings(parts []VarPart) []VarPart {
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

	result := make([]VarPart, 0, estimatedLen)
	builder := strings.Builder{}
	builder.Grow(mxLen)
	processed := 0
	for i, part := range parts {
		if part.Type == PartConst {
			processed++
			builder.WriteString(part.Value)
		} else {
			if processed > 0 {
				result = append(result, VarPart{
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
		result = append(result, VarPart{
			Value: builder.String(),
			Pos:   parts[len(parts)-processed].Pos,
			End:   parts[len(parts)-1].End,
			Type:  PartConst,
		})
	}
	return result
}
