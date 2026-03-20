package chunk

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSplitFile(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "video.bin")

	data := make([]byte, 2500)
	for i := range data {
		data[i] = byte(i % 251)
	}
	if err := os.WriteFile(p, data, 0644); err != nil {
		t.Fatalf("write temp file: %v", err)
	}

	chunks, err := SplitFile(p, 1024)
	if err != nil {
		t.Fatalf("split file: %v", err)
	}
	if len(chunks) != 3 {
		t.Fatalf("expected 3 chunks, got %d", len(chunks))
	}
	if chunks[0].Index != 0 || chunks[1].Index != 1 || chunks[2].Index != 2 {
		t.Fatalf("unexpected chunk indexes: %+v", chunks)
	}
}
