package main

import (
	"antonovxx/go-musthave-shortener/internal/config"
	"antonovxx/go-musthave-shortener/internal/handler"
	"antonovxx/go-musthave-shortener/internal/middleware"
	"antonovxx/go-musthave-shortener/internal/repository"
	"antonovxx/go-musthave-shortener/internal/service"
	"errors"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

func main() {
	cfg, err := config.NewConfig()

	if err != nil {
		log.Fatal(err)
	}

	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Sync()

	repos := repository.NewURLRepository()
	shortener := service.NewShortenerService(repos, cfg.BaseURL)
	urlHandler := handler.NewURLHandler(shortener)

	r := chi.NewRouter()
	r.Use(middleware.Logger(logger))
	r.Post("/", urlHandler.HandleShortenURL)
	r.Post("/api/shorten", urlHandler.HandleShortenURLJSON)
	r.Get("/{id}", urlHandler.HandleExpandURL)

	logger.Info("Starting server", zap.String("address", cfg.ServerAddress))

	if err := http.ListenAndServe(cfg.ServerAddress, r); !errors.Is(err, http.ErrServerClosed) {
		logger.Fatal("server failed", zap.Error(err))
	}
}
