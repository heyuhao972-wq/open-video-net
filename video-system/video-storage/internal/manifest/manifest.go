package manifest

import (
	"encoding/json"
	"os"
)

type ChunkRef struct {
	Hash  string `json:"hash"`
	Index int    `json:"index"`
	Size  int    `json:"size"`
}

type VideoManifest struct {
	VideoID         string     `json:"video_id"`
	VideoHash       string     `json:"video_hash"`
	Timestamp       int64      `json:"timestamp"`
	AuthorPublicKey string     `json:"author_public_key"`
	Signature       string     `json:"signature"`
	Chunks          []ChunkRef `json:"chunks"`
}

func (m *VideoManifest) Save(path string) error {

	data, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}
