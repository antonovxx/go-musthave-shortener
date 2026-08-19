package repository

import (
	"bufio"
	"encoding/json"
	"errors"
	"os"
	"strconv"
)

type urlRecord struct {
	UUID        string `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

type fileStorage struct {
	file    *os.File
	writer  *bufio.Writer
	encoder *json.Encoder
}

func newFileStorage(path string) (*fileStorage, error) {
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)

	if err != nil {
		return nil, err
	}

	writer := bufio.NewWriter(file)

	return &fileStorage{
		file:    file,
		writer:  writer,
		encoder: json.NewEncoder(writer),
	}, nil
}

func (fs *fileStorage) write(uuid, shortURL, originalURL string) error {
	if err := fs.encoder.Encode(urlRecord{
		UUID:        uuid,
		ShortURL:    shortURL,
		OriginalURL: originalURL,
	}); err != nil {
		return err
	}

	return fs.writer.Flush()
}

func (fs *fileStorage) close() error {
	return fs.file.Close()
}

func loadURLsFromFile(path string) (map[string]string, int, error) {
	file, err := os.OpenFile(path, os.O_RDONLY, 0o644)
	if errors.Is(err, os.ErrNotExist) {
		return make(map[string]string), 0, nil
	}

	if err != nil {
		return nil, 0, err
	}

	defer file.Close()
	urls := make(map[string]string)
	maxUUID := 0
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		var rec urlRecord
		if err := json.Unmarshal(scanner.Bytes(), &rec); err != nil {
			return nil, 0, err
		}

		urls[rec.ShortURL] = rec.OriginalURL
		if id, err := strconv.Atoi(rec.UUID); err == nil && id > maxUUID {
			maxUUID = id
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, 0, err
	}

	return urls, maxUUID, nil
}
