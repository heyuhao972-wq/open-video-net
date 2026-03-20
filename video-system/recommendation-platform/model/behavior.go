package model

type Behavior struct {
	UserID    string `json:"user_id"`
	VideoID   string `json:"video_id"`
	Type      string `json:"type"`
	Timestamp int64  `json:"timestamp"`
}
