package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/Sene4ka/log-linter/pkg/analyzer"

	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	analyzer.RegisterConfigFlag(flag.CommandLine)

	configPath := flag.Lookup("config").Value.String()
	cfg, err := analyzer.LoadConfigFromFile(configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "loglint: config error: %v\n", err)
		cfg = analyzer.DefaultConfig()
	}

	analyzer.UseConfig(cfg)

	singlechecker.Main(analyzer.Analyzer)
}
