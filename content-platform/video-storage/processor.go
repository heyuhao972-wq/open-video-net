package videostorage

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
	"path/filepath"
	"time"

	"video-storage/internal/chunk"
	"video-storage/internal/manifest"
	"video-storage/internal/storage"
)

type ProcessResult struct {
	VideoID      string
	ChunkHashes  []string
	ChunkDir     string
	ManifestPath string
	VideoHash    string
	Timestamp    int64
}

type Processor struct {
	basePath  string
	chunkSize int
}

func NewProcessor(basePath string, chunkSize int) (*Processor, error) {
	if chunkSize <= 0 {
		chunkSize = 1024 * 1024
	}

	if err := os.MkdirAll(basePath, os.ModePerm); err != nil {
		return nil, err
	}

	return &Processor{
		basePath:  basePath,
		chunkSize: chunkSize,
	}, nil
}

func (p *Processor) StoreVideo(filePath string) (ProcessResult, error) {
	videoHash, err := computeFileHash(filePath)
	if err != nil {
		return ProcessResult{}, err
	}
	timestamp := time.Now().Unix()

	chunks, err := chunk.SplitFile(filePath, p.chunkSize)
	if err != nil {
		return ProcessResult{}, err
	}

	chunkDir := filepath.Join(p.basePath, "chunks")
	store, err := storage.NewLocalStore(chunkDir)
	if err != nil {
		return ProcessResult{}, err
	}

	hashes := make([]string, 0, len(chunks))
	for _, c := range chunks {
		if err := store.Save(c); err != nil {
			return ProcessResult{}, err
		}
		hashes = append(hashes, c.Hash)
	}

	m, err := manifest.BuildManifest(chunks, videoHash, timestamp)
	if err != nil {
		return ProcessResult{}, err
	}

	manifestDir := filepath.Join(p.basePath, "manifests")
	if err := os.MkdirAll(manifestDir, os.ModePerm); err != nil {
		return ProcessResult{}, err
	}

	manifestPath := filepath.Join(manifestDir, m.VideoID+".json")
	if err := m.Save(manifestPath); err != nil {
		return ProcessResult{}, err
	}

	return ProcessResult{
		VideoID:      m.VideoID,
		ChunkHashes:  hashes,
		ChunkDir:     chunkDir,
		ManifestPath: manifestPath,
		VideoHash:    videoHash,
		Timestamp:    timestamp,
	}, nil
}

func (p *Processor) GetChunk(hash string) ([]byte, error) {
	path := filepath.Join(p.basePath, "chunks", hash)
	return os.ReadFile(path)
}

func computeFileHash(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}
