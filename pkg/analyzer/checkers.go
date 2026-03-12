package analyzer

import (
	"go/token"
	"strings"
	"unicode"
	"unicode/utf8"

	"golang.org/x/tools/go/analysis"
)

func checkShouldStartWithLowercase(parts []ExprPart, pass *analysis.Pass) {
	for _, part := range parts {
		if part.Type == PartVar {
			continue
		}

		msg := part.Value
		pos := part.Pos

		if len(msg) == 0 {
			continue
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
}

func checkShouldContainOnlyEnglish(parts []ExprPart, pass *analysis.Pass) {
	for _, part := range parts {
		if part.Type == PartVar {
			continue
		}

		msg := part.Value
		pos := part.Pos
		if len(msg) == 0 {
			return
		}

		for i, r := range msg {
			if unicode.IsLetter(r) && !unicode.Is(unicode.Latin, r) {
				charPos := pos + token.Pos(i)
				diag := analysis.Diagnostic{
					Pos:      charPos,
					Category: "Warning",
					Message:  "log message should contain only English letters",
				}
				pass.Report(diag)
				return
			}
		}
	}
}

func checkShouldNotContainSpecialSymbols(parts []ExprPart, pass *analysis.Pass) {
	for _, part := range parts {
		if part.Type == PartVar {
			continue
		}
		msg := part.Value
		pos := part.Pos
		if len(msg) == 0 {
			return
		}
		for i, r := range msg {
			if !unicode.In(r, unicode.Letter, unicode.Digit, unicode.Space) && !config.allowedSymbolsMap[r] {
				charPos := pos + token.Pos(i)
				diag := analysis.Diagnostic{
					Pos:      charPos,
					Category: "Warning",
					Message:  "log message should not contain special symbols or emojis",
				}
				pass.Report(diag)
				return
			}
		}
	}
}

func checkShouldNotContainSensitiveInformation(parts []ExprPart, pass *analysis.Pass) {
	varsCnt := 0
	for _, part := range parts {
		if part.Type == PartVar {
			varsCnt++
		}
	}

	vars := make([]ExprPart, 0, varsCnt)
	diagnosedVars := make([]bool, varsCnt)
	varIdx := make(map[int]int, varsCnt)
	varReverseIdx := make(map[int]int, varsCnt)

	for i, part := range parts {
		if part.Type == PartVar {
			varIdx[i] = len(vars)
			varReverseIdx[len(vars)] = i
			vars = append(vars, part)
		}
	}

	for i, part := range parts {
		if part.Type == PartConst {
			for kIndex, r := range config.sensitiveKeywordsRegex {
				if matches := r.FindAllStringIndex(part.Value, -1); len(matches) > 0 && varsCnt > 0 {
					for _, match := range matches {
						pass.Report(analysis.Diagnostic{
							Pos:      part.Pos + token.Pos(match[0]),
							End:      part.Pos + token.Pos(match[1]),
							Category: "Warning",
							Message:  "log message should not contain sensitive information",
						})
					}
					for j, v := range vars {
						if !diagnosedVars[j] {
							found := strings.Contains(v.Value, config.SensitiveKeywords[kIndex])
							if found {
								diagnosedVars[j] = true
								pass.Report(analysis.Diagnostic{
									Pos:      part.Pos,
									End:      part.End,
									Category: "Warning",
									Message:  "log message should not contain sensitive information",
								})
							}
						}
					}
				}
			}
		} else if part.Type == PartVar {
			if res, ok := varIdx[i]; ok && !diagnosedVars[res] {
				for _, kw := range config.SensitiveKeywords {
					if part.Value == kw {
						pass.Report(analysis.Diagnostic{
							Pos:      part.Pos,
							End:      part.End,
							Category: "Warning",
							Message:  "log message should not contain sensitive information",
						})
					}
				}
			}
		}
	}
}
