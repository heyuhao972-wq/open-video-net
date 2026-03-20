package service

import (
	"time"

	"github.com/google/uuid"

	"video-platform/model"
	"video-platform/repository"
)

type VideoService struct {
	repo *repository.VideoRepository
}

func NewVideoService(repo *repository.VideoRepository) *VideoService {

	return &VideoService{
		repo: repo,
	}

}

func (s *VideoService) CreateVideo(title string, filename string) model.Video {
	return s.CreateVideoWithStorage(title, "", nil, filename, "", "", nil, "", "", "", "", "", "", 0, "")
}

func (s *VideoService) CreateVideoWithStorage(
	title string,
	description string,
	tags []string,
	filename string,
	filePath string,
	storageID string,
	chunks []string,
	manifest string,
	authorID string,
	authorPublicKey string,
	manifestHash string,
	platformID string,
	videoHash string,
	proofTimestamp int64,
	authorSignature string,
) model.Video {
	video := model.Video{
		ID:              uuid.New().String(),
		PlatformID:      platformID,
		StorageID:       storageID,
		Title:           title,
		Description:     description,
		Filename:        filename,
		FilePath:        filePath,
		Tags:            tags,
		AuthorID:        authorID,
		AuthorPublicKey: authorPublicKey,
		AuthorSignature: authorSignature,
		ProofTimestamp:  proofTimestamp,
		VideoHash:       videoHash,
		Chunks:          chunks,
		Manifest:        manifest,
		ManifestHash:    manifestHash,
		CreatedAt:       time.Now().Unix(),
	}

	_ = s.repo.Save(video)

	return video
}

func (s *VideoService) ListVideos() []model.Video {

	return s.repo.FindAll()

}

func (s *VideoService) GetVideo(id string) (model.Video, bool) {

	return s.repo.FindByID(id)

}
