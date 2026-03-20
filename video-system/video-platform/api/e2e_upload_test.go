package api

import (
	"bytes"
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/gin-gonic/gin"

	"video-platform/db"
	"video-platform/handler"
	"video-platform/index"
	"video-platform/repository"
	"video-platform/service"
	"video-platform/storage"
)

type uploadResp struct {
	Video struct {
		ID              string   `json:"id"`
		StorageID       string   `json:"storage_id"`
		Chunks          []string `json:"chunks"`
		Manifest        string   `json:"manifest"`
		AuthorPublicKey string   `json:"author_public_key"`
		ManifestHash    string   `json:"manifest_hash"`
	} `json:"video"`
}

func TestUploadAndGetVideoE2E(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tmp := t.TempDir()
	storeClient, err := storage.NewStorageClient(filepath.Join(tmp, "storage"), 1024)
	if err != nil {
		t.Fatalf("new storage client: %v", err)
	}

	indexServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost && (r.URL.Path == "/video" || r.URL.Path == "/video/views") {
			w.WriteHeader(http.StatusOK)
			return
		}
		http.NotFound(w, r)
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
	userRepo := repository.NewUserRepository(database)
	videoService := service.NewVideoService(repo)
	userService := service.NewUserService(userRepo, "dev-secret")
	uploadService := service.NewUploadService(videoService, storeClient, indexClient)
	handler.InitServices(videoService, uploadService, userService, storeClient)

	pub, priv, err := ed25519.GenerateKey(nil)
	if err != nil {
		t.Fatalf("generate key: %v", err)
	}
	pubB64 := base64.StdEncoding.EncodeToString(pub)
	_, err = userService.Register(pubB64)
	if err != nil {
		t.Fatalf("register: %v", err)
	}
	nonce, err := userService.CreateChallenge(pubB64)
	if err != nil {
		t.Fatalf("challenge: %v", err)
	}
	sig := base64.StdEncoding.EncodeToString(ed25519.Sign(priv, mustDecodeB64(nonce)))
	token, _, err := userService.Login(pubB64, nonce, sig)
	if err != nil {
		t.Fatalf("login: %v", err)
	}

	r := gin.New()
	RegisterRoutes(r)

	videoPath := filepath.Join(tmp, "video.bin")
	videoData := make([]byte, 2500)
	for i := range videoData {
		videoData[i] = byte(i % 251)
	}
	if err := os.WriteFile(videoPath, videoData, 0644); err != nil {
		t.Fatalf("write test video: %v", err)
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	if err := writer.WriteField("title", "test video"); err != nil {
		t.Fatalf("write title field: %v", err)
	}
	if err := writer.WriteField("description", "test desc"); err != nil {
		t.Fatalf("write description field: %v", err)
	}
	if err := writer.WriteField("tags", "tag1,tag2"); err != nil {
		t.Fatalf("write tags field: %v", err)
	}

	part, err := writer.CreateFormFile("file", "video.bin")
	if err != nil {
		t.Fatalf("create form file: %v", err)
	}
	if _, err := part.Write(videoData); err != nil {
		t.Fatalf("write form file: %v", err)
	}
	hash := sha256.Sum256(videoData)
	videoHash := hex.EncodeToString(hash[:])
	ts := time.Now().Unix()
	msg := []byte(videoHash + "|" + fmt.Sprintf("%d", ts) + "|" + pubB64)
	sigUpload := base64.StdEncoding.EncodeToString(ed25519.Sign(priv, msg))

	if err := writer.WriteField("video_hash", videoHash); err != nil {
		t.Fatalf("write video_hash: %v", err)
	}
	if err := writer.WriteField("author_timestamp", fmt.Sprintf("%d", ts)); err != nil {
		t.Fatalf("write author_timestamp: %v", err)
	}
	if err := writer.WriteField("author_signature", sigUpload); err != nil {
		t.Fatalf("write author_signature: %v", err)
	}

	if err := writer.Close(); err != nil {
		t.Fatalf("close writer: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("upload status = %d, body = %s", w.Code, w.Body.String())
	}

	var up uploadResp
	if err := json.Unmarshal(w.Body.Bytes(), &up); err != nil {
		t.Fatalf("decode upload response: %v", err)
	}

	if up.Video.ID == "" || up.Video.StorageID == "" {
		t.Fatalf("missing ids in upload response: %+v", up.Video)
	}
	if len(up.Video.Chunks) != 3 {
		t.Fatalf("expected 3 chunks, got %d", len(up.Video.Chunks))
	}
	if _, err := os.Stat(up.Video.Manifest); err != nil {
		t.Fatalf("manifest not found: %v", err)
	}
	if up.Video.AuthorPublicKey == "" {
		t.Fatalf("missing author_public_key in upload response")
	}
	if up.Video.ManifestHash == "" {
		t.Fatalf("missing manifest_hash in upload response")
	}

	manifestData, err := os.ReadFile(up.Video.Manifest)
	if err != nil {
		t.Fatalf("read manifest: %v", err)
	}
	var manifest struct {
		AuthorPublicKey string `json:"author_public_key"`
		Signature       string `json:"signature"`
	}
	if err := json.Unmarshal(manifestData, &manifest); err != nil {
		t.Fatalf("parse manifest: %v", err)
	}
	if manifest.AuthorPublicKey == "" || manifest.Signature == "" {
		t.Fatalf("manifest missing author signature fields")
	}

	getReq := httptest.NewRequest(http.MethodGet, "/video/"+up.Video.ID, nil)
	getW := httptest.NewRecorder()
	r.ServeHTTP(getW, getReq)
	if getW.Code != http.StatusOK {
		t.Fatalf("get video status = %d, body = %s", getW.Code, getW.Body.String())
	}

	streamReq := httptest.NewRequest(http.MethodGet, "/video/"+up.Video.ID+"/stream", nil)
	streamW := httptest.NewRecorder()
	r.ServeHTTP(streamW, streamReq)
	if streamW.Code != http.StatusOK {
		t.Fatalf("stream video status = %d, body = %s", streamW.Code, streamW.Body.String())
	}
}

func mustDecodeB64(v string) []byte {
	b, _ := base64.StdEncoding.DecodeString(v)
	return b
}
