package service

import (
	"regexp"
)

var (
	privateRegex = regexp.MustCompile(`(?s)<private>.*?</private>`)
	emailRegex   = regexp.MustCompile(`[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}`)
	// Pattern for potential API keys or tokens (e.g., sk-..., secret_..., token-...)
	secretRegex  = regexp.MustCompile(`(?i)(sk-|secret_|token-|key-)[a-z0-9]{8,}`)
)

// PrivacyService handles sensitive data redaction.
type PrivacyService struct{}

// NewPrivacyService creates a new PrivacyService instance.
func NewPrivacyService() *PrivacyService {
	return &PrivacyService{}
}

// Redact replaces any sensitive content with [REDACTED].
func (s *PrivacyService) Redact(content string) string {
	content = privateRegex.ReplaceAllString(content, "[REDACTED]")
	content = emailRegex.ReplaceAllString(content, "[EMAIL_REDACTED]")
	content = secretRegex.ReplaceAllString(content, "[SECRET_REDACTED]")
	return content
}
