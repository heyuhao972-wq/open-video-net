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
	}

	s.repo.Save(video)

	return video
}

func (s *VideoService) ListVideos() []model.Video {

	return s.repo.FindAll()

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

func (s *VideoService) Search(q string, tag string) []model.Video {
	return s.repo.Search(q, tag)
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
