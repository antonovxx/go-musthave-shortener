package service

import (
	"errors"
	"math/rand"
	"net/url"
)

const (
	idLength            = 8
	maxGenerateAttempts = 5
	charset             = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

type Repository interface {
	SaveIfNotExists(id, originalURL string) error
	Get(id string) (originalURL string, err error)
}

type ShortenerService struct {
	repository Repository
	baseURL    string
}

func NewShortenerService(repository Repository, baseURL string) *ShortenerService {
	return &ShortenerService{
		repository: repository,
		baseURL:    baseURL,
	}
}

func (s *ShortenerService) Shorten(originalURL string) (string, error) {
	for range maxGenerateAttempts {
		id := generateID(idLength)
		if err := s.repository.SaveIfNotExists(id, originalURL); err == nil {
			return url.JoinPath(s.baseURL, id)
		}
	}
	return "", errors.New("failed to generate unique id")
}

func (s *ShortenerService) Resolve(id string) (string, error) {
	return s.repository.Get(id)
}

func generateID(length int) string {
	b := make([]byte, length)

	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}

	return string(b)
}
