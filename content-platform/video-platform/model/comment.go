package model

type Comment struct {
	ID           int    `json:"id"`
	VideoID      string `json:"video_id"`
	UserID       string `json:"user_id"`
	Content      string `json:"content"`
	ParentID     int    `json:"parent_id"`
	Likes        int    `json:"likes"`
	CreatedAt    int64  `json:"created_at"`
	Status       string `json:"status"`
	ReviewReason string `json:"review_reason"`
	ReviewedBy   string `json:"reviewed_by"`
	ReviewedAt   int64  `json:"reviewed_at"`
}
