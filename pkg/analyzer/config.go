package analyzer

import (
	"regexp"
)

type Config struct {
	AllowedSpecialSymbols  []rune
	allowedSymbolsMap      map[rune]bool
	SensitiveKeywords      []string
	sensitiveKeywordsRegex []*regexp.Regexp
}

func NewConfig(allowedSpecialSymbols []rune, sensitiveKeywords []string) *Config {
	symbolsMap := make(map[rune]bool, len(allowedSpecialSymbols))
	for _, symbol := range allowedSpecialSymbols {
		symbolsMap[symbol] = true
	}

	keywordsRegexp := make([]*regexp.Regexp, 0, len(sensitiveKeywords))
	for _, keyword := range sensitiveKeywords {
		escaped := regexp.QuoteMeta(keyword)
		pattern := `(?i)\b` + escaped + `\b`
		keywordsRegexp = append(keywordsRegexp, regexp.MustCompile(pattern))
	}
	return &Config{
		AllowedSpecialSymbols:  allowedSpecialSymbols,
		allowedSymbolsMap:      symbolsMap,
		SensitiveKeywords:      sensitiveKeywords,
		sensitiveKeywordsRegex: keywordsRegexp,
	}
}

var config Config

func UseConfig(cfg *Config) {
	config = *cfg
}
