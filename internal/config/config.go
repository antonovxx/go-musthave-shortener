package config

import (
	"flag"
	"os"
	"path/filepath"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	ServerAddress   string `env:"SERVER_ADDRESS"`
	BaseURL         string `env:"BASE_URL"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
}

func NewConfig() (*Config, error) {
	cfg := &Config{}

	defaultStoragePath := filepath.Join(os.TempDir(), "short-url-db.json")

	flag.StringVar(&cfg.ServerAddress, "a", "localhost:8080", "server address")
	flag.StringVar(&cfg.BaseURL, "b", "http://localhost:8080", "base URL")
	flag.StringVar(&cfg.FileStoragePath, "f", defaultStoragePath, "file storage path")

	flag.Parse()

	if err := env.Parse(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
