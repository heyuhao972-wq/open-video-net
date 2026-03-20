package videostorage

import (
	"os"
	"path/filepath"
	"testing"
)

func TestProcessorStoreVideo(t *testing.T) {
	tmp := t.TempDir()
	videoPath := filepath.Join(tmp, "video.bin")

	data := make([]byte, 2100)
	for i := range data {
		data[i] = byte(i % 251)
	}
	if err := os.WriteFile(videoPath, data, 0644); err != nil {
		t.Fatalf("write video file: %v", err)
	}

	p, err := NewProcessor(filepath.Join(tmp, "storage"), 1024)
	if err != nil {
		t.Fatalf("new processor: %v", err)
	}

	result, err := p.StoreVideo(videoPath)
	if err != nil {
		t.Fatalf("store video: %v", err)
	}

	if result.VideoID == "" {
		t.Fatalf("expected non-empty video id")
	}
	if len(result.ChunkHashes) != 3 {
		t.Fatalf("expected 3 chunk hashes, got %d", len(result.ChunkHashes))
	}
	if _, err := os.Stat(result.ManifestPath); err != nil {
		t.Fatalf("manifest not found: %v", err)
	}
}
