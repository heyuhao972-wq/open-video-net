package service

import (
	"strings"
	"time"

	"github.com/google/uuid"

	"video-platform/config"
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
	return s.CreateVideoWithStorage(title, "", nil, filename, "", "", "", nil, "", "", "", "", 0, "", "", "")
}

func (s *VideoService) CreateVideoWithStorage(
	title string,
	description string,
	tags []string,
	filename string,
	filePath string,
	coverPath string,
	storageID string,
	chunks []string,
	manifest string,
	authorID string,
	authorPublicKey string,
	authorSignature string,
	authorTimestamp int64,
	videoHash string,
	manifestHash string,
	platformID string,
) model.Video {
	status := "approved"
	cfg := config.LoadConfig()
	if strings.EqualFold(cfg.ModerationMode, "review") {
		status = "pending"
	}
	video := model.Video{
		ID:              uuid.New().String(),
		PlatformID:      platformID,
		StorageID:       storageID,
		Title:           title,
		Description:     description,
		Filename:        filename,
		FilePath:        filePath,
		CoverPath:       coverPath,
		Tags:            tags,
		AuthorID:        authorID,
		AuthorPublicKey: authorPublicKey,
		AuthorSignature: authorSignature,
		AuthorTimestamp: authorTimestamp,
		VideoHash:       videoHash,
		Chunks:          chunks,
		Manifest:        manifest,
		ManifestHash:    manifestHash,
		CreatedAt:       time.Now().Unix(),
		Status:          status,
	}

	s.repo.Save(video)

	return video
}

func (s *VideoService) ListVideos() []model.Video {

	return s.repo.FindAll()

}

func (s *VideoService) ListVideosByStatus(status string) []model.Video {
	return s.repo.FindAllByStatus(status)
}

func (s *VideoService) GetVideo(id string) (model.Video, bool) {

	return s.repo.FindByID(id)

}

func (s *VideoService) ListByAuthor(authorID string) []model.Video {
	if authorID == "" {
		return nil
	}
	return s.repo.FindByAuthor(authorID)
}

func (s *VideoService) ListByAuthorAndStatus(authorID string, status string) []model.Video {
	if authorID == "" {
		return nil
	}
	return s.repo.FindByAuthorAndStatus(authorID, status)
}

func (s *VideoService) Search(q string, tag string) []model.Video {
	return s.repo.Search(q, tag)
}

func (s *VideoService) SearchByStatus(q string, tag string, status string) []model.Video {
	return s.repo.SearchByStatus(q, tag, status)
}

func (s *VideoService) DeleteVideo(id string) (model.Video, bool) {
	if id == "" {
		return model.Video{}, false
	}
	return s.repo.Delete(id)
}

func (s *VideoService) UpdateVideo(id string, title string, description string, tags []string) (model.Video, bool) {
	if id == "" {
		return model.Video{}, false
	}
	return s.repo.UpdateMeta(id, title, description, tags)
}

func (s *VideoService) ReviewVideo(id string, status string, reason string, reviewer string) (model.Video, bool) {
	if id == "" {
		return model.Video{}, false
	}
	if status != "approved" && status != "rejected" && status != "pending" {
		return model.Video{}, false
	}
	return s.repo.UpdateReview(id, status, reason, reviewer, time.Now().Unix())
}
