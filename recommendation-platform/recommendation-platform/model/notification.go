package model

type Notification struct {
	ID        int    `json:"id"`
	UserID    string `json:"user_id"`
	ActorID   string `json:"actor_id"`
	VideoID   string `json:"video_id"`
	Type      string `json:"type"`
	Read      bool   `json:"read"`
	CreatedAt int64  `json:"created_at"`
}
