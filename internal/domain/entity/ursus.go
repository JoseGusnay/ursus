package entity

import "time"

// Ursus represents the core memory entity.
type Ursus struct {
	ID        string    `json:"id"`
	Content   string    `json:"content"`
	Metadata  string    `json:"metadata"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
