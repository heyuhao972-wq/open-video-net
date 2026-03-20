package model

type Video struct {
	ID         string   `json:"id"`
	PlatformID string   `json:"platform_id"`
	Title      string   `json:"title"`
	Views      int      `json:"views"`
	CreatedAt  int64    `json:"created_at"`
	Tags       []string `json:"tags"`
	AuthorID   string   `json:"author_id"`
}
