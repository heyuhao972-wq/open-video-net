package chunk

import (
	"errors"
	"os"
	"path/filepath"
)

type Store struct {
	dir string
}

func NewStore(dir string) (*Store, error) {
	if dir == "" {
		return nil, errors.New("chunk dir required")
	}
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return nil, err
	}
	return &Store{dir: dir}, nil
}

func (s *Store) Get(hash string) ([]byte, error) {
	path := filepath.Join(s.dir, hash)
	return os.ReadFile(path)
}

func (s *Store) Put(hash string, data []byte) error {
	path := filepath.Join(s.dir, hash)
	return os.WriteFile(path, data, 0644)
}
