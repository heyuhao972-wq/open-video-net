package db

import (
	"database/sql"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

func Open(path string) (*sql.DB, error) {
	if path == "" {
		path = "./data/platform.db"
	}
	if err := os.MkdirAll(filepath.Dir(path), os.ModePerm); err != nil {
		return nil, err
	}
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}
	if err := migrate(db); err != nil {
		return nil, err
	}
	return db, nil
}

func migrate(db *sql.DB) error {
	_, err := db.Exec(`
CREATE TABLE IF NOT EXISTS users (
  user_id TEXT PRIMARY KEY,
  public_key TEXT NOT NULL,
  created_at INTEGER NOT NULL
);
CREATE UNIQUE INDEX IF NOT EXISTS idx_users_public_key ON users(public_key);

CREATE TABLE IF NOT EXISTS videos (
  video_id TEXT PRIMARY KEY,
  platform_id TEXT,
  storage_id TEXT,
  title TEXT,
  description TEXT,
  filename TEXT,
  file_path TEXT,
  tags TEXT,
  author_id TEXT,
  author_public_key TEXT,
  author_signature TEXT,
  proof_timestamp INTEGER,
  video_hash TEXT,
  chunks TEXT,
  manifest TEXT,
  manifest_hash TEXT,
  created_at INTEGER
);
`)
	return err
}
