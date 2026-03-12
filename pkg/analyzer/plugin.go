package analyzer

import (
	"strings"

	"github.com/golangci/plugin-module-register/register"
	"golang.org/x/tools/go/analysis"
)

func init() {
	register.Plugin("loglint", New)
}

type PluginSettings struct {
	AllowedSpecialSymbols string   `json:"allowed_special_symbols"`
	SensitiveKeywords     []string `json:"sensitive_keywords"`
	Rules                 struct {
		RuleLowercase bool `json:"rule_lowercase"`
		RuleEnglish   bool `json:"rule_english"`
		RuleSymbols   bool `json:"rule_symbols"`
		RuleSensitive bool `json:"rule_sensitive"`
	} `json:"rules"`
}

type loglintPlugin struct {
	cfg *Config
}

func New(settings any) (register.LinterPlugin, error) {
	s, err := register.DecodeSettings[PluginSettings](settings)
	if err != nil {
		return nil, err
	}

	cfg := LoadConfigFromFlags(
		s.AllowedSpecialSymbols,
		strings.Join(s.SensitiveKeywords, ","),
		map[string]bool{
			"shouldStartWithLowercase":             s.Rules.RuleLowercase,
			"shouldContainOnlyEnglish":             s.Rules.RuleEnglish,
			"shouldNotContainSpecialSymbols":       s.Rules.RuleSymbols,
			"shouldNotContainSensitiveInformation": s.Rules.RuleSensitive,
		},
	)

	UseConfig(cfg)

	return &loglintPlugin{cfg: cfg}, nil
}

func (p *loglintPlugin) BuildAnalyzers() ([]*analysis.Analyzer, error) {
	return []*analysis.Analyzer{Analyzer}, nil
}

func (p *loglintPlugin) GetLoadMode() string {
	return register.LoadModeTypesInfo
}
