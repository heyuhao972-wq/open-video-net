package service

import (
	"errors"
	"strings"
	"time"

	"video-platform/config"
	"video-platform/model"
	"video-platform/repository"
)

type CommentService struct {
	repo *repository.CommentRepository
}

func NewCommentService(repo *repository.CommentRepository) *CommentService {
	return &CommentService{repo: repo}
}

func (s *CommentService) Create(videoID string, userID string, content string, parentID int) (model.Comment, error) {
	if s.repo == nil {
		return model.Comment{}, errors.New("comment repo not initialized")
	}
	if strings.TrimSpace(videoID) == "" {
		return model.Comment{}, errors.New("video_id required")
	}
	if strings.TrimSpace(userID) == "" {
		return model.Comment{}, errors.New("user_id required")
	}
	content = strings.TrimSpace(content)
	if content == "" {
		return model.Comment{}, errors.New("content required")
	}
	if len(content) > 1000 {
		return model.Comment{}, errors.New("content too long")
	}

	status := "approved"
	if strings.EqualFold(config.LoadConfig().ModerationMode, "review") {
		status = "pending"
	}
	if parentID < 0 {
		parentID = 0
	}
	if parentID > 0 {
		parent, ok := s.repo.Get(parentID)
		if !ok || parent.VideoID != videoID {
			return model.Comment{}, errors.New("invalid parent_id")
		}
	}

	c, err := s.repo.Add(videoID, userID, content, parentID, status)
	if err != nil {
		return model.Comment{}, err
	}
	return c, nil
}

func (s *CommentService) List(videoID string) ([]model.Comment, error) {
	if s.repo == nil {
		return nil, errors.New("comment repo not initialized")
	}
	if strings.TrimSpace(videoID) == "" {
		return nil, errors.New("video_id required")
	}
	if strings.EqualFold(config.LoadConfig().ModerationMode, "review") {
		return s.repo.ListByVideo(videoID), nil
	}
	return s.repo.ListByVideoAll(videoID), nil
}

func (s *CommentService) Count(videoID string) (int, error) {
	if s.repo == nil {
		return 0, errors.New("comment repo not initialized")
	}
	if strings.TrimSpace(videoID) == "" {
		return 0, errors.New("video_id required")
	}
	if strings.EqualFold(config.LoadConfig().ModerationMode, "review") {
		return s.repo.CountByVideo(videoID), nil
	}
	return s.repo.CountByVideoAll(videoID), nil
}

func (s *CommentService) Delete(id int, userID string, allowAdmin bool) error {
	if s.repo == nil {
		return errors.New("comment repo not initialized")
	}
	if id <= 0 {
		return errors.New("invalid id")
	}
	c, ok := s.repo.Get(id)
	if !ok {
		return errors.New("comment not found")
	}
	if c.UserID != userID && !allowAdmin {
		return errors.New("permission denied")
	}
	s.repo.Delete(id)
	return nil
}

func (s *CommentService) Like(id int, userID string) (model.Comment, bool, error) {
	if s.repo == nil {
		return model.Comment{}, false, errors.New("comment repo not initialized")
	}
	if id <= 0 {
		return model.Comment{}, false, errors.New("invalid id")
	}
	if strings.TrimSpace(userID) == "" {
		return model.Comment{}, false, errors.New("user_id required")
	}
	c, ok, liked := s.repo.Like(id, userID)
	if !ok {
		return model.Comment{}, false, errors.New("comment not found")
	}
	if !liked {
		return c, false, errors.New("already liked")
	}
	return c, true, nil
}

func (s *CommentService) ListByStatus(status string) ([]model.Comment, error) {
	if s.repo == nil {
		return nil, errors.New("comment repo not initialized")
	}
	return s.repo.ListByStatus(status), nil
}

func (s *CommentService) Review(id int, status string, reason string, reviewer string) (model.Comment, bool, error) {
	if s.repo == nil {
		return model.Comment{}, false, errors.New("comment repo not initialized")
	}
	if id <= 0 {
		return model.Comment{}, false, errors.New("invalid id")
	}
	if status != "approved" && status != "rejected" && status != "pending" {
		return model.Comment{}, false, errors.New("invalid status")
	}
	c, ok := s.repo.UpdateReview(id, status, reason, reviewer, time.Now().Unix())
	if !ok {
		return model.Comment{}, false, errors.New("comment not found")
	}
	return c, true, nil
}
