package analyzer

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestAnalyzerBasic(t *testing.T) {
	testdata := analysistest.TestData()

	UseConfig(DefaultConfig())

	analysistest.Run(t, testdata, Analyzer, "basic")
}

func TestAnalyzerAdvanced(t *testing.T) {
	testdata := analysistest.TestData()

	UseConfig(DefaultConfig())

	analysistest.Run(t, testdata, Analyzer, "advanced")
}
