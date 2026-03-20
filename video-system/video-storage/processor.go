package videostorage

import (
	"os"
	"path/filepath"
	"time"

	"video-storage/internal/chunk"
	"video-storage/internal/manifest"
	"video-storage/internal/storage"
)

type ProcessResult struct {
	VideoID      string
	VideoHash    string
	Timestamp    int64
	ChunkHashes  []string
	ChunkDir     string
	ManifestPath string
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
	videoHash, err := HashFile(filePath)
	if err != nil {
		return ProcessResult{}, err
	}

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

	ts := time.Now().Unix()
	m, err := manifest.BuildManifest(chunks, videoHash, ts)
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
		VideoHash:    videoHash,
		Timestamp:    ts,
		ChunkHashes:  hashes,
		ChunkDir:     chunkDir,
		ManifestPath: manifestPath,
	}, nil
}

func (p *Processor) GetChunk(hash string) ([]byte, error) {
	path := filepath.Join(p.basePath, "chunks", hash)
	return os.ReadFile(path)
}
