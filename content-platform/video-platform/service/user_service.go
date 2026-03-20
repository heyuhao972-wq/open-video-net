package service

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"video-platform/model"
	"video-platform/repository"
)

type UserService struct {
	repo       *repository.UserRepository
	jwtSecret  []byte
	challenges map[string]challenge
	lock       sync.Mutex
}

type challenge struct {
	Nonce     string
	ExpiresAt time.Time
}

func NewUserService(repo *repository.UserRepository, jwtSecret string) *UserService {
	return &UserService{
		repo:       repo,
		jwtSecret:  []byte(jwtSecret),
		challenges: make(map[string]challenge),
	}
}

func (s *UserService) Register(publicKey string) (model.User, error) {
	if publicKey == "" {
		return model.User{}, errors.New("public_key required")
	}

	if _, ok := s.repo.FindByPublicKey(publicKey); ok {
		return model.User{}, errors.New("public_key already exists")
	}

	if _, err := parsePublicKey(publicKey); err != nil {
		return model.User{}, errors.New("invalid public_key")
	}

	userID := userIDFromPublicKey(publicKey)
	user := model.User{
		ID:        userID,
		PublicKey: publicKey,
		CreatedAt: time.Now().Unix(),
	}

	s.repo.Save(user)
	return user, nil
}

func (s *UserService) CreateChallenge(publicKey string) (string, error) {
	user, ok := s.repo.FindByPublicKey(publicKey)
	if !ok {
		return "", errors.New("user not found")
	}

	_ = user
	nonceBytes := make([]byte, 32)
	if _, err := rand.Read(nonceBytes); err != nil {
		return "", err
	}
	nonce := base64.StdEncoding.EncodeToString(nonceBytes)

	s.lock.Lock()
	s.challenges[publicKey] = challenge{
		Nonce:     nonce,
		ExpiresAt: time.Now().Add(5 * time.Minute),
	}
	s.lock.Unlock()

	return nonce, nil
}

func (s *UserService) Login(publicKey string, nonce string, signature string) (string, model.User, error) {
	user, ok := s.repo.FindByPublicKey(publicKey)
	if !ok {
		return "", model.User{}, errors.New("invalid credentials")
	}

	if nonce == "" || signature == "" {
		return "", model.User{}, errors.New("nonce and signature required")
	}

	s.lock.Lock()
	ch, ok := s.challenges[publicKey]
	if ok && time.Now().After(ch.ExpiresAt) {
		delete(s.challenges, publicKey)
		ok = false
	}
	s.lock.Unlock()
	if !ok || ch.Nonce != nonce {
		return "", model.User{}, errors.New("invalid challenge")
	}

	nonceBytes, err := base64.StdEncoding.DecodeString(nonce)
	if err != nil {
		return "", model.User{}, errors.New("invalid nonce")
	}

	sigBytes, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		return "", model.User{}, errors.New("invalid signature")
	}

	pub, err := parsePublicKey(user.PublicKey)
	if err != nil {
		return "", model.User{}, errors.New("invalid public key")
	}

	if !ed25519.Verify(pub, nonceBytes, sigBytes) {
		return "", model.User{}, errors.New("invalid credentials")
	}

	claims := jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(s.jwtSecret)
	if err != nil {
		return "", model.User{}, err
	}

	return signed, user, nil
}

func (s *UserService) GetByID(id string) (model.User, bool) {
	return s.repo.FindByID(id)
}

func (s *UserService) UpdateNickname(id string, nickname string) (model.User, error) {
	if id == "" {
		return model.User{}, errors.New("user_id required")
	}
	if len(nickname) > 64 {
		return model.User{}, errors.New("nickname too long")
	}
	u, ok := s.repo.UpdateProfile(id, nickname, "", "")
	if !ok {
		return model.User{}, errors.New("user not found")
	}
	return u, nil
}

func (s *UserService) UpdateProfile(id string, nickname string, avatarURL string, bio string) (model.User, error) {
	if id == "" {
		return model.User{}, errors.New("user_id required")
	}
	if len(nickname) > 64 {
		return model.User{}, errors.New("nickname too long")
	}
	if len(avatarURL) > 512 {
		return model.User{}, errors.New("avatar url too long")
	}
	if len(bio) > 1000 {
		return model.User{}, errors.New("bio too long")
	}
	u, ok := s.repo.UpdateProfile(id, nickname, avatarURL, bio)
	if !ok {
		return model.User{}, errors.New("user not found")
	}
	return u, nil
}

func (s *UserService) SearchUsers(q string) []model.User {
	return s.repo.Search(q)
}

func (s *UserService) BanUser(id string, reason string) error {
	if id == "" {
		return errors.New("user_id required")
	}
	s.repo.Ban(id, reason)
	return nil
}

func (s *UserService) UnbanUser(id string) error {
	if id == "" {
		return errors.New("user_id required")
	}
	s.repo.Unban(id)
	return nil
}

func (s *UserService) IsBanned(id string) (bool, string) {
	return s.repo.IsBanned(id)
}

func (s *UserService) ListBans() map[string]string {
	return s.repo.ListBans()
}

func userIDFromPublicKey(publicKey string) string {
	sum := sha256.Sum256([]byte(publicKey))
	return hex.EncodeToString(sum[:])
}

func parsePublicKey(publicKey string) (ed25519.PublicKey, error) {
	raw, err := base64.StdEncoding.DecodeString(publicKey)
	if err != nil {
		return nil, err
	}
	if len(raw) == ed25519.PublicKeySize {
		return ed25519.PublicKey(raw), nil
	}
	parsed, err := x509.ParsePKIXPublicKey(raw)
	if err != nil {
		return nil, err
	}
	pub, ok := parsed.(ed25519.PublicKey)
	if !ok {
		return nil, errors.New("unsupported public key")
	}
	return pub, nil
}
