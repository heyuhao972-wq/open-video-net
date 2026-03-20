package service

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"video-platform/db"
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

	dbPath := filepath.Join(tmp, "platform.db")
	database, err := db.Open(dbPath)
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	defer database.Close()
	repo := repository.NewVideoRepository(database)
	videoService := NewVideoService(repo)
	uploadService := NewUploadService(videoService, storeClient, indexClient)

	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatalf("generate key: %v", err)
	}
	pubB64 := base64.StdEncoding.EncodeToString(pub)
	ts := time.Now().Unix()
	hash := sha256.Sum256(data)
	videoHash := hex.EncodeToString(hash[:])
	msg := []byte(videoHash + "|" + fmt.Sprintf("%d", ts) + "|" + pubB64)
	sig := base64.StdEncoding.EncodeToString(ed25519.Sign(priv, msg))

	video, err := uploadService.UploadVideo("demo", "desc", []string{"tag1"}, videoPath, "video.bin", "author-1", pubB64, sig, videoHash, ts, "platformA")
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
