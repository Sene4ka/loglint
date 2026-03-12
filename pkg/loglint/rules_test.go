package loglint

import (
	"go/token"
	"testing"

	"golang.org/x/tools/go/analysis"
)

type mockPass struct {
	diagnostics []analysis.Diagnostic
}

func (m *mockPass) ReportWarning(d analysis.Diagnostic) {
	m.diagnostics = append(m.diagnostics, d)
}

func constPart(value string, pos, end token.Pos) ExprPart {
	return ExprPart{Value: value, Pos: pos, End: end, Type: PartConst}
}

func varPart(value string, pos, end token.Pos) ExprPart {
	return ExprPart{Value: value, Pos: pos, End: end, Type: PartVar}
}

func otherPart(value string, pos, end token.Pos) ExprPart {
	return ExprPart{Value: value, Pos: pos, End: end, Type: PartOther}
}

func assertDiagnosticsCount(t *testing.T, pass *mockPass, expected int) {
	t.Helper()
	if len(pass.diagnostics) != expected {
		t.Errorf("expected %d diagnostics, got %d:\n%+v", expected, len(pass.diagnostics), pass.diagnostics)
	}
}

func assertDiagnosticMessage(t *testing.T, diag analysis.Diagnostic, expected string) {
	t.Helper()
	if diag.Message != expected {
		t.Errorf("expected message %q, got %q", expected, diag.Message)
	}
}

func assertDiagnosticPosition(t *testing.T, diag analysis.Diagnostic, expectedPos token.Pos) {
	t.Helper()
	if diag.Pos != expectedPos {
		t.Errorf("expected Pos %v, got %v", expectedPos, diag.Pos)
	}
}

func assertHasSuggestedFix(t *testing.T, diag analysis.Diagnostic) {
	t.Helper()
	if len(diag.SuggestedFixes) == 0 {
		t.Error("expected SuggestedFixes, got none")
	}
}

func TestCheckShouldStartWithLowercase(t *testing.T) {
	UseConfig(DefaultConfig())

	tests := []struct {
		name        string
		parts       []ExprPart
		wantCount   int
		wantMessage string
		wantPos     token.Pos
		wantFix     bool
	}{
		{
			name:        "uppercase first letter in const",
			parts:       []ExprPart{constPart("Starting server", 10, 25)},
			wantCount:   1,
			wantMessage: "log message should start with lowercase letter",
			wantPos:     10,
			wantFix:     true,
		},
		{
			name:      "lowercase first letter - OK",
			parts:     []ExprPart{constPart("starting server", 10, 25)},
			wantCount: 0,
		},
		{
			name:      "starts with variable - undetermined, skip",
			parts:     []ExprPart{varPart("server", 10, 16)},
			wantCount: 0,
		},
		{
			name:      "empty const value",
			parts:     []ExprPart{constPart("", 10, 10)},
			wantCount: 0,
		},
		{
			name:        "Cyrillic uppercase - should flag",
			parts:       []ExprPart{constPart("Запуск сервера", 10, 30)},
			wantCount:   1,
			wantMessage: "log message should start with lowercase letter",
			wantPos:     10,
			wantFix:     true,
		},
		{
			name:      "emoji at start - not a letter, skip",
			parts:     []ExprPart{constPart("🚀 Starting", 10, 25)},
			wantCount: 0,
		},
		{
			name:      "var at start - not a const, skip",
			parts:     []ExprPart{varPart("someVar", 10, 25)},
			wantCount: 0,
		},
		{
			name:        "multiple parts - only first const checked",
			parts:       []ExprPart{constPart("Starting", 10, 18), varPart("server", 20, 26)},
			wantCount:   1,
			wantMessage: "log message should start with lowercase letter",
			wantPos:     10,
			wantFix:     true,
		},
		{
			name:      "PartOther at start - skip",
			parts:     []ExprPart{otherPart("", 10, 10)},
			wantCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pass := &mockPass{}
			checkShouldStartWithLowercase(tt.parts, pass)
			assertDiagnosticsCount(t, pass, tt.wantCount)
			if tt.wantCount > 0 {
				d := pass.diagnostics[0]
				assertDiagnosticMessage(t, d, tt.wantMessage)
				assertDiagnosticPosition(t, d, tt.wantPos)
				if tt.wantFix {
					assertHasSuggestedFix(t, d)
				}
			}
		})
	}
}

func TestCheckShouldContainOnlyEnglish(t *testing.T) {
	UseConfig(DefaultConfig())

	tests := []struct {
		name        string
		parts       []ExprPart
		wantCount   int
		wantMessage string
		wantPos     token.Pos
	}{
		{
			name:      "english only - OK",
			parts:     []ExprPart{constPart("server started", 10, 25)},
			wantCount: 0,
		},
		{
			name:        "Cyrillic letter - flag first non-Latin",
			parts:       []ExprPart{constPart("запуск сервера", 10, 30)},
			wantCount:   1,
			wantMessage: "log message should contain only English letters",
			wantPos:     10,
		},
		{
			name:        "mixed english+cyrillic - flag at cyrillic start",
			parts:       []ExprPart{constPart("server запуск", 10, 25)},
			wantCount:   1,
			wantMessage: "log message should contain only English letters",
			wantPos:     17,
		},
		{
			name:      "variable with cyrillic - skipped",
			parts:     []ExprPart{varPart("запуск", 10, 20)},
			wantCount: 0,
		},
		{
			name:      "digits and spaces - OK",
			parts:     []ExprPart{constPart("server123 test", 10, 25)},
			wantCount: 0,
		},
		{
			name:        "multiple consts - stops at first violation",
			parts:       []ExprPart{constPart("ok ", 10, 13), varPart("someVar", 15, 19), constPart("запуск", 21, 35)},
			wantCount:   1,
			wantMessage: "log message should contain only English letters",
			wantPos:     21,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pass := &mockPass{}
			checkShouldContainOnlyEnglish(tt.parts, pass)
			assertDiagnosticsCount(t, pass, tt.wantCount)
			if tt.wantCount > 0 {
				d := pass.diagnostics[0]
				assertDiagnosticMessage(t, d, tt.wantMessage)
				assertDiagnosticPosition(t, d, tt.wantPos)
			}
		})
	}
}

func TestCheckShouldNotContainSpecialSymbols(t *testing.T) {
	UseConfig(DefaultConfig())

	tests := []struct {
		name        string
		parts       []ExprPart
		wantCount   int
		wantMessage string
		wantPos     token.Pos
	}{
		{
			name:      "allowed symbols only - OK",
			parts:     []ExprPart{constPart("api_key=value:100%", 10, 28)},
			wantCount: 0,
		},
		{
			name:        "emoji - flag",
			parts:       []ExprPart{constPart("server started 🚀", 10, 28)},
			wantCount:   1,
			wantMessage: "log message should not contain special symbols or emojis",
			wantPos:     25,
		},
		{
			name:        "exclamation marks - flag first",
			parts:       []ExprPart{constPart("failed!!!", 10, 19)},
			wantCount:   1,
			wantMessage: "log message should not contain special symbols or emojis",
			wantPos:     16,
		},
		{
			name:        "parentheses - flag",
			parts:       []ExprPart{constPart("server (prod)", 10, 24)},
			wantCount:   1,
			wantMessage: "log message should not contain special symbols or emojis",
			wantPos:     17,
		},
		{
			name:      "variable with emoji - skipped",
			parts:     []ExprPart{varPart("🚀", 10, 14)},
			wantCount: 0,
		},
		{
			name:      "letters+digits+space - OK",
			parts:     []ExprPart{constPart("server 123 test", 10, 25)},
			wantCount: 0,
		},
		{
			name:        "multiple consts - stops at first violation",
			parts:       []ExprPart{constPart("ok ", 10, 13), constPart("fail!!!", 15, 22)},
			wantCount:   1,
			wantMessage: "log message should not contain special symbols or emojis",
			wantPos:     19,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pass := &mockPass{}
			checkShouldNotContainSpecialSymbols(tt.parts, pass)
			assertDiagnosticsCount(t, pass, tt.wantCount)
			if tt.wantCount > 0 {
				d := pass.diagnostics[0]
				assertDiagnosticMessage(t, d, tt.wantMessage)
				assertDiagnosticPosition(t, d, tt.wantPos)
			}
		})
	}
}

func TestCheckShouldNotContainSensitiveInformation(t *testing.T) {
	UseConfig(DefaultConfig())

	tests := []struct {
		name          string
		parts         []ExprPart
		wantCount     int
		wantMessages  []string
		wantPositions []token.Pos
	}{
		{
			name:      "keyword in const + variable with matching name",
			parts:     []ExprPart{constPart("auth: ", 10, 16), varPart("authToken", 20, 29)},
			wantCount: 3,
			wantMessages: []string{
				"log message should not contain sensitive information",
				"log variable name suggests sensitive data",
				"log variable name suggests sensitive data",
			},
			wantPositions: []token.Pos{10, 20},
		},
		{
			name:      "keyword alone without variables - no report",
			parts:     []ExprPart{constPart("password ok", 10, 21)},
			wantCount: 0,
		},
		{
			name:          "variable name equals keyword",
			parts:         []ExprPart{varPart("password", 10, 18)},
			wantCount:     1,
			wantMessages:  []string{"log variable name suggests sensitive data"},
			wantPositions: []token.Pos{10},
		},
		{
			name:          "variable name contains keyword (case-insensitive)",
			parts:         []ExprPart{varPart("UserPassword", 10, 22)},
			wantCount:     1,
			wantMessages:  []string{"log variable name suggests sensitive data"},
			wantPositions: []token.Pos{10},
		},
		{
			name:      "keyword inside word - no boundary, skip",
			parts:     []ExprPart{constPart("mypassword", 10, 20), varPart("x", 25, 26)},
			wantCount: 0,
		},
		{
			name:      "CamelCase: keyword+Upper in var name",
			parts:     []ExprPart{constPart("Token: ", 10, 17), varPart("accessToken", 20, 31)},
			wantCount: 2,
			wantMessages: []string{
				"log message should not contain sensitive information",
				"log variable name suggests sensitive data",
			},
			wantPositions: []token.Pos{10, 20},
		},
		{
			name:      "split keyword in const + matching var",
			parts:     []ExprPart{constPart("ap"+"i_key=", 10, 18), varPart("apiKey", 20, 26)},
			wantCount: 2,
			wantMessages: []string{
				"log message should not contain sensitive information",
				"log variable name suggests sensitive data",
			},
			wantPositions: []token.Pos{14, 20},
		},
		{
			name:      "no variables - keyword alone is OK",
			parts:     []ExprPart{constPart("just a password mention", 10, 33)},
			wantCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pass := &mockPass{}
			checkShouldNotContainSensitiveInformation(tt.parts, pass)
			assertDiagnosticsCount(t, pass, tt.wantCount)
			for i, expectedMsg := range tt.wantMessages {
				if i < len(pass.diagnostics) {
					assertDiagnosticMessage(t, pass.diagnostics[i], expectedMsg)
					if i < len(tt.wantPositions) {
						assertDiagnosticPosition(t, pass.diagnostics[i], tt.wantPositions[i])
					}
				}
			}
		})
	}
}
