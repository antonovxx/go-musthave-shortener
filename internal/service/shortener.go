package service

import (
	"errors"
	"math/rand"
	"net/url"
)

const (
	idLength = 8
	charset  = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

type Repository interface {
	Save(id, originalURL string)
	Get(id string) (originalURL string, ok bool)
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
	var id string

	for {
		id = generateID(idLength)
		if _, exists := s.repository.Get(id); !exists {
			break
		}
	}

	s.repository.Save(id, originalURL)
	return url.JoinPath(s.baseURL, id)
}

func (s *ShortenerService) Resolve(id string) (string, error) {
	originalURL, ok := s.repository.Get(id)
	if !ok {
		return "", errors.New("not found")
	}

	return originalURL, nil
}

func generateID(length int) string {
	b := make([]byte, length)

	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}

	return string(b)
}
