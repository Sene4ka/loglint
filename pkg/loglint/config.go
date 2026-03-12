package loglint

import (
	"regexp"
)

type Config struct {
	AllowedSpecialSymbols  []rune
	allowedSymbolsMap      map[rune]bool
	SensitiveKeywords      []string
	sensitiveKeywordsRegex []*regexp.Regexp
	Rules                  map[string]bool
	rulesSet               ruleSet
}

type ruleSet struct {
	shouldStartWithLowercase             bool
	shouldContainOnlyEnglish             bool
	shouldNotContainSpecialSymbols       bool
	shouldNotContainSensitiveInformation bool
}

func NewConfig(allowedSpecialSymbols []rune, sensitiveKeywords []string, rules map[string]bool) *Config {
	symbolsMap := make(map[rune]bool, len(allowedSpecialSymbols))
	for _, symbol := range allowedSpecialSymbols {
		symbolsMap[symbol] = true
	}

	keywordsRegexp := make([]*regexp.Regexp, 0, len(sensitiveKeywords))
	for _, keyword := range sensitiveKeywords {
		escaped := regexp.QuoteMeta(keyword)
		pattern := `(?i)` + escaped
		keywordsRegexp = append(keywordsRegexp, regexp.MustCompile(pattern))
	}

	rulesSet := ruleSet{
		shouldStartWithLowercase:             rules["shouldStartWithLowercase"],
		shouldContainOnlyEnglish:             rules["shouldContainOnlyEnglish"],
		shouldNotContainSpecialSymbols:       rules["shouldNotContainSpecialSymbols"],
		shouldNotContainSensitiveInformation: rules["shouldNotContainSensitiveInformation"],
	}

	return &Config{
		AllowedSpecialSymbols:  allowedSpecialSymbols,
		allowedSymbolsMap:      symbolsMap,
		SensitiveKeywords:      sensitiveKeywords,
		sensitiveKeywordsRegex: keywordsRegexp,
		Rules:                  rules,
		rulesSet:               rulesSet,
	}
}

func DefaultConfig() *Config {
	return NewConfig(
		[]rune{':', '_', '-', '=', '%'},
		[]string{"key", "password", "secret", "auth", "token"},
		map[string]bool{
			"shouldStartWithLowercase":             true,
			"shouldContainOnlyEnglish":             true,
			"shouldNotContainSpecialSymbols":       true,
			"shouldNotContainSensitiveInformation": true,
		})
}

var config *Config

func UseConfig(cfg *Config) {
	if config != nil {
		return
	}
	config = cfg
}
