package analyzer

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestAnalyzer(t *testing.T) {
	testdata := analysistest.TestData()

	config := NewConfig(
		[]rune{':', '_', '='},
		[]string{"key", "password", "secret", "apiKey", "api_key", "auth", "token", "auth_token", "authToken"})
	UseConfig(config)

	analysistest.Run(t, testdata, Analyzer, "example")
}
