package repository

import (
	"database/sql"
	"errors"
	"time"

	"video-platform/model"
)

type CommentRepository struct {
	db *sql.DB
}

func NewCommentRepository(db *sql.DB) *CommentRepository {
	return &CommentRepository{db: db}
}

func (r *CommentRepository) Add(videoID string, userID string, content string, parentID int, status string) (model.Comment, error) {
	if r.db == nil {
		return model.Comment{}, errors.New("db not initialized")
	}
	if status == "" {
		status = "approved"
	}
	res, err := r.db.Exec(
		"INSERT INTO comments (video_id, user_id, content, parent_id, likes, status, created_at) VALUES (?, ?, ?, ?, 0, ?, strftime('%s','now'))",
		videoID, userID, content, parentID, status,
	)
	if err != nil {
		return model.Comment{}, err
	}
	id, _ := res.LastInsertId()
	return model.Comment{
		ID:        int(id),
		VideoID:   videoID,
		UserID:    userID,
		Content:   content,
		ParentID:  parentID,
		Likes:     0,
		Status:    status,
		CreatedAt: time.Now().Unix(),
	}, nil
}

func (r *CommentRepository) Get(id int) (model.Comment, bool) {
	if r.db == nil {
		return model.Comment{}, false
	}
	row := r.db.QueryRow("SELECT id, video_id, user_id, content, parent_id, likes, created_at, status, COALESCE(review_reason,''), COALESCE(reviewed_by,''), COALESCE(reviewed_at,0) FROM comments WHERE id = ?", id)
	var c model.Comment
	if err := row.Scan(&c.ID, &c.VideoID, &c.UserID, &c.Content, &c.ParentID, &c.Likes, &c.CreatedAt, &c.Status, &c.ReviewReason, &c.ReviewedBy, &c.ReviewedAt); err != nil {
		return model.Comment{}, false
	}
	return c, true
}

func (r *CommentRepository) ListByVideo(videoID string) []model.Comment {
	if r.db == nil {
		return nil
	}
	rows, err := r.db.Query("SELECT id, video_id, user_id, content, parent_id, likes, created_at, status, COALESCE(review_reason,''), COALESCE(reviewed_by,''), COALESCE(reviewed_at,0) FROM comments WHERE video_id = ? AND status = 'approved' ORDER BY created_at ASC", videoID)
	if err != nil {
		return nil
	}
	defer rows.Close()

	out := []model.Comment{}
	for rows.Next() {
		var c model.Comment
		if err := rows.Scan(&c.ID, &c.VideoID, &c.UserID, &c.Content, &c.ParentID, &c.Likes, &c.CreatedAt, &c.Status, &c.ReviewReason, &c.ReviewedBy, &c.ReviewedAt); err == nil {
			out = append(out, c)
		}
	}
	return out
}

func (r *CommentRepository) ListByVideoAll(videoID string) []model.Comment {
	if r.db == nil {
		return nil
	}
	rows, err := r.db.Query("SELECT id, video_id, user_id, content, parent_id, likes, created_at, status, COALESCE(review_reason,''), COALESCE(reviewed_by,''), COALESCE(reviewed_at,0) FROM comments WHERE video_id = ? ORDER BY created_at ASC", videoID)
	if err != nil {
		return nil
	}
	defer rows.Close()

	out := []model.Comment{}
	for rows.Next() {
		var c model.Comment
		if err := rows.Scan(&c.ID, &c.VideoID, &c.UserID, &c.Content, &c.ParentID, &c.Likes, &c.CreatedAt, &c.Status, &c.ReviewReason, &c.ReviewedBy, &c.ReviewedAt); err == nil {
			out = append(out, c)
		}
	}
	return out
}

func (r *CommentRepository) CountByVideo(videoID string) int {
	if r.db == nil {
		return 0
	}
	row := r.db.QueryRow("SELECT COUNT(1) FROM comments WHERE video_id = ? AND status = 'approved'", videoID)
	var count int
	if err := row.Scan(&count); err != nil {
		return 0
	}
	return count
}

func (r *CommentRepository) CountByVideoAll(videoID string) int {
	if r.db == nil {
		return 0
	}
	row := r.db.QueryRow("SELECT COUNT(1) FROM comments WHERE video_id = ?", videoID)
	var count int
	if err := row.Scan(&count); err != nil {
		return 0
	}
	return count
}

func (r *CommentRepository) Delete(id int) {
	if r.db == nil {
		return
	}
	_, _ = r.db.Exec("DELETE FROM comment_likes WHERE comment_id = ?", id)
	_, _ = r.db.Exec("DELETE FROM comments WHERE parent_id = ?", id)
	_, _ = r.db.Exec("DELETE FROM comments WHERE id = ?", id)
}

func (r *CommentRepository) Like(id int, userID string) (model.Comment, bool, bool) {
	if r.db == nil {
		return model.Comment{}, false, false
	}

	tx, err := r.db.Begin()
	if err != nil {
		return model.Comment{}, false, false
	}
	defer tx.Rollback()

	var c model.Comment
	row := tx.QueryRow("SELECT id, video_id, user_id, content, parent_id, likes, created_at, status, COALESCE(review_reason,''), COALESCE(reviewed_by,''), COALESCE(reviewed_at,0) FROM comments WHERE id = ?", id)
	if err := row.Scan(&c.ID, &c.VideoID, &c.UserID, &c.Content, &c.ParentID, &c.Likes, &c.CreatedAt, &c.Status, &c.ReviewReason, &c.ReviewedBy, &c.ReviewedAt); err != nil {
		return model.Comment{}, false, false
	}

	res, err := tx.Exec("INSERT OR IGNORE INTO comment_likes (user_id, comment_id) VALUES (?, ?)", userID, id)
	if err != nil {
		return model.Comment{}, true, false
	}
	affected, _ := res.RowsAffected()
	if affected == 0 {
		_ = tx.Commit()
		return c, true, false
	}

	if _, err := tx.Exec("UPDATE comments SET likes = likes + 1 WHERE id = ?", id); err != nil {
		return model.Comment{}, true, false
	}
	c.Likes++
	if err := tx.Commit(); err != nil {
		return model.Comment{}, true, false
	}
	return c, true, true
}

func (r *CommentRepository) ListByStatus(status string) []model.Comment {
	if r.db == nil {
		return nil
	}
	if status == "" {
		status = "pending"
	}
	rows, err := r.db.Query("SELECT id, video_id, user_id, content, parent_id, likes, created_at, status, COALESCE(review_reason,''), COALESCE(reviewed_by,''), COALESCE(reviewed_at,0) FROM comments WHERE status = ? ORDER BY created_at ASC", status)
	if err != nil {
		return nil
	}
	defer rows.Close()
	out := []model.Comment{}
	for rows.Next() {
		var c model.Comment
		if err := rows.Scan(&c.ID, &c.VideoID, &c.UserID, &c.Content, &c.ParentID, &c.Likes, &c.CreatedAt, &c.Status, &c.ReviewReason, &c.ReviewedBy, &c.ReviewedAt); err == nil {
			out = append(out, c)
		}
	}
	return out
}

func (r *CommentRepository) UpdateReview(id int, status string, reason string, reviewer string, reviewedAt int64) (model.Comment, bool) {
	if r.db == nil {
		return model.Comment{}, false
	}
	_, err := r.db.Exec("UPDATE comments SET status = ?, review_reason = ?, reviewed_by = ?, reviewed_at = ? WHERE id = ?", status, reason, reviewer, reviewedAt, id)
	if err != nil {
		return model.Comment{}, false
	}
	return r.Get(id)
}
