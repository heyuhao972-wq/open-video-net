package service

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"video-platform/index"
	"video-platform/repository"
	"video-platform/storage"
)

func TestUploadServiceStoresChunksAndManifest(t *testing.T) {
	tmp := t.TempDir()
	videoPath := filepath.Join(tmp, "video.bin")

	data := make([]byte, 2500)
	for i := range data {
		data[i] = byte(i % 251)
	}
	if err := os.WriteFile(videoPath, data, 0644); err != nil {
		t.Fatalf("write temp video: %v", err)
	}

	storeClient, err := storage.NewStorageClient(filepath.Join(tmp, "store"), 1024)
	if err != nil {
		t.Fatalf("new storage client: %v", err)
	}

	indexServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/video" {
			http.NotFound(w, r)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer indexServer.Close()
	indexClient := index.NewClient(indexServer.URL)

	repo := repository.NewVideoRepository()
	videoService := NewVideoService(repo)
	uploadService := NewUploadService(videoService, storeClient, indexClient)

	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatalf("generate key: %v", err)
	}
	pubB64 := base64.StdEncoding.EncodeToString(pub)

	hash := sha256.Sum256(data)
	videoHash := hex.EncodeToString(hash[:])
	timestamp := int64(1234567890)
	proofMsg := videoHash + "|" + "1234567890" + "|" + pubB64
	signature := base64.StdEncoding.EncodeToString(ed25519.Sign(priv, []byte(proofMsg)))

	video, err := uploadService.UploadVideo("demo", "desc", []string{"tag1"}, videoPath, "video.bin", "", "author-1", pubB64, signature, timestamp, videoHash, "platformA")
	if err != nil {
		t.Fatalf("upload video: %v", err)
	}

	if video.StorageID == "" {
		t.Fatalf("expected non-empty storage id")
	}
	if len(video.Chunks) != 3 {
		t.Fatalf("expected 3 chunks, got %d", len(video.Chunks))
	}
	if video.Manifest == "" {
		t.Fatalf("expected manifest path")
	}
	if _, err := os.Stat(video.Manifest); err != nil {
		t.Fatalf("manifest path invalid: %v", err)
	}
	if video.ManifestHash == "" || video.AuthorPublicKey == "" {
		t.Fatalf("expected manifest hash and author public key")
	}
}
