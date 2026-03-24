package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	videostorage "video-storage"
)

type config struct {
	port      string
	basePath  string
	chunkSize int
}

func loadConfig() config {
	return config{
		port:      envString("STORAGE_PORT", "8085"),
		basePath:  envString("STORAGE_PATH", "./data/video-storage"),
		chunkSize: envInt("CHUNK_SIZE", 1024*1024),
	}
}

func envString(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func envInt(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			return n
		}
	}
	return fallback
}

func main() {
	cfg := loadConfig()

	mux := http.NewServeMux()

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	mux.HandleFunc("/store", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			writeCORS(w)
			w.WriteHeader(http.StatusNoContent)
			return
		}
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		writeCORS(w)
		if err := r.ParseMultipartForm(500 << 20); err != nil {
			http.Error(w, "invalid form", http.StatusBadRequest)
			return
		}
		file, header, err := r.FormFile("file")
		if err != nil {
			http.Error(w, "file required", http.StatusBadRequest)
			return
		}
		defer file.Close()

		tmpDir := filepath.Join(cfg.basePath, "tmp")
		if err := os.MkdirAll(tmpDir, os.ModePerm); err != nil {
			http.Error(w, "create temp dir failed", http.StatusInternalServerError)
			return
		}
		tmpFile, err := os.CreateTemp(tmpDir, "upload-*")
		if err != nil {
			http.Error(w, "create temp file failed", http.StatusInternalServerError)
			return
		}
		if _, err := io.Copy(tmpFile, file); err != nil {
			tmpFile.Close()
			http.Error(w, "save temp failed", http.StatusInternalServerError)
			return
		}
		tmpFile.Close()

		processor, err := videostorage.NewProcessor(cfg.basePath, cfg.chunkSize)
		if err != nil {
			http.Error(w, "init storage failed", http.StatusInternalServerError)
			return
		}

		result, err := processor.StoreVideo(tmpFile.Name())
		if err != nil {
			http.Error(w, "store failed", http.StatusInternalServerError)
			return
		}

		reqVideoHash := strings.TrimSpace(r.FormValue("video_hash"))
		if reqVideoHash != "" && result.VideoHash != reqVideoHash {
			http.Error(w, "video_hash mismatch", http.StatusBadRequest)
			return
		}
		videoHash := result.VideoHash
		if reqVideoHash != "" {
			videoHash = reqVideoHash
		}

		pubKey := strings.TrimSpace(r.FormValue("author_public_key"))
		signature := strings.TrimSpace(r.FormValue("author_signature"))
		tsRaw := strings.TrimSpace(r.FormValue("author_timestamp"))
		authorTs := int64(0)
		if tsRaw != "" {
			if ts, err := strconv.ParseInt(tsRaw, 10, 64); err == nil {
				authorTs = ts
			}
		}
		if pubKey != "" && signature != "" && authorTs > 0 {
			_ = videostorage.SetManifestProof(result.ManifestPath, pubKey, signature, videoHash, authorTs)
		}
		manifestHash, _ := videostorage.ComputeManifestHash(result.ManifestPath, pubKey, videoHash, authorTs)

		originalDir := filepath.Join(cfg.basePath, "originals")
		_ = os.MkdirAll(originalDir, os.ModePerm)
		origPath := filepath.Join(originalDir, result.VideoID)
		_ = os.Rename(tmpFile.Name(), origPath)

		writeJSON(w, http.StatusOK, map[string]interface{}{
			"storage_id":    result.VideoID,
			"chunks":        result.ChunkHashes,
			"manifest_url":  "/manifest/" + result.VideoID,
			"stream_url":    "/stream/" + result.VideoID,
			"video_hash":    videoHash,
			"timestamp":     authorTs,
			"manifest_hash": manifestHash,
			"filename":      header.Filename,
		})
	})

	mux.HandleFunc("/manifest/", func(w http.ResponseWriter, r *http.Request) {
		writeCORS(w)
		id := strings.TrimPrefix(r.URL.Path, "/manifest/")
		if id == "" {
			http.Error(w, "id required", http.StatusBadRequest)
			return
		}
		path := filepath.Join(cfg.basePath, "manifests", id+".json")
		http.ServeFile(w, r, path)
	})

	mux.HandleFunc("/chunk/", func(w http.ResponseWriter, r *http.Request) {
		writeCORS(w)
		hash := strings.TrimPrefix(r.URL.Path, "/chunk/")
		if hash == "" {
			http.Error(w, "hash required", http.StatusBadRequest)
			return
		}
		path := filepath.Join(cfg.basePath, "chunks", hash)
		http.ServeFile(w, r, path)
	})

	mux.HandleFunc("/stream/", func(w http.ResponseWriter, r *http.Request) {
		writeCORS(w)
		id := strings.TrimPrefix(r.URL.Path, "/stream/")
		if id == "" {
			http.Error(w, "id required", http.StatusBadRequest)
			return
		}
		path := filepath.Join(cfg.basePath, "originals", id)
		f, err := os.Open(path)
		if err != nil {
			http.Error(w, "file not found", http.StatusNotFound)
			return
		}
		defer f.Close()
		buf := make([]byte, 512)
		n, _ := f.Read(buf)
		contentType := http.DetectContentType(buf[:n])
		w.Header().Set("Content-Type", contentType)
		f.Seek(0, 0)
		io.Copy(w, f)
	})

	fmt.Println("storage service running on port:", cfg.port)
	http.ListenAndServe(":"+cfg.port, mux)
}

func writeCORS(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization")
}

func writeJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(payload)
}
