package storage

import "video-storage/internal/chunk"

type ChunkStore interface {
	Save(chunk chunk.Chunk) error

	Get(hash string) ([]byte, error)
}
