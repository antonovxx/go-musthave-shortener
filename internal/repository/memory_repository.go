package repository

import (
	"antonovxx/go-musthave-shortener/internal/model"
	"strconv"
	"sync"
)

type FileRepository struct {
	mutex    sync.Mutex
	urls     map[string]string
	storage  *fileStorage
	nextUUID int
}

func NewFileRepository(filePath string) (*FileRepository, error) {
	repo := &FileRepository{
		urls:     make(map[string]string),
		nextUUID: 1,
	}

	if filePath == "" {
		return repo, nil
	}

	urls, maxUUID, err := loadURLsFromFile(filePath)
	if err != nil {
		return nil, err
	}
	repo.urls = urls
	repo.nextUUID = maxUUID + 1

	storage, err := newFileStorage(filePath)
	if err != nil {
		return nil, err
	}
	repo.storage = storage

	return repo, nil
}

func (r *FileRepository) Close() error {
	if r.storage == nil {
		return nil
	}
	return r.storage.close()
}

func (r *FileRepository) Get(id string) (string, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	url, ok := r.urls[id]
	if !ok {
		return "", model.ErrNotFound
	}
	return url, nil
}

func (r *FileRepository) SaveIfNotExists(id, originalURL string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.urls[id]; exists {
		return model.ErrIDExists
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
