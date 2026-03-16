package storage

import (
	"context"
	"database/sql"

	"github.com/JoseGusnay/ursus/internal/domain/entity"
)

// SQLitePromptRepository implements the PromptRepository interface using SQLite.
type SQLitePromptRepository struct {
	db *sql.DB
}

// NewSQLitePromptRepository creates a new SQLitePromptRepository instance.
func NewSQLitePromptRepository(db *sql.DB) *SQLitePromptRepository {
	return &SQLitePromptRepository{db: db}
}

// Save stores a new prompt in the database.
func (r *SQLitePromptRepository) Save(ctx context.Context, p *entity.Prompt) error {
	query := `INSERT INTO ursus_prompts (id, input, session_id, created_at) VALUES (?, ?, ?, ?)`
	var sessionID interface{}
	if p.SessionID != "" {
		sessionID = p.SessionID
	}
	_, err := r.db.ExecContext(ctx, query, p.ID, p.Input, sessionID, p.CreatedAt)
	return err
}

// GetByID retrieves a prompt by its ID.
func (r *SQLitePromptRepository) GetByID(ctx context.Context, id string) (*entity.Prompt, error) {
	query := `SELECT id, input, session_id, created_at FROM ursus_prompts WHERE id = ?`
	p := &entity.Prompt{}
	var sessionID *string
	err := r.db.QueryRowContext(ctx, query, id).Scan(&p.ID, &p.Input, &sessionID, &p.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	if sessionID != nil {
		p.SessionID = *sessionID
	}
	return p, nil
}
