package service

import (
	"regexp"
)

var privateRegex = regexp.MustCompile(`(?s)<private>.*?</private>`)

// PrivacyService handles sensitive data redaction.
type PrivacyService struct{}

// NewPrivacyService creates a new PrivacyService instance.
func NewPrivacyService() *PrivacyService {
	return &PrivacyService{}
}

// Redact replaces any content within <private> tags with [REDACTED].
func (s *PrivacyService) Redact(content string) string {
	return privateRegex.ReplaceAllString(content, "[REDACTED]")
}
