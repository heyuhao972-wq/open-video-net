package model

type Report struct {
	ID        int    `json:"id"`
	TargetID  string `json:"target_id"`
	TargetType string `json:"target_type"`
	UserID    string `json:"user_id"`
	Reason    string `json:"reason"`
	CreatedAt int64  `json:"created_at"`
}
