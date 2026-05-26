package repository

import (
	"errors"
	"sync"
)

type URLRepository struct {
	mutex sync.Mutex
	urls  map[string]string
}

func NewURLRepository() *URLRepository {
	return &URLRepository{
		urls: make(map[string]string)}
}

func (r *URLRepository) Get(id string) (string, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	url, ok := r.urls[id]
	if !ok {
		return "", errors.New("url not found")
	}
	return url, nil
}

func (r *URLRepository) SaveIfNotExists(id, originalURL string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.urls[id]; exists {
		return errors.New("url already exists")
	}

	r.urls[id] = originalURL
	return nil
}
