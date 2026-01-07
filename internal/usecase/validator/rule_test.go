package validator

import (
	"testing"

	"github.com/mew-ton/kex/internal/domain"
)

func TestIDRequiredRule_Validate(t *testing.T) {
	rule := &IDRequiredRule{}

	tests := []struct {
		name    string
		doc     *domain.Document
		wantErr bool
	}{
		{
			name:    "valid when ID is present",
			doc:     &domain.Document{ID: "doc-1"},
			wantErr: false,
		},
		{
			name:    "invalid when ID is empty",
			doc:     &domain.Document{ID: ""},
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

func TestFilenameMatchRule_Validate(t *testing.T) {
	rule := &FilenameMatchRule{}

	tests := []struct {
		name    string
		doc     *domain.Document
		wantErr bool
	}{
		{
			name:    "valid when filename matches ID",
			doc:     &domain.Document{ID: "doc-1", Path: "/path/to/doc-1.md"},
			wantErr: false,
		},
		{
			name:    "invalid when filename does not match ID",
			doc:     &domain.Document{ID: "doc-1", Path: "/path/to/other.md"},
			wantErr: true,
		},
		{
			name: "invalid when ID is empty (match empty basename?)",
			// ID empty is invalid here because basename("foo.md") != ""
			// If both empty? Path "" -> basename ".". ID "". "." != ""
			doc:     &domain.Document{ID: "", Path: "doc.md"},
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
