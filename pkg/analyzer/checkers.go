package analyzer

import (
	"fmt"
	"go/token"
	"unicode"
	"unicode/utf8"

	"golang.org/x/tools/go/analysis"
)

type Pass interface {
	ReportWarning(d analysis.Diagnostic)
}

type passWrapper struct {
	pass *analysis.Pass
}

func newPassWrapper(pass *analysis.Pass) passWrapper {
	return passWrapper{pass: pass}
}

func (p passWrapper) ReportWarning(d analysis.Diagnostic) {
	p.pass.Report(d)
}

func createLowercaseFix(pos token.Pos, msg string, size int) analysis.SuggestedFix {
	fixed := string(unicode.ToLower(rune(msg[0]))) + msg[size:]
	return analysis.SuggestedFix{
		Message: "Convert first letter to lowercase",
		TextEdits: []analysis.TextEdit{{
			Pos:     pos,
			End:     pos + token.Pos(size),
			NewText: []byte(fixed),
		}},
	}
}

func createRemoveNonLatinFix(pos token.Pos, msg string, byteIdx int) analysis.SuggestedFix {
	r, size := utf8.DecodeRuneInString(msg[byteIdx:])
	fixed := msg[:byteIdx] + msg[byteIdx+size:]
	return analysis.SuggestedFix{
		Message: fmt.Sprintf("Remove non-Latin character: %q", r),
		TextEdits: []analysis.TextEdit{{
			Pos:     pos + token.Pos(byteIdx),
			End:     pos + token.Pos(byteIdx+size),
			NewText: []byte(fixed),
		}},
	}
}

func createRemoveSpecialSymbolFix(pos token.Pos, msg string, byteIdx int) analysis.SuggestedFix {
	r, size := utf8.DecodeRuneInString(msg[byteIdx:])
	fixed := msg[:byteIdx] + msg[byteIdx+size:]
	return analysis.SuggestedFix{
		Message: fmt.Sprintf("Remove special symbol: %q", r),
		TextEdits: []analysis.TextEdit{{
			Pos:     pos + token.Pos(byteIdx),
			End:     pos + token.Pos(byteIdx+size),
			NewText: []byte(fixed),
		}},
	}
}

func createRedactFix(pos, end token.Pos) analysis.SuggestedFix {
	return analysis.SuggestedFix{
		Message: "Replace with [REDACTED]",
		TextEdits: []analysis.TextEdit{{
			Pos:     pos,
			End:     end,
			NewText: []byte("[REDACTED]"),
		}},
	}
}

func createRemoveVariableFix(pos, end token.Pos) analysis.SuggestedFix {
	return analysis.SuggestedFix{
		Message: "Remove the variable",
	}
}

/*
checkShouldStartWithLowercase checks only first ExprPart
if it is of type PartConst - check first rune
if it is other/variable type - then it is undetermined if it starts with uppercase or no
*/
func checkShouldStartWithLowercase(parts []ExprPart, pass Pass) {
	if len(parts) == 0 {
		return
	}

	part := parts[0]

	if part.Type != PartConst {
		return
	}

	msg := part.Value
	pos := part.Pos

	// should be merged before, so no other PartConst coming next, only other types
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

		diag.SuggestedFixes = []analysis.SuggestedFix{
			createLowercaseFix(pos, msg, size),
		}

		pass.ReportWarning(diag)
	}
}

// check that all direct and underlying const parts is English only
func checkShouldContainOnlyEnglish(parts []ExprPart, pass Pass) {
	for _, part := range parts {
		if part.Type == PartVar {
			continue
		}

		msg := part.Value
		pos := part.Pos
		if len(msg) == 0 {
			continue
		}

		for i, r := range msg {
			if unicode.IsLetter(r) && !unicode.Is(unicode.Latin, r) {
				charPos := pos + token.Pos(i)
				diag := analysis.Diagnostic{
					Pos:      charPos,
					Category: "Warning",
					Message:  "log message should contain only English letters",
					SuggestedFixes: []analysis.SuggestedFix{
						createRemoveNonLatinFix(pos, msg, i),
					},
				}
				pass.ReportWarning(diag)
				return
			}
		}
	}
}

// checks that all direct and underlying const parts do not contain special symbols
func checkShouldNotContainSpecialSymbols(parts []ExprPart, pass Pass) {
	for _, part := range parts {
		if part.Type == PartVar {
			continue
		}
		msg := part.Value
		pos := part.Pos
		if len(msg) == 0 {
			continue
		}
		for i, r := range msg {
			if !unicode.In(r, unicode.Letter, unicode.Digit, unicode.Space) && !config.allowedSymbolsMap[r] {
				charPos := pos + token.Pos(i)
				diag := analysis.Diagnostic{
					Pos:      charPos,
					Category: "Warning",
					Message:  "log message should not contain special symbols or emojis",
					SuggestedFixes: []analysis.SuggestedFix{
						createRemoveSpecialSymbolFix(pos, msg, i),
					},
				}
				pass.ReportWarning(diag)
				return
			}
		}
	}
}

/*
Checks that all const parts and variables do not contain sensitive keywords
Reports if sensitive keyword found in text as separate words and,
if found, also reports variables which contain, but might not be equal to, this keyword
*/
// Also separately checks for variables which names are equal to sensitive keywords and reports them
//
func checkShouldNotContainSensitiveInformation(parts []ExprPart, pass Pass) {
	var vars []ExprPart
	for _, part := range parts {
		if part.Type == PartVar {
			vars = append(vars, part)
		}
	}

	diagnosedVars := make([]bool, len(vars))

	for _, part := range parts {
		if part.Type != PartConst {
			continue
		}

		for _, r := range config.sensitiveKeywordsRegex {
			if matches := r.FindAllStringIndex(part.Value, -1); len(matches) > 0 && len(vars) > 0 {
				for _, match := range matches {
					start, end := match[0], match[1]

					if start > 0 {
						before, _ := utf8.DecodeLastRuneInString(part.Value[:start])
						first, _ := utf8.DecodeRuneInString(part.Value[start:])
						if unicode.IsLetter(before) && !unicode.IsUpper(first) {
							continue
						}
					}

					if end < len(part.Value) {
						after, _ := utf8.DecodeRuneInString(part.Value[end:])
						if unicode.IsLetter(after) && !unicode.IsUpper(after) {
							continue
						}
					}

					pass.ReportWarning(analysis.Diagnostic{
						Pos:      part.Pos + token.Pos(start),
						End:      part.Pos + token.Pos(end),
						Category: "Warning",
						Message:  "log message should not contain sensitive information",
						SuggestedFixes: []analysis.SuggestedFix{
							createRedactFix(part.Pos+token.Pos(start), part.End+token.Pos(end)),
						},
					})
				}
			}
		}
	}

	for i, v := range vars {
		if diagnosedVars[i] {
			continue
		}

		for _, r := range config.sensitiveKeywordsRegex {
			if r.MatchString(v.Value) {
				pass.ReportWarning(analysis.Diagnostic{
					Pos:      v.Pos,
					End:      v.End,
					Category: "Warning",
					Message:  "log variable name suggests sensitive data",
					SuggestedFixes: []analysis.SuggestedFix{
						createRemoveVariableFix(v.Pos, v.End),
					},
				})
			}
		}
	}
}
