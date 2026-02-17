package service

import "math/rand"

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
}

func NewShortenerService(repository Repository) *ShortenerService {
	return &ShortenerService{
		repository: repository,
	}
}

func (s *ShortenerService) Shorten(originalURL string) string {
	id := generateID(idLength)
	s.repository.Save(id, originalURL)
	return id
}

func (s *ShortenerService) Resolve(id string) (string, bool) {
	return s.repository.Get(id)
}

func generateID(length int) string {
	b := make([]byte, length)

	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}

	return string(b)
}
