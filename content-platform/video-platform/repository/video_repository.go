package repository

import (
	"strings"
	"sync"

	"video-platform/model"
)

type VideoRepository struct {
	videos map[string]model.Video
	lock   sync.RWMutex
}

func NewVideoRepository() *VideoRepository {

	return &VideoRepository{
		videos: make(map[string]model.Video),
	}

}

func (r *VideoRepository) Save(video model.Video) {

	r.lock.Lock()
	defer r.lock.Unlock()

	r.videos[video.ID] = video

}

func (r *VideoRepository) FindAll() []model.Video {

	r.lock.RLock()
	defer r.lock.RUnlock()

	list := []model.Video{}

	for _, v := range r.videos {
		list = append(list, v)
	}

	return list
}

func (r *VideoRepository) FindByID(id string) (model.Video, bool) {

	r.lock.RLock()
	defer r.lock.RUnlock()

	v, ok := r.videos[id]

	return v, ok
}

func (r *VideoRepository) FindByAuthor(authorID string) []model.Video {
	r.lock.RLock()
	defer r.lock.RUnlock()

	list := []model.Video{}
	for _, v := range r.videos {
		if v.AuthorID == authorID {
			list = append(list, v)
		}
	}
	return list
}

func (r *VideoRepository) Delete(id string) (model.Video, bool) {
	r.lock.Lock()
	defer r.lock.Unlock()
	v, ok := r.videos[id]
	if !ok {
		return model.Video{}, false
	}
	delete(r.videos, id)
	return v, true
}

func (r *VideoRepository) UpdateMeta(id string, title string, description string, tags []string) (model.Video, bool) {
	r.lock.Lock()
	defer r.lock.Unlock()
	v, ok := r.videos[id]
	if !ok {
		return model.Video{}, false
	}
	if title != "" {
		v.Title = title
	}
	v.Description = description
	v.Tags = tags
	r.videos[id] = v
	return v, true
}

func (r *VideoRepository) Search(q string, tag string) []model.Video {
	r.lock.RLock()
	defer r.lock.RUnlock()

	out := []model.Video{}
	for _, v := range r.videos {
		if tag != "" {
			if !hasTag(v.Tags, tag) {
				continue
			}
		}
		if q == "" {
			out = append(out, v)
			continue
		}
		if containsFoldVideo(v.Title, q) || containsFoldVideo(v.Description, q) || containsFoldVideo(v.AuthorID, q) {
			out = append(out, v)
		}
	}
	return out
}

func hasTag(tags []string, tag string) bool {
	for _, t := range tags {
		if strings.EqualFold(t, tag) {
			return true
		}
	}
	return false
}

func containsFoldVideo(haystack string, needle string) bool {
	if needle == "" {
		return true
	}
	h := strings.ToLower(haystack)
	n := strings.ToLower(needle)
	return strings.Contains(h, n)
}
