package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/Sene4ka/loglint/pkg/loglint"

	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	loglint.RegisterConfigFlag(flag.CommandLine)

	configPath := flag.Lookup("config").Value.String()
	cfg, err := loglint.LoadConfigFromFile(configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "loglint: config error: %v\n", err)
		cfg = loglint.DefaultConfig()
	}

	loglint.UseConfig(cfg)

	singlechecker.Main(loglint.Analyzer)
}
