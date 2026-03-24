package db

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

func Open(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

func InitCommentTables(db *sql.DB) error {
	if db == nil {
		return fmt.Errorf("db is nil")
	}
	_, err := db.Exec(`
CREATE TABLE IF NOT EXISTS comments (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	video_id TEXT,
	user_id TEXT,
	content TEXT,
	parent_id INTEGER DEFAULT 0,
	likes INTEGER DEFAULT 0,
	status TEXT DEFAULT 'approved',
	review_reason TEXT,
	reviewed_by TEXT,
	reviewed_at INTEGER,
	created_at INTEGER
);
CREATE TABLE IF NOT EXISTS comment_likes (
	user_id TEXT,
	comment_id INTEGER,
	PRIMARY KEY (user_id, comment_id)
);
`)
	if err != nil {
		return err
	}
	// backfill older databases
	_, _ = db.Exec("ALTER TABLE comments ADD COLUMN parent_id INTEGER DEFAULT 0")
	_, _ = db.Exec("ALTER TABLE comments ADD COLUMN status TEXT DEFAULT 'approved'")
	_, _ = db.Exec("ALTER TABLE comments ADD COLUMN review_reason TEXT")
	_, _ = db.Exec("ALTER TABLE comments ADD COLUMN reviewed_by TEXT")
	_, _ = db.Exec("ALTER TABLE comments ADD COLUMN reviewed_at INTEGER")
	return nil
}
