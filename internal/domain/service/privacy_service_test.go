package service

import (
	"testing"
)

func TestPrivacyService_Redact(t *testing.T) {
	s := NewPrivacyService()

	tests := []struct {
		name    string
		content string
		want    string
	}{
		{
			name:    "No tags",
			content: "Normal text without secrets",
			want:    "Normal text without secrets",
		},
		{
			name:    "Single tag",
			content: "My secret is <private>12345</private>",
			want:    "My secret is [REDACTED]",
		},
		{
			name:    "Multiple tags",
			content: "User: <private>Jose</private>, Pass: <private>pass123</private>",
			want:    "User: [REDACTED], Pass: [REDACTED]",
		},
		{
			name:    "Multiline tag",
			content: "Private block:\n<private>\nline 1\nline 2\n</private>",
			want:    "Private block:\n[REDACTED]",
		},
		{
			name:    "Nested-like tags (simple match)",
			content: "<private>outer <private>inner</private></private>",
			want:    "[REDACTED]</private>", // Regular regex is non-recursive, it matches first </private>
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := s.Redact(tt.content); got != tt.want {
				t.Errorf("Redact() = %v, want %v", got, tt.want)
			}
		})
	}
}
