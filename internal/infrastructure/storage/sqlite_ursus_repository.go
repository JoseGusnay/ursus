package storage

import (
	"context"
	"database/sql"

	"github.com/JoseGusnay/ursus/internal/domain/entity"
	_ "modernc.org/sqlite"
)

// SQLiteUrsusRepository implements the UrsusRepository interface using SQLite and FTS5.
type SQLiteUrsusRepository struct {
	db *sql.DB
}

// NewSQLiteUrsusRepository creates a new instance of SQLiteUrsusRepository.
func NewSQLiteUrsusRepository(db *sql.DB) *SQLiteUrsusRepository {
	return &SQLiteUrsusRepository{db: db}
}

// Migrate sets up the database schema with FTS5 support.
func (r *SQLiteUrsusRepository) Migrate(ctx context.Context) error {
	query := `
	CREATE TABLE IF NOT EXISTS ursus_data (
		id TEXT PRIMARY KEY,
		content TEXT,
		metadata TEXT,
		created_at DATETIME,
		updated_at DATETIME
	);
	CREATE VIRTUAL TABLE IF NOT EXISTS ursus_fts USING fts5(
		id UNINDEXED,
		content,
		metadata
	);
	
	-- Triggers to keep FTS index in sync
	CREATE TRIGGER IF NOT EXISTS ursus_ai AFTER INSERT ON ursus_data BEGIN
		INSERT INTO ursus_fts(id, content, metadata) VALUES (new.id, new.content, new.metadata);
	END;
	CREATE TRIGGER IF NOT EXISTS ursus_ad AFTER DELETE ON ursus_data BEGIN
		DELETE FROM ursus_fts WHERE id = old.id;
	END;
	CREATE TRIGGER IF NOT EXISTS ursus_au AFTER UPDATE ON ursus_data BEGIN
		UPDATE ursus_fts SET content = new.content, metadata = new.metadata WHERE id = new.id;
	END;
	`
	_, err := r.db.ExecContext(ctx, query)
	return err
}

func (r *SQLiteUrsusRepository) Save(ctx context.Context, u *entity.Ursus) error {
	query := `INSERT OR REPLACE INTO ursus_data (id, content, metadata, created_at, updated_at) VALUES (?, ?, ?, ?, ?)`
	_, err := r.db.ExecContext(ctx, query, u.ID, u.Content, u.Metadata, u.CreatedAt, u.UpdatedAt)
	return err
}

func (r *SQLiteUrsusRepository) Search(ctx context.Context, query string) ([]*entity.Ursus, error) {
	sqlQuery := `
		SELECT d.id, d.content, d.metadata, d.created_at, d.updated_at 
		FROM ursus_data d
		JOIN ursus_fts f ON d.id = f.id
		WHERE ursus_fts MATCH ?
		ORDER BY rank`
	
	rows, err := r.db.QueryContext(ctx, sqlQuery, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []*entity.Ursus
	for rows.Next() {
		u := &entity.Ursus{}
		if err := rows.Scan(&u.ID, &u.Content, &u.Metadata, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, err
		}
		results = append(results, u)
	}
	return results, nil
}

func (r *SQLiteUrsusRepository) List(ctx context.Context) ([]*entity.Ursus, error) {
	query := `SELECT id, content, metadata, created_at, updated_at FROM ursus_data ORDER BY created_at DESC`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []*entity.Ursus
	for rows.Next() {
		u := &entity.Ursus{}
		if err := rows.Scan(&u.ID, &u.Content, &u.Metadata, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, err
		}
		results = append(results, u)
	}
	return results, nil
}

func (r *SQLiteUrsusRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM ursus_data WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
