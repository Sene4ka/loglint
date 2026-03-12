package analyzer

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestAnalyzer(t *testing.T) {
	testdata := analysistest.TestData()

	testConfig := NewConfig(
		[]rune{':', '_', '='},
		[]string{"key", "password", "secret", "apiKey", "api_key", "auth", "token", "auth_token", "authToken"},
		map[string]bool{
			"shouldStartWithLowercase":             true,
			"shouldContainOnlyEnglish":             true,
			"shouldNotContainSpecialSymbols":       true,
			"shouldNotContainSensitiveInformation": true,
		},
	)
	UseConfig(testConfig)

	analysistest.Run(t, testdata, Analyzer, "example")
}
