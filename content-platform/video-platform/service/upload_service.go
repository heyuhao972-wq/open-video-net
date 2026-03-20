package service

import (
	"crypto/ed25519"
	"encoding/base64"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"video-platform/index"
	"video-platform/model"
	"video-platform/storage"
)

type UploadService struct {
	videoService *VideoService
	storage      *storage.StorageClient
	indexClient  *index.Client
}

func NewUploadService(videoService *VideoService, storageClient *storage.StorageClient, indexClient *index.Client) *UploadService {

	return &UploadService{
		videoService: videoService,
		storage:      storageClient,
		indexClient:  indexClient,
	}

}

func (s *UploadService) UploadVideo(title string, description string, tags []string, filePath string, filename string, coverPath string, authorID string, authorPublicKey string, authorSignature string, authorTimestamp int64, videoHash string, platformID string) (model.Video, error) {
	result, err := s.storage.Upload(filePath)
	if err != nil {
		return model.Video{}, fmt.Errorf("store video failed: %w", err)
	}

	if result.VideoHash != "" && videoHash != "" && result.VideoHash != videoHash {
		return model.Video{}, fmt.Errorf("video_hash mismatch")
	}
	if videoHash == "" {
		videoHash = result.VideoHash
	}

	if err := s.verifyAuthorProof(authorPublicKey, authorSignature, authorTimestamp, videoHash); err != nil {
		return model.Video{}, fmt.Errorf("author signature invalid: %w", err)
	}

	if err := s.storage.SetManifestProof(result.ManifestPath, authorPublicKey, authorSignature, videoHash, authorTimestamp); err != nil {
		return model.Video{}, fmt.Errorf("set manifest proof failed: %w", err)
	}

	manifestHash, err := s.storage.ComputeManifestHash(result.ManifestPath, authorPublicKey, videoHash, authorTimestamp)
	if err != nil {
		return model.Video{}, fmt.Errorf("manifest hash failed: %w", err)
	}

	video := s.videoService.CreateVideoWithStorage(
		title,
		description,
		tags,
		filename,
		filePath,
		coverPath,
		result.StorageID,
		result.Chunks,
		result.ManifestPath,
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
