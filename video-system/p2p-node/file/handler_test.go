package file

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestReceiveFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "recv.bin")

	input := []byte("peer data")
	if err := ReceiveFile(bytes.NewReader(input), path); err != nil {
		t.Fatalf("receive file: %v", err)
	}

	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read received file: %v", err)
	}
	if string(got) != string(input) {
		t.Fatalf("expected %q, got %q", string(input), string(got))
	}
}
