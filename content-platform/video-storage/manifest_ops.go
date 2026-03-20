package videostorage

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"

	"video-storage/internal/manifest"
)

func ComputeManifestHash(path string, authorPublicKey string, videoHash string, timestamp int64) (string, error) {
	m, err := loadManifest(path)
	if err != nil {
		return "", err
	}

	m.AuthorPublicKey = authorPublicKey
	m.VideoHash = videoHash
	m.Timestamp = timestamp
	m.Signature = ""

	data, err := json.Marshal(m)
	if err != nil {
		return "", err
	}

	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:]), nil
}

func SetManifestProof(path string, authorPublicKey string, signature string, videoHash string, timestamp int64) error {
	m, err := loadManifest(path)
	if err != nil {
		return err
	}

	m.AuthorPublicKey = authorPublicKey
	m.Signature = signature
	m.VideoHash = videoHash
	m.Timestamp = timestamp
	return m.Save(path)
}

func loadManifest(path string) (manifest.VideoManifest, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return manifest.VideoManifest{}, err
	}

	var m manifest.VideoManifest
	if err := json.Unmarshal(data, &m); err != nil {
		return manifest.VideoManifest{}, err
	}

	return m, nil
}
