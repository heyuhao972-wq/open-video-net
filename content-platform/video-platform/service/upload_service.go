package service

import (
	"crypto/ed25519"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

	videostorage "video-storage"

	"video-platform/index"
	"video-platform/model"
)

type UploadService struct {
	videoService *VideoService
	indexClient  *index.Client
}

func NewUploadService(videoService *VideoService, indexClient *index.Client) *UploadService {

	return &UploadService{
		videoService: videoService,
		indexClient:  indexClient,
	}

}

func (s *UploadService) RegisterVideoFromStorage(
	title string,
	description string,
	tags []string,
	filename string,
	coverPath string,
	storageID string,
	chunks []string,
	manifestURL string,
	manifestHash string,
	authorID string,
	authorPublicKey string,
	authorSignature string,
	authorTimestamp int64,
	videoHash string,
	platformID string,
) (model.Video, error) {
	if storageID == "" || manifestURL == "" {
		return model.Video{}, fmt.Errorf("storage_id and manifest_url required")
	}

	if err := s.verifyAuthorProof(authorPublicKey, authorSignature, authorTimestamp, videoHash); err != nil {
		return model.Video{}, fmt.Errorf("author signature invalid: %w", err)
	}

	if manifestHash == "" {
		hash, err := computeManifestHashFromURL(manifestURL, authorPublicKey, videoHash, authorTimestamp)
		if err != nil {
			return model.Video{}, fmt.Errorf("manifest hash failed: %w", err)
		}
		manifestHash = hash
	}

	video := s.videoService.CreateVideoWithStorage(
		title,
		description,
		tags,
		filename,
		"",
		coverPath,
		storageID,
		chunks,
		manifestURL,
		authorID,
		authorPublicKey,
		authorSignature,
		authorTimestamp,
		videoHash,
		manifestHash,
		platformID,
	)

	if s.indexClient != nil {
		if err := s.indexClient.UpsertVideo(video); err != nil {
			fmt.Printf("index update failed: %v\n", err)
		}
	}

	return video, nil

}

func (s *UploadService) UpdateIndex(video model.Video) error {
	if s.indexClient == nil {
		return nil
	}
	return s.indexClient.UpsertVideo(video)
}

func (s *UploadService) verifyAuthorProof(publicKey string, signatureB64 string, timestamp int64, videoHash string) error {
	if publicKey == "" || signatureB64 == "" || videoHash == "" || timestamp <= 0 {
		return errors.New("missing proof fields")
	}

	pub, err := parsePublicKey(publicKey)
	if err != nil {
		return err
	}

	msg := buildProofMessage(videoHash, timestamp, publicKey)
	sigBytes, err := base64.StdEncoding.DecodeString(signatureB64)
	if err != nil {
		return err
	}
	if !ed25519.Verify(pub, []byte(msg), sigBytes) {
		return errors.New("signature verify failed")
	}
	return nil
}

func buildProofMessage(videoHash string, timestamp int64, publicKey string) string {
	return strings.Join([]string{videoHash, strconv.FormatInt(timestamp, 10), publicKey}, "|")
}

func computeManifestHashFromURL(url string, authorPublicKey string, videoHash string, timestamp int64) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("manifest fetch status: %d", resp.StatusCode)
	}
	tmp, err := os.CreateTemp("", "manifest-*.json")
	if err != nil {
		return "", err
	}
	defer os.Remove(tmp.Name())
	if _, err := io.Copy(tmp, resp.Body); err != nil {
		_ = tmp.Close()
		return "", err
	}
	if err := tmp.Close(); err != nil {
		return "", err
	}
	return videostorage.ComputeManifestHash(tmp.Name(), authorPublicKey, videoHash, timestamp)
}
