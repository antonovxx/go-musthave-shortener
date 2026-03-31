package repository

import (
	"sync"
)

type URLRepository struct {
	mutex sync.RWMutex
	urls  map[string]string
}

func NewURLRepository() *URLRepository {
	return &URLRepository{
		urls: make(map[string]string)}
}

func (r *URLRepository) Save(id, originalURL string) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.urls[id] = originalURL
}

func (r *URLRepository) Get(id string) (string, bool) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	url, ok := r.urls[id]
	return url, ok
}
