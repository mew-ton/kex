package generator

import (
	"reflect"
	"testing"

	"github.com/mew-ton/kex/internal/infrastructure/config"
)

func TestResolveStrategy(t *testing.T) {
	tests := []struct {
		name       string
		path       string
		strategies config.Strategies
		wantType   string
	}{
		{
			name:       "kex documentation overwrites",
			path:       "contents/documentation/kex/choose-effective-keywords.md",
			strategies: config.Strategies{ChooseEffectiveKeywords: "overwrite"},
			wantType:   "*generator.OverwriteStrategy",
		},
		{
			name:       "kex coding skip matches config (agent)",
			path:       ".agent/rules/kex-coding.md",
			strategies: config.Strategies{AgentKexCoding: "skip"},
			wantType:   "*generator.CreateStrategy",
		},
		{
			name:       "kex coding skip matches config (claude)",
			path:       ".claude/rules/kex/follow-coding-rules.md",
			strategies: config.Strategies{ClaudeKexCoding: "skip"},
			wantType:   "*generator.CreateStrategy",
		},

		{
			name:       "default fallback is overwrite",
			path:       "random.txt",
			strategies: config.Strategies{}, // Empty config matches nothing -> fallback
			wantType:   "*generator.OverwriteStrategy",
		},
		{
			name:       "unconfigured file defaults to overwrite match",
			path:       ".agent/rules/kex-coding.md",
			strategies: config.Strategies{AgentKexCoding: ""},
			// Empty string -> default case in switch (default overwrite)
			wantType: "*generator.OverwriteStrategy",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			strategy := ResolveStrategy(tt.path, tt.strategies)
			gotType := reflect.TypeOf(strategy).String()
			if gotType != tt.wantType {
				t.Errorf("ResolveStrategy() type = %v, want %v", gotType, tt.wantType)
			}
		})
	}
}
