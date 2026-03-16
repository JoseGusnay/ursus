package storage

import (
	"context"
	"database/sql"
	"github.com/JoseGusnay/ursus/internal/domain/entity"
)

type SQLiteSessionRepository struct {
	db *sql.DB
}

func NewSQLiteSessionRepository(db *sql.DB) *SQLiteSessionRepository {
	return &SQLiteSessionRepository{db: db}
}

func (r *SQLiteSessionRepository) Save(ctx context.Context, s *entity.Session) error {
	query := `INSERT OR REPLACE INTO ursus_sessions (id, title, start_time, end_time, is_active) VALUES (?, ?, ?, ?, ?)`
	isActive := 0
	if s.IsActive {
		isActive = 1
	}
	_, err := r.db.ExecContext(ctx, query, s.ID, s.Title, s.StartTime, s.EndTime, isActive)
	return err
}

func (r *SQLiteSessionRepository) GetActive(ctx context.Context) (*entity.Session, error) {
	query := `SELECT id, title, start_time, end_time, is_active FROM ursus_sessions WHERE is_active = 1 LIMIT 1`
	row := r.db.QueryRowContext(ctx, query)
	
	s := &entity.Session{}
	var isActive int
	err := row.Scan(&s.ID, &s.Title, &s.StartTime, &s.EndTime, &isActive)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	s.IsActive = isActive == 1
	return s, nil
}

func (r *SQLiteSessionRepository) GetByID(ctx context.Context, id string) (*entity.Session, error) {
	query := `SELECT id, title, start_time, end_time, is_active FROM ursus_sessions WHERE id = ?`
	row := r.db.QueryRowContext(ctx, query, id)
	
	s := &entity.Session{}
	var isActive int
	err := row.Scan(&s.ID, &s.Title, &s.StartTime, &s.EndTime, &isActive)
	if err != nil {
		return nil, err
	}
	s.IsActive = isActive == 1
	return s, nil
}

func (r *SQLiteSessionRepository) List(ctx context.Context) ([]*entity.Session, error) {
	query := `SELECT id, title, start_time, end_time, is_active FROM ursus_sessions ORDER BY start_time DESC`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []*entity.Session
	for rows.Next() {
		s := &entity.Session{}
		var isActive int
		if err := rows.Scan(&s.ID, &s.Title, &s.StartTime, &s.EndTime, &isActive); err != nil {
			return nil, err
		}
		s.IsActive = isActive == 1
		results = append(results, s)
	}
	return results, nil
}

func (r *SQLiteSessionRepository) DeactivateAll(ctx context.Context) error {
	query := `UPDATE ursus_sessions SET is_active = 0 WHERE is_active = 1`
	_, err := r.db.ExecContext(ctx, query)
	return err
}
