package loglint

import (
	"testing"

	"github.com/golangci/plugin-module-register/register"
	"github.com/stretchr/testify/require"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestAnalyzerBasic(t *testing.T) {
	testdata := analysistest.TestData()

	newPlugin, err := register.GetPlugin("loglint")
	require.NoError(t, err)

	plugin, err := newPlugin(nil)
	require.NoError(t, err)

	analyzers, err := plugin.BuildAnalyzers()
	require.NoError(t, err)

	analysistest.Run(t, testdata, analyzers[0], "basic")
}

func TestAnalyzerAdvanced(t *testing.T) {
	testdata := analysistest.TestData()

	newPlugin, err := register.GetPlugin("loglint")
	require.NoError(t, err)

	plugin, err := newPlugin(nil)
	require.NoError(t, err)

	analyzers, err := plugin.BuildAnalyzers()
	require.NoError(t, err)

	analysistest.Run(t, testdata, analyzers[0], "advanced")
}
