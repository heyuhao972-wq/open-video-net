package model

type Video struct {
	ID              string   `json:"id"`
	PlatformID      string   `json:"platform_id"`
	StorageID       string   `json:"storage_id"`
	Title           string   `json:"title"`
	Description     string   `json:"description"`
	Filename        string   `json:"filename"`
	FilePath        string   `json:"file_path"`
	Tags            []string `json:"tags"`
	AuthorID        string   `json:"author_id"`
	AuthorPublicKey string   `json:"author_public_key"`
	AuthorSignature string   `json:"author_signature"`
	ProofTimestamp  int64    `json:"proof_timestamp"`
	VideoHash       string   `json:"video_hash"`
	Chunks          []string `json:"chunks"`
	Manifest        string   `json:"manifest"`
	ManifestHash    string   `json:"manifest_hash"`
	CreatedAt       int64    `json:"created_at"`
}
