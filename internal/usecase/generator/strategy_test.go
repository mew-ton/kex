package generator

import (
	"reflect"
	"testing"
)

func TestResolveStrategy(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		config   map[string]string
		wantType string
	}{
		{
			name:     "kex documentation overwrites",
			path:     "contents/documentation/kex/rule.md",
			wantType: "*generator.OverwriteStrategy",
		},
		{
			name:     "AGENTS.md uses marker strategy",
			path:     "AGENTS.md",
			wantType: "*generator.MarkerUpdateStrategy",
		},
		{
			name:     "CLAUDE.md uses marker strategy",
			path:     "CLAUDE.md",
			wantType: "*generator.MarkerUpdateStrategy",
		},
		{
			name:     "config override takes precedence",
			path:     "README.md",
			config:   map[string]string{"README.md": "append"},
			wantType: "*generator.AppendStrategy",
		},
		{
			name:     "default is skip",
			path:     "random.txt",
			wantType: "*generator.SkipStrategy",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			strategy := ResolveStrategy(tt.path, tt.config)
			gotType := reflect.TypeOf(strategy).String()
			if gotType != tt.wantType {
				t.Errorf("ResolveStrategy() type = %v, want %v", gotType, tt.wantType)
			}
		})
	}
}
