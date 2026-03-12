package loglint

import (
	"strings"

	"github.com/creasty/defaults"
	"github.com/golangci/plugin-module-register/register"
	"golang.org/x/tools/go/analysis"
)

func init() {
	register.Plugin("loglint", New)
}

type PluginSettings struct {
	AllowedSpecialSymbols string   `json:"allowed_special_symbols" default:": _ - = %"`
	SensitiveKeywords     []string `json:"sensitive_keywords" default:"[\"key\",\"password\",\"secret\",\"auth\",\"token\"]"`
	Rules                 struct {
		RuleLowercase bool `json:"rule_lowercase" default:"true"`
		RuleEnglish   bool `json:"rule_english" default:"true"`
		RuleSymbols   bool `json:"rule_symbols" default:"true"`
		RuleSensitive bool `json:"rule_sensitive" default:"true"`
	} `json:"rules"`
}

type PluginLoglint struct {
	cfg *Config
}

func New(settings any) (register.LinterPlugin, error) {
	s, err := register.DecodeSettings[PluginSettings](settings)
	if err != nil {
		return nil, err
	}

	if err = defaults.Set(&s); err != nil {
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

	return &PluginLoglint{cfg: cfg}, nil
}

func (p *PluginLoglint) BuildAnalyzers() ([]*analysis.Analyzer, error) {
	return []*analysis.Analyzer{Analyzer}, nil
}

func (p *PluginLoglint) GetLoadMode() string {
	return register.LoadModeTypesInfo
}
