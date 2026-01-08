package validator

import (
	"testing"

	"github.com/mew-ton/kex/internal/domain"
)

func TestTitleRequiredRule_Validate(t *testing.T) {
	rule := &TitleRequiredRule{}

	tests := []struct {
		name    string
		doc     *domain.Document
		wantErr bool
	}{
		{
			name:    "valid when Title is present",
			doc:     &domain.Document{Title: "Some Title"},
			wantErr: false,
		},
		{
			name:    "invalid when Title is empty",
			doc:     &domain.Document{Title: ""},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := rule.Validate(tt.doc); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
