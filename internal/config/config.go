package config

import (
	"flag"

	"github.com/caarlos0/env/v11"
	_ "github.com/caarlos0/env/v11"
)

type Config struct {
	ServerAddress string `env:"SERVER_ADDRESS"`
	BaseURL       string `env:"BASE_URL"`
}

func NewConfig() (*Config, error) {
	cfg := &Config{}

	flag.StringVar(&cfg.ServerAddress, "a", "localhost:8080", "server address")
	flag.StringVar(&cfg.BaseURL, "b", "http://localhost:8080", "base URL")

	flag.Parse()

	if err := env.Parse(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
