package entity

import "time"

type Session struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time,omitempty"`
	IsActive  bool      `json:"is_active"`
}

func NewSession(id, title string) *Session {
	return &Session{
		ID:        id,
		Title:     title,
		StartTime: time.Now(),
		IsActive:  true,
	}
}
