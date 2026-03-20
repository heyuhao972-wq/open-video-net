package model

type Video struct {
	ID              string   `json:"id"`
	PlatformID      string   `json:"platform_id"`
	Title           string   `json:"title"`
	Description     string   `json:"description"`
	Tags            []string `json:"tags"`
	Views           int      `json:"views"`
	CreatedAt       int64    `json:"created_at"`
	AuthorID        string   `json:"author_id"`
	AuthorPublicKey string   `json:"author_public_key"`
	ManifestHash    string   `json:"manifest_hash"`
}
