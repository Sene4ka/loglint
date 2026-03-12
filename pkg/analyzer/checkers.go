package analyzer

import (
	"go/token"
	"unicode"
	"unicode/utf8"

	"golang.org/x/tools/go/analysis"
)

func checkLowercaseStart(parts []VarPart, pass *analysis.Pass) {
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

func checkOnlyLatinLetters(parts []VarPart, pass *analysis.Pass) {
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

func checkNoSpecialSymbols(parts []VarPart, pass *analysis.Pass) {
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

func checkNoSensitiveKeywordsAndVariables(parts []VarPart, pass *analysis.Pass) {
	for i, part := range parts {
		if part.Type == PartConst {
			for _, r := range config.sensitiveKeywordsRegex {
				if matches := r.FindAllStringIndex(part.Value, -1); len(matches) > 0 && i+1 < len(parts) && parts[i+1].Type == PartVar {
					for _, match := range matches {
						pass.Report(analysis.Diagnostic{
							Pos:      part.Pos + token.Pos(match[0]),
							End:      part.Pos + token.Pos(match[1]),
							Category: "Warning",
							Message:  "log message should not contain sensitive information",
						})
					}
				}
			}
		} else if part.Type == PartVar {
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
