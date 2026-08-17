package repository

import (
	"errors"
	"strconv"
	"sync"
)

type URLRepository struct {
	mutex    sync.Mutex
	urls     map[string]string
	storage  *fileStorage
	nextUUID int
}

func NewURLRepository(filePath string) (*URLRepository, error) {
	repo := &URLRepository{
		urls:     make(map[string]string),
		nextUUID: 1,
	}

	if filePath == "" {
		return repo, nil
	}

	urls, err := loadURLsFromFile(filePath)
	if err != nil {
		return nil, err
	}
	repo.urls = urls
	repo.nextUUID = len(urls) + 1

	storage, err := newFileStorage(filePath)
	if err != nil {
		return nil, err
	}
	repo.storage = storage

	return repo, nil
}

func (r *URLRepository) Close() error {
	if r.storage == nil {
		return nil
	}
	return r.storage.close()
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

	if r.storage != nil {
		uuid := strconv.Itoa(r.nextUUID)
		if err := r.storage.write(uuid, id, originalURL); err != nil {
			delete(r.urls, id)
			return err
		}
		r.nextUUID++
	}

	return nil
}
