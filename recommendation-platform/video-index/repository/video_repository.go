package repository

import (
	"strings"
	"sync"

	"video-index/model"
)

type VideoRepository struct {
	byID map[string]model.Video
	lock sync.RWMutex
}

func NewVideoRepository() *VideoRepository {
	return &VideoRepository{
		byID: make(map[string]model.Video),
	}
}

func (r *VideoRepository) Save(v model.Video) {
	r.lock.Lock()
	defer r.lock.Unlock()

	r.byID[v.ID] = v
}

func (r *VideoRepository) FindByID(id string) (model.Video, bool) {
	r.lock.RLock()
	defer r.lock.RUnlock()

	v, ok := r.byID[id]
	return v, ok
}

func (r *VideoRepository) FindAll() []model.Video {
	r.lock.RLock()
	defer r.lock.RUnlock()

	out := make([]model.Video, 0, len(r.byID))
	for _, v := range r.byID {
		out = append(out, v)
	}
	return out
}

func (r *VideoRepository) Search(q string) []model.Video {
	r.lock.RLock()
	defer r.lock.RUnlock()

	q = strings.ToLower(strings.TrimSpace(q))
	if q == "" {
		return []model.Video{}
	}

	out := make([]model.Video, 0)
	for _, v := range r.byID {
		if strings.Contains(strings.ToLower(v.Title), q) ||
			strings.Contains(strings.ToLower(v.Description), q) ||
			matchTag(v.Tags, q) {
			out = append(out, v)
		}
	}
	return out
}

func matchTag(tags []string, q string) bool {
	for _, t := range tags {
		if strings.Contains(strings.ToLower(t), q) {
			return true
		}
	}
	return false
}
