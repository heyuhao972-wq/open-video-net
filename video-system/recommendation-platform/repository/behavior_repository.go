package repository

import (
	"database/sql"
	"time"

	"recommendation-platform/model"
)

type Store struct {
	db *sql.DB
}

var store *Store

func Init(db *sql.DB) {
	store = &Store{db: db}
}

func AddBehavior(b model.Behavior) {
	if store == nil {
		return
	}
	if b.Timestamp == 0 {
		b.Timestamp = time.Now().Unix()
	}
	_, _ = store.db.Exec(
		`INSERT INTO behaviors (user_id, video_id, type, timestamp) VALUES (?, ?, ?, ?)`,
		b.UserID, b.VideoID, b.Type, b.Timestamp,
	)
}

func GetBehaviors() []model.Behavior {
	if store == nil {
		return []model.Behavior{}
	}
	rows, err := store.db.Query(`SELECT user_id, video_id, type, timestamp FROM behaviors`)
	if err != nil {
		return []model.Behavior{}
	}
	defer rows.Close()
	out := []model.Behavior{}
	for rows.Next() {
		var b model.Behavior
		_ = rows.Scan(&b.UserID, &b.VideoID, &b.Type, &b.Timestamp)
		out = append(out, b)
	}
	return out
}

func AddFollow(f model.Follow) {
	if store == nil {
		return
	}
	if !f.Active {
		f.Active = true
	}
	if f.Timestamp == 0 {
		f.Timestamp = time.Now().Unix()
	}
	_, _ = store.db.Exec(
		`INSERT INTO follows (follower_id, followee_id, active, timestamp)
		 VALUES (?, ?, ?, ?)
		 ON CONFLICT(follower_id, followee_id) DO UPDATE SET active=excluded.active, timestamp=excluded.timestamp`,
		f.UserID, f.AuthorID, f.Active, f.Timestamp,
	)
}

func GetFollows() []model.Follow {
	if store == nil {
		return []model.Follow{}
	}
	rows, err := store.db.Query(`SELECT follower_id, followee_id, active, timestamp FROM follows`)
	if err != nil {
		return []model.Follow{}
	}
	defer rows.Close()
	out := []model.Follow{}
	for rows.Next() {
		var f model.Follow
		_ = rows.Scan(&f.UserID, &f.AuthorID, &f.Active, &f.Timestamp)
		out = append(out, f)
	}
	return out
}

func RemoveFollow(userID string, authorID string) {
	if store == nil {
		return
	}
	_, _ = store.db.Exec(
		`UPDATE follows SET active=0, timestamp=? WHERE follower_id=? AND followee_id=?`,
		time.Now().Unix(), userID, authorID,
	)
}
