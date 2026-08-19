package repository

import (
	"path/filepath"
	"testing"
)

func TestURLRepository_PersistsAndRestores(t *testing.T) {
	path := filepath.Join(t.TempDir(), "storage.json")

	repo1, err := NewURLRepository(path)
	if err != nil {
		t.Fatalf("failed to create repository: %v", err)
	}

	if err := repo1.SaveIfNotExists("abc123", "https://practicum.yandex.ru"); err != nil {
		t.Fatalf("failed to save: %v", err)
	}

	if err := repo1.Close(); err != nil {
		t.Fatalf("failed to close: %v", err)
	}

	repo2, err := NewURLRepository(path)
	if err != nil {
		t.Fatalf("failed to reopen repository: %v", err)
	}
	defer repo2.Close()

	originalURL, err := repo2.Get("abc123")
	if err != nil {
		t.Fatalf("expected url to be restored, got error: %v", err)
	}

	if originalURL != "https://practicum.yandex.ru" {
		t.Errorf("got %v, want %v", originalURL, "https://practicum.yandex.ru")
	}
}

func TestURLRepository_EmptyPath_InMemoryOnly(t *testing.T) {
	repo, err := NewURLRepository("")
	if err != nil {
		t.Fatalf("failed to create repository: %v", err)
	}

	if err := repo.SaveIfNotExists("xyz789", "https://ya.ru"); err != nil {
		t.Fatalf("failed to save: %v", err)
	}

	originalURL, err := repo.Get("xyz789")
	if err != nil || originalURL != "https://ya.ru" {
		t.Errorf("got (%v, %v), want (https://ya.ru, nil)", originalURL, err)
	}
}
