package generator

import (
	"reflect"
	"testing"
)

func TestResolveStrategy(t *testing.T) {
	tests := []struct {
		name         string
		strategyName string
		wantType     string
	}{
		{
			name:         "overwrite strategy",
			strategyName: "overwrite",
			wantType:     "*generator.OverwriteStrategy",
		},
		{
			name:         "ignore strategy",
			strategyName: "ignore",
			wantType:     "<nil>",
		},
		{
			name:         "unknown strategy",
			strategyName: "unknown",
			wantType:     "<nil>",
		},
		{
			name:         "empty strategy",
			strategyName: "",
			wantType:     "<nil>",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			strategy := ResolveStrategy(tt.strategyName)
			var gotType string
			if strategy == nil {
				gotType = "<nil>"
			} else {
				gotType = reflect.TypeOf(strategy).String()
			}
			if gotType != tt.wantType {
				t.Errorf("ResolveStrategy() type = %v, want %v", gotType, tt.wantType)
			}
		})
	}
}
