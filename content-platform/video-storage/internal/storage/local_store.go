package storage

import (
	"os"
	"path/filepath"

	"video-storage/internal/chunk"
)

type LocalStore struct {
	basePath string
}

func NewLocalStore(path string) (*LocalStore, error) {
	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		return nil, err
	}

	return &LocalStore{
		basePath: path,
	}, nil
}

func (s *LocalStore) Save(c chunk.Chunk) error {

	path := filepath.Join(s.basePath, c.Hash)

	return os.WriteFile(path, c.Data, 0644)
}

func (s *LocalStore) Get(hash string) ([]byte, error) {

	path := filepath.Join(s.basePath, hash)

	return os.ReadFile(path)
}
