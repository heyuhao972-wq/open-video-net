package model

type Comment struct {
	ID        int    `json:"id"`
	VideoID   string `json:"video_id"`
	UserID    string `json:"user_id"`
	Content   string `json:"content"`
	Likes     int    `json:"likes"`
	CreatedAt int64  `json:"created_at"`
}
