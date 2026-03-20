package service

import "video-index/model"

type VideoService struct {
	repo VideoRepository
}

type VideoRepository interface {
	Save(v model.Video)
	FindByID(id string) (model.Video, bool)
	FindAll() []model.Video
	Search(q string) []model.Video
}

func NewVideoService(repo VideoRepository) *VideoService {
	return &VideoService{repo: repo}
}

func (s *VideoService) Save(v model.Video) {
	s.repo.Save(v)
}

func (s *VideoService) Get(id string) (model.Video, bool) {
	return s.repo.FindByID(id)
}

func (s *VideoService) List() []model.Video {
	return s.repo.FindAll()
}

func (s *VideoService) Search(q string) []model.Video {
	return s.repo.Search(q)
}
