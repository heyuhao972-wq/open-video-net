package model

type Follow struct {
	UserID    string `json:"user_id"`
	AuthorID  string `json:"author_id"`
	Active    bool   `json:"active"`
	Timestamp int64  `json:"timestamp"`
}
