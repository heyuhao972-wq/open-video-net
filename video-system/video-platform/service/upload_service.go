package service

import (
	"encoding/hex"
	"fmt"
	"log"

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

func (s *UploadService) UploadVideo(title string, description string, tags []string, filePath string, filename string, authorID string, authorPublicKey string, authorSignature string, videoHash string, proofTimestamp int64, platformID string) (model.Video, error) {
	result, err := s.storage.Upload(filePath)
	if err != nil {
		return model.Video{}, fmt.Errorf("store video failed: %w", err)
	}

	if err := s.VerifyHash(videoHash, result.VideoHash); err != nil {
		return model.Video{}, err
	}

	if err := s.storage.SetManifestProof(result.ManifestPath, authorPublicKey, authorSignature, videoHash, proofTimestamp); err != nil {
		return model.Video{}, fmt.Errorf("set manifest signature failed: %w", err)
	}

	manifestHash, err := s.storage.ComputeManifestHash(result.ManifestPath, authorPublicKey)
	if err != nil {
		return model.Video{}, fmt.Errorf("compute manifest hash failed: %w", err)
	}

	video := s.videoService.CreateVideoWithStorage(
		title,
		description,
		tags,
		filename,
		filePath,
		result.StorageID,
		result.Chunks,
		result.ManifestPath,
		authorID,
		authorPublicKey,
		manifestHash,
		platformID,
		videoHash,
		proofTimestamp,
		authorSignature,
	)

	if s.indexClient != nil {
		if err := s.indexClient.UpsertVideo(video); err != nil {
			log.Printf("index update failed: %v", err)
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

func (s *UploadService) VerifyHash(expected string, actual string) error {
	if expected == "" || actual == "" {
		return fmt.Errorf("missing video hash")
	}
	if expected != actual {
		return fmt.Errorf("video hash mismatch")
	}
	_, err := hex.DecodeString(expected)
	if err != nil {
		return fmt.Errorf("invalid video hash")
	}
	return nil
}
