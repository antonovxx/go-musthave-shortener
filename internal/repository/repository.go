package repository

import (
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

func (r *URLRepository) Get(id string) (string, bool) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	url, ok := r.urls[id]
	return url, ok
}

func (r *URLRepository) SaveIfNotExists(id, originalURL string) bool {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.urls[id]; exists {
		return false
	}

	r.urls[id] = originalURL
	return true
}
