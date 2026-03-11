package analyzer

import (
	"go/token"
	"regexp"
	"unicode"
	"unicode/utf8"

	"golang.org/x/tools/go/analysis"
)

var keywords []string = []string{"key", "password", "token", "secret", "api_key", "auth"}

var keywordsRegexp []*regexp.Regexp

func init() {
	for _, word := range keywords {
		escaped := regexp.QuoteMeta(word)
		pattern := `(?i).*` + escaped + `.*`
		keywordsRegexp = append(keywordsRegexp, regexp.MustCompile(pattern))
	}
}

func checkLowercaseStart(msg string, pos token.Pos, pass *analysis.Pass) {
	if len(msg) == 0 {
		return
	}

	r, size := utf8.DecodeRuneInString(msg)
	if r == utf8.RuneError && size == 1 {
		pass.Reportf(pos, "%s: invalid UTF-8", msg)
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

func checkOnlyLatinLetters(msg string, pos token.Pos, pass *analysis.Pass) {
	if len(msg) == 0 {
		return
	}

	for i, r := range msg {
		if unicode.IsLetter(r) && !unicode.Is(unicode.Latin, r) {
			charPos := pos + token.Pos(i)
			diag := analysis.Diagnostic{
				Pos:      charPos,
				Category: "Warning",
				Message:  "log message should be in English",
			}
			pass.Report(diag)
			return
		}
	}
}

func checkNoSpecialSymbols(msg string, pos token.Pos, pass *analysis.Pass) {
	for i, r := range msg {
		if !unicode.In(r, unicode.Letter, unicode.Digit, unicode.Space) {
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

func checkNoKeywords(msg string, pos token.Pos, pass *analysis.Pass) {
	for _, keywordRegexp := range keywordsRegexp {
		if matches := keywordRegexp.FindAllStringIndex(msg, -1); len(matches) > 0 {
			for _, match := range matches {
				startPos := pos + token.Pos(match[0])
				endPos := pos + token.Pos(match[1])
				diag := analysis.Diagnostic{
					Pos:      startPos,
					End:      endPos,
					Category: "Warning",
					Message:  "log message should not contain sensitive information",
				}
				pass.Report(diag)
			}
		}
	}
}
