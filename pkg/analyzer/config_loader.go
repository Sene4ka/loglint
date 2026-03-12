package analyzer

import (
	"flag"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

type FileConfig struct {
	AllowedSpecialSymbols string    `yaml:"allowed_special_symbols"`
	SensitiveKeywords     []string  `yaml:"sensitive_keywords"`
	Rules                 RuleFlags `yaml:"rules"`
}

type RuleFlags struct {
	ShouldStartWithLowercase             *bool `yaml:"rule-lowercase"`
	ShouldContainOnlyEnglish             *bool `yaml:"rule-english"`
	ShouldNotContainSpecialSymbols       *bool `yaml:"rule-symbols"`
	ShouldNotContainSensitiveInformation *bool `yaml:"rule-sensitive"`
}

func coalesce(ptr *bool, defaultVal bool) bool {
	if ptr != nil {
		return *ptr
	}
	return defaultVal
}

func DefaultRuleFlags() RuleFlags {
	t := true
	return RuleFlags{
		ShouldStartWithLowercase:             &t,
		ShouldContainOnlyEnglish:             &t,
		ShouldNotContainSpecialSymbols:       &t,
		ShouldNotContainSensitiveInformation: &t,
	}
}

func parseAllowedSymbols(s string) []rune {
	if s == "" {
		return []rune{':', '_', '=', '%'} // дефолт
	}

	var runes []rune
	for _, token := range strings.Fields(s) {
		for _, r := range token {
			runes = append(runes, r)
			break
		}
	}

	if len(runes) == 0 {
		return []rune{':', '_', '=', '%'}
	}
	return runes
}

func LoadConfigFromFile(path string) (*Config, error) {
	if path == "" {
		if env := os.Getenv("LOGLINT_CONFIG"); env != "" {
			path = env
		} else {
			cwd, _ := os.Getwd()
			path = filepath.Join(cwd, ".loglint.yml")
		}
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return DefaultConfig(), nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var fc FileConfig
	if err := yaml.Unmarshal(data, &fc); err != nil {
		return nil, err
	}

	allowedSymbols := parseAllowedSymbols(fc.AllowedSpecialSymbols)

	keywords := fc.SensitiveKeywords
	if len(keywords) == 0 {
		keywords = []string{"key", "password", "secret", "auth", "token"}
	}

	rules := fc.Rules

	return NewConfig(
		allowedSymbols,
		keywords,
		map[string]bool{
			"shouldStartWithLowercase":             coalesce(rules.ShouldStartWithLowercase, true),
			"shouldContainOnlyEnglish":             coalesce(rules.ShouldContainOnlyEnglish, true),
			"shouldNotContainSpecialSymbols":       coalesce(rules.ShouldNotContainSpecialSymbols, true),
			"shouldNotContainSensitiveInformation": coalesce(rules.ShouldNotContainSensitiveInformation, true),
		},
	), nil
}

func RegisterConfigFlag(fs *flag.FlagSet) {
	fs.String("config", "", "Path to .loglint.yml config file")
}

func LoadConfigFromFlags(
	allowedSymbols string,
	keywords string,
	rules map[string]bool,
) *Config {
	allowedRunes := parseAllowedSymbols(allowedSymbols)
	kwList := strings.Split(keywords, ",")
	for i := range kwList {
		kwList[i] = strings.TrimSpace(kwList[i])
	}

	return NewConfig(allowedRunes, kwList, rules)
}
