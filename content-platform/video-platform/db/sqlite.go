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
	likes INTEGER DEFAULT 0,
	created_at INTEGER
);
CREATE TABLE IF NOT EXISTS comment_likes (
	user_id TEXT,
	comment_id INTEGER,
	PRIMARY KEY (user_id, comment_id)
);
`)
	return err
}
