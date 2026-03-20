package db

import (
	"database/sql"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

func Open(path string) (*sql.DB, error) {
	if path == "" {
		path = "./data/recommendation.db"
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
CREATE TABLE IF NOT EXISTS behaviors (
  user_id TEXT,
  video_id TEXT,
  type TEXT,
  timestamp INTEGER
);

CREATE TABLE IF NOT EXISTS follows (
  follower_id TEXT,
  followee_id TEXT,
  active INTEGER,
  timestamp INTEGER,
  PRIMARY KEY (follower_id, followee_id)
);
`)
	return err
}
