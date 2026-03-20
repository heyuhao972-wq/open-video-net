package manifest

import (
	"crypto/rand"
	"encoding/hex"
	"errors"

	"video-storage/internal/chunk"
)

func randomID() string {

	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return ""
	}

	return hex.EncodeToString(b)
}

func BuildManifest(chunks []chunk.Chunk, videoHash string, timestamp int64) (VideoManifest, error) {

	var refs []ChunkRef

	for _, c := range chunks {

		ref := ChunkRef{
			Hash:  c.Hash,
			Index: c.Index,
			Size:  c.Size,
		}

		refs = append(refs, ref)
	}

	id := randomID()
	if id == "" {
		return VideoManifest{}, errors.New("failed to generate video id")
	}

	return VideoManifest{
		VideoID:   id,
		VideoHash: videoHash,
		Timestamp: timestamp,
		Chunks:    refs,
	}, nil
}
