package entity

import "time"

// Prompt represents the original user instruction that led to a memory.
type Prompt struct {
	ID        string    `json:"id"`
	Input     string    `json:"input"`
	SessionID string    `json:"session_id,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}
