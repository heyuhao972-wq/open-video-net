package manifest

import (
	"testing"

	"video-storage/internal/chunk"
)

func TestBuildManifest(t *testing.T) {
	chunks := []chunk.Chunk{
		{Hash: "h1", Index: 0, Size: 10},
		{Hash: "h2", Index: 1, Size: 20},
	}

	m, err := BuildManifest(chunks, "hash", 123)
	if err != nil {
		t.Fatalf("build manifest: %v", err)
	}
	if m.VideoID == "" {
		t.Fatalf("expected non-empty video id")
	}
	if len(m.Chunks) != 2 {
		t.Fatalf("expected 2 chunks, got %d", len(m.Chunks))
	}
}
