package storage

import (
	"testing"

	"video-storage/internal/chunk"
)

func TestLocalStoreSaveGet(t *testing.T) {
	store, err := NewLocalStore(t.TempDir())
	if err != nil {
		t.Fatalf("new local store: %v", err)
	}

	c := chunk.Chunk{
		Hash: "abc",
		Data: []byte("hello"),
	}

	if err := store.Save(c); err != nil {
		t.Fatalf("save: %v", err)
	}

	data, err := store.Get("abc")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if string(data) != "hello" {
		t.Fatalf("unexpected data: %s", string(data))
	}
}
