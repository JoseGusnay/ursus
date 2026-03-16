package entity

import "time"

const (
	ScopeProject  = "project"
	ScopePersonal = "personal"
)

// Ursus represents the core memory entity.
type Ursus struct {
	ID             string     `json:"id"`
	Content        string     `json:"content"`
	Metadata       string     `json:"metadata"`
	SessionID      string     `json:"session_id,omitempty"`
	TopicKey       string     `json:"topic_key,omitempty"`
	PromptID       string     `json:"prompt_id,omitempty"`
	Scope          string     `json:"scope,omitempty"`
	DuplicateCount int        `json:"duplicate_count,omitempty"`
	RevisionCount  int        `json:"revision_count,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	LastSeenAt     time.Time  `json:"last_seen_at,omitempty"`
	DeletedAt      *time.Time `json:"deleted_at,omitempty"`
}
