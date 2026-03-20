package repository

import (
	"strings"
	"sync"

	"video-platform/model"
)

type UserRepository struct {
	byID       map[string]model.User
	byPublic   map[string]model.User
	banned     map[string]string
	lock       sync.RWMutex
}

func NewUserRepository() *UserRepository {
	return &UserRepository{
		byID:     make(map[string]model.User),
		byPublic: make(map[string]model.User),
		banned:   make(map[string]string),
	}
}

func (r *UserRepository) Save(user model.User) {
	r.lock.Lock()
	defer r.lock.Unlock()

	r.byID[user.ID] = user
	r.byPublic[user.PublicKey] = user
}

func (r *UserRepository) FindByPublicKey(publicKey string) (model.User, bool) {
	r.lock.RLock()
	defer r.lock.RUnlock()

	u, ok := r.byPublic[publicKey]
	return u, ok
}

func (r *UserRepository) FindByID(id string) (model.User, bool) {
	r.lock.RLock()
	defer r.lock.RUnlock()

	u, ok := r.byID[id]
	return u, ok
}

func (r *UserRepository) UpdateProfile(id string, nickname string, avatarURL string, bio string) (model.User, bool) {
	r.lock.Lock()
	defer r.lock.Unlock()

	u, ok := r.byID[id]
	if !ok {
		return model.User{}, false
	}
	u.Nickname = nickname
	u.AvatarURL = avatarURL
	u.Bio = bio
	r.byID[id] = u
	r.byPublic[u.PublicKey] = u
	return u, true
}

func (r *UserRepository) Search(q string) []model.User {
	r.lock.RLock()
	defer r.lock.RUnlock()

	out := []model.User{}
	for _, u := range r.byID {
		if q == "" {
			out = append(out, u)
			continue
		}
		if containsFold(u.Nickname, q) || containsFold(u.PublicKey, q) || containsFold(u.ID, q) {
			out = append(out, u)
		}
	}
	return out
}

func (r *UserRepository) Ban(userID string, reason string) {
	r.lock.Lock()
	defer r.lock.Unlock()
	if userID == "" {
		return
	}
	r.banned[userID] = reason
}

func (r *UserRepository) Unban(userID string) {
	r.lock.Lock()
	defer r.lock.Unlock()
	delete(r.banned, userID)
}

func (r *UserRepository) IsBanned(userID string) (bool, string) {
	r.lock.RLock()
	defer r.lock.RUnlock()
	reason, ok := r.banned[userID]
	return ok, reason
}

func (r *UserRepository) ListBans() map[string]string {
	r.lock.RLock()
	defer r.lock.RUnlock()
	out := map[string]string{}
	for k, v := range r.banned {
		out[k] = v
	}
	return out
}

func containsFold(haystack string, needle string) bool {
	if needle == "" {
		return true
	}
	h := strings.ToLower(haystack)
	n := strings.ToLower(needle)
	return strings.Contains(h, n)
}
