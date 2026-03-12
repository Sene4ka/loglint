package analyzer

import (
	"flag"
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	Name:  "log-linter",
	Doc:   "Static analysis tool for detecting style violations and sensitive data leaks in Go log messages.",
	Run:   run,
	Flags: flag.FlagSet{},
}

var (
	allowedSpecialSymbolsFlag string
	sensitiveKeywordsFlag     string
	ruleLowercaseFlag         bool
	ruleEnglishFlag           bool
	ruleSymbolsFlag           bool
	ruleSensitiveFlag         bool
)

func init() {
	Analyzer.Flags.StringVar(&allowedSpecialSymbolsFlag, "allowed-special-symbols",
		": _ = %", "Allowed special symbols (space-separated)")
	Analyzer.Flags.StringVar(&sensitiveKeywordsFlag, "sensitive-keywords",
		"key,password,secret,auth,token", "Sensitive keywords (comma-separated)")
	Analyzer.Flags.BoolVar(&ruleLowercaseFlag, "rule-lowercase", true,
		"Enable shouldStartWithLowercase rule")
	Analyzer.Flags.BoolVar(&ruleEnglishFlag, "rule-english", true,
		"Enable shouldContainOnlyEnglish rule")
	Analyzer.Flags.BoolVar(&ruleSymbolsFlag, "rule-symbols", true,
		"Enable shouldNotContainSpecialSymbols rule")
	Analyzer.Flags.BoolVar(&ruleSensitiveFlag, "rule-sensitive", true,
		"Enable shouldNotContainSensitiveInformation rule")
}

func run(pass *analysis.Pass) (interface{}, error) {
	if !configInitialized {
		rules := map[string]bool{
			"shouldStartWithLowercase":             ruleLowercaseFlag,
			"shouldContainOnlyEnglish":             ruleEnglishFlag,
			"shouldNotContainSpecialSymbols":       ruleSymbolsFlag,
			"shouldNotContainSensitiveInformation": ruleSensitiveFlag,
		}

		cfg := LoadConfigFromFlags(
			allowedSpecialSymbolsFlag,
			sensitiveKeywordsFlag,
			rules,
		)
		UseConfig(cfg)
	}

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

			partsArr := make([][]ExprPart, 0, len(call.Args))

			ln := 0
			for _, arg := range call.Args {
				argParts := getExpressionParts(arg, nil)
				argParts = foldConstantStrings(argParts)
				ln += len(argParts)
				partsArr = append(partsArr, argParts)
			}

			parts := make([]ExprPart, 0, ln)
			for _, argParts := range partsArr {
				parts = append(parts, argParts...)
			}

			p := newPassWrapper(pass)

			if config.rulesSet.shouldStartWithLowercase {
				checkShouldStartWithLowercase(parts, p)
			}
			if config.rulesSet.shouldContainOnlyEnglish {
				checkShouldContainOnlyEnglish(parts, p)
			}
			if config.rulesSet.shouldNotContainSpecialSymbols {
				checkShouldNotContainSpecialSymbols(parts, p)
			}
			if config.rulesSet.shouldNotContainSensitiveInformation {
				checkShouldNotContainSensitiveInformation(parts, p)
			}

			return true
		})
	}
	return nil, nil
}
