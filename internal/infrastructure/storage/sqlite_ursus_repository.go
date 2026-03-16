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

func NewSQLiteUrsusRepository(db *sql.DB) *SQLiteUrsusRepository {
	return &SQLiteUrsusRepository{db: db}
}

func (r *SQLiteUrsusRepository) DB() *sql.DB {
	return r.db
}

// Migrate sets up the database schema with FTS5 support.
func (r *SQLiteUrsusRepository) Migrate(ctx context.Context) error {
	// 1. Ensure ursus_sessions exists
	query := `
	CREATE TABLE IF NOT EXISTS ursus_sessions (
		id TEXT PRIMARY KEY,
		title TEXT,
		start_time DATETIME,
		end_time DATETIME,
		is_active INTEGER DEFAULT 1
	);`
	if _, err := r.db.ExecContext(ctx, query); err != nil {
		return err
	}

	// 2. Safely add session_id to ursus_data if it doesn't exist
	// We check if the column exists by querying the table info
	rows, err := r.db.QueryContext(ctx, "PRAGMA table_info(ursus_data)")
	if err == nil {
		hasSessionID := false
		for rows.Next() {
			var cid int
			var name, dtype string
			var notnull, pk int
			var dflt interface{}
			if err := rows.Scan(&cid, &name, &dtype, &notnull, &dflt, &pk); err != nil {
				continue
			}
			if name == "session_id" {
				hasSessionID = true
				break
			}
		}
		rows.Close()

		if !hasSessionID {
			// Check if table exists before altering
			var name string
			err := r.db.QueryRowContext(ctx, "SELECT name FROM sqlite_master WHERE type='table' AND name='ursus_data'").Scan(&name)
			if err == nil {
				_, _ = r.db.ExecContext(ctx, "ALTER TABLE ursus_data ADD COLUMN session_id TEXT REFERENCES ursus_sessions(id)")
			}
		}
	}

	// 2.5 Ensure ursus_prompts exists
	query = `
	CREATE TABLE IF NOT EXISTS ursus_prompts (
		id TEXT PRIMARY KEY,
		input TEXT,
		session_id TEXT,
		created_at DATETIME,
		FOREIGN KEY(session_id) REFERENCES ursus_sessions(id)
	);`
	if _, err := r.db.ExecContext(ctx, query); err != nil {
		return err
	}

	// 3. Create ursus_data if not exists
	query = `
	CREATE TABLE IF NOT EXISTS ursus_data (
		id TEXT PRIMARY KEY,
		content TEXT,
		metadata TEXT,
		session_id TEXT,
		topic_key TEXT,
		prompt_id TEXT,
		scope TEXT DEFAULT 'project',
		duplicate_count INTEGER DEFAULT 1,
		revision_count INTEGER DEFAULT 0,
		created_at DATETIME,
		updated_at DATETIME,
		last_seen_at DATETIME,
		deleted_at DATETIME,
		FOREIGN KEY(session_id) REFERENCES ursus_sessions(id),
		FOREIGN KEY(prompt_id) REFERENCES ursus_prompts(id)
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
	_, err = r.db.ExecContext(ctx, query)
	if err != nil {
		return err
	}

	// Handle schema evolution: Add missing columns if they don't exist
	evolutionQueries := []string{
		"ALTER TABLE ursus_data ADD COLUMN topic_key TEXT",
		"ALTER TABLE ursus_data ADD COLUMN prompt_id TEXT",
		"ALTER TABLE ursus_data ADD COLUMN scope TEXT DEFAULT 'project'",
		"ALTER TABLE ursus_data ADD COLUMN duplicate_count INTEGER DEFAULT 1",
		"ALTER TABLE ursus_data ADD COLUMN revision_count INTEGER DEFAULT 0",
		"ALTER TABLE ursus_data ADD COLUMN last_seen_at DATETIME",
		"ALTER TABLE ursus_data ADD COLUMN deleted_at DATETIME",
	}

	for _, q := range evolutionQueries {
		// We ignore errors here because SQLite doesn't have "IF NOT EXISTS" for ADD COLUMN
		// and it will error if the column already exists.
		_, _ = r.db.ExecContext(ctx, q)
	}

	return nil
}

func (r *SQLiteUrsusRepository) Save(ctx context.Context, u *entity.Ursus) error {
	query := `INSERT INTO ursus_data (id, content, metadata, session_id, topic_key, prompt_id, scope, duplicate_count, revision_count, created_at, updated_at, last_seen_at, deleted_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	_, err := r.db.ExecContext(ctx, query, u.ID, u.Content, u.Metadata, u.SessionID, u.TopicKey, u.PromptID, u.Scope, u.DuplicateCount, u.RevisionCount, u.CreatedAt, u.UpdatedAt, u.LastSeenAt, u.DeletedAt)
	return err
}

func (r *SQLiteUrsusRepository) Search(ctx context.Context, query string) ([]*entity.Ursus, error) {
	sqlQuery := `
		SELECT d.id, d.content, d.metadata, d.session_id, d.topic_key, d.prompt_id, d.scope, d.duplicate_count, d.revision_count, d.created_at, d.updated_at, d.last_seen_at, d.deleted_at
		FROM ursus_data d
		JOIN ursus_fts f ON d.id = f.id
		WHERE ursus_fts MATCH ? AND d.deleted_at IS NULL
		ORDER BY rank`
	
	rows, err := r.db.QueryContext(ctx, sqlQuery, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []*entity.Ursus
	for rows.Next() {
		u, err := r.scanUrsus(rows)
		if err != nil {
			return nil, err
		}
		results = append(results, u)
	}
	return results, nil
}

func (r *SQLiteUrsusRepository) List(ctx context.Context) ([]*entity.Ursus, error) {
	query := `SELECT id, content, metadata, session_id, topic_key, prompt_id, scope, duplicate_count, revision_count, created_at, updated_at, last_seen_at, deleted_at FROM ursus_data WHERE deleted_at IS NULL ORDER BY created_at DESC`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []*entity.Ursus
	for rows.Next() {
		u, err := r.scanUrsus(rows)
		if err != nil {
			return nil, err
		}
		results = append(results, u)
	}
	return results, nil
}

func (r *SQLiteUrsusRepository) ListBySession(ctx context.Context, sessionID string) ([]*entity.Ursus, error) {
	query := `SELECT id, content, metadata, session_id, topic_key, prompt_id, scope, duplicate_count, revision_count, created_at, updated_at, last_seen_at, deleted_at FROM ursus_data WHERE session_id = ? AND deleted_at IS NULL ORDER BY created_at ASC`
	rows, err := r.db.QueryContext(ctx, query, sessionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []*entity.Ursus
	for rows.Next() {
		u, err := r.scanUrsus(rows)
		if err != nil {
			return nil, err
		}
		results = append(results, u)
	}
	return results, nil
}

func (r *SQLiteUrsusRepository) GetByID(ctx context.Context, id string) (*entity.Ursus, error) {
	query := `SELECT id, content, metadata, session_id, topic_key, prompt_id, scope, duplicate_count, revision_count, created_at, updated_at, last_seen_at, deleted_at FROM ursus_data WHERE id = ? AND deleted_at IS NULL`
	row := r.db.QueryRowContext(ctx, query, id)
	return r.scanUrsusRow(row)
}

func (r *SQLiteUrsusRepository) GetByTopicKey(ctx context.Context, topicKey string) (*entity.Ursus, error) {
	query := `SELECT id, content, metadata, session_id, topic_key, prompt_id, scope, duplicate_count, revision_count, created_at, updated_at, last_seen_at, deleted_at FROM ursus_data WHERE topic_key = ? AND deleted_at IS NULL`
	row := r.db.QueryRowContext(ctx, query, topicKey)
	return r.scanUrsusRow(row)
}

func (r *SQLiteUrsusRepository) scanUrsus(rows *sql.Rows) (*entity.Ursus, error) {
	u := &entity.Ursus{}
	var sessionID, topicKey, promptID *string
	var lastSeen, updated, deleted sql.NullTime
	err := rows.Scan(
		&u.ID, &u.Content, &u.Metadata, &sessionID, &topicKey, &promptID,
		&u.Scope, &u.DuplicateCount, &u.RevisionCount,
		&u.CreatedAt, &updated, &lastSeen, &deleted,
	)
	if err != nil {
		return nil, err
	}
	u.UpdatedAt = updated.Time
	if sessionID != nil {
		u.SessionID = *sessionID
	}
	if topicKey != nil {
		u.TopicKey = *topicKey
	}
	if promptID != nil {
		u.PromptID = *promptID
	}
	if lastSeen.Valid {
		u.LastSeenAt = lastSeen.Time
	}
	if deleted.Valid {
		u.DeletedAt = &deleted.Time
	}
	return u, nil
}

func (r *SQLiteUrsusRepository) scanUrsusRow(row *sql.Row) (*entity.Ursus, error) {
	u := &entity.Ursus{}
	var sessionID, topicKey, promptID *string
	var lastSeen, updated, deleted sql.NullTime
	err := row.Scan(
		&u.ID, &u.Content, &u.Metadata, &sessionID, &topicKey, &promptID,
		&u.Scope, &u.DuplicateCount, &u.RevisionCount,
		&u.CreatedAt, &updated, &lastSeen, &deleted,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	u.UpdatedAt = updated.Time
	if sessionID != nil {
		u.SessionID = *sessionID
	}
	if topicKey != nil {
		u.TopicKey = *topicKey
	}
	if promptID != nil {
		u.PromptID = *promptID
	}
	if lastSeen.Valid {
		u.LastSeenAt = lastSeen.Time
	}
	if deleted.Valid {
		u.DeletedAt = &deleted.Time
	}
	return u, nil
}

func (r *SQLiteUrsusRepository) Update(ctx context.Context, u *entity.Ursus) error {
	query := `UPDATE ursus_data SET content = ?, metadata = ?, session_id = ?, topic_key = ?, prompt_id = ?, scope = ?, duplicate_count = ?, revision_count = ?, updated_at = ?, last_seen_at = ?, deleted_at = ? WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, u.Content, u.Metadata, u.SessionID, u.TopicKey, u.PromptID, u.Scope, u.DuplicateCount, u.RevisionCount, u.UpdatedAt, u.LastSeenAt, u.DeletedAt, u.ID)
	return err
}

func (r *SQLiteUrsusRepository) Delete(ctx context.Context, id string) error {
	query := `UPDATE ursus_data SET deleted_at = CURRENT_TIMESTAMP WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
