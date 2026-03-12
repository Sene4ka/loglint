package main

import (
	"github.com/Sene4ka/log-linter/pkg/analyzer"

	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	config := analyzer.NewConfig(
		[]rune{':', '_', '='},
		[]string{"key", "password", "secret", "apiKey", "api_key", "auth", "token", "auth_token", "authToken"})
	analyzer.UseConfig(config)

	singlechecker.Main(analyzer.Analyzer)
}
