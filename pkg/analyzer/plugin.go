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
		Lowercase bool `json:"rule-lowercase"`
		English   bool `json:"rule-english"`
		Symbols   bool `json:"rule-symbols"`
		Sensitive bool `json:"rule-sensitive"`
	} `json:"rules"`
}

type Plugin struct {
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
			"shouldStartWithLowercase":             s.Rules.Lowercase,
			"shouldContainOnlyEnglish":             s.Rules.English,
			"shouldNotContainSpecialSymbols":       s.Rules.Symbols,
			"shouldNotContainSensitiveInformation": s.Rules.Sensitive,
		},
	)

	UseConfig(cfg)

	return &Plugin{cfg: cfg}, nil
}

func (p *Plugin) BuildAnalyzers() ([]*analysis.Analyzer, error) {
	return []*analysis.Analyzer{Analyzer}, nil
}

func (p *Plugin) GetLoadMode() string {
	return register.LoadModeTypesInfo
}
