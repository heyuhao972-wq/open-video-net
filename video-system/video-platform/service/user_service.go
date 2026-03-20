package service

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/sha256"
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
	if _, err := base64.StdEncoding.DecodeString(publicKey); err != nil {
		return model.User{}, errors.New("invalid public_key")
	}

	userID := deriveUserID(publicKey)
	if _, ok := s.repo.FindByID(userID); ok {
		return model.User{}, errors.New("user already exists")
	}

	user := model.User{
		ID:        userID,
		PublicKey: publicKey,
		CreatedAt: time.Now().Unix(),
	}

	if err := s.repo.Save(user); err != nil {
		return model.User{}, err
	}
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
	s.challenges[user.ID] = challenge{
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
	ch, ok := s.challenges[user.ID]
	if ok && time.Now().After(ch.ExpiresAt) {
		delete(s.challenges, user.ID)
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

	pubBytes, err := base64.StdEncoding.DecodeString(user.PublicKey)
	if err != nil {
		return "", model.User{}, errors.New("invalid public key")
	}

	if !ed25519.Verify(ed25519.PublicKey(pubBytes), nonceBytes, sigBytes) {
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

func deriveUserID(publicKey string) string {
	sum := sha256.Sum256([]byte(publicKey))
	return hex.EncodeToString(sum[:])
}
