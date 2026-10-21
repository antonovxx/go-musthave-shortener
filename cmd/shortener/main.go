package main

import (
	"antonovxx/go-musthave-shortener/internal/config"
	"antonovxx/go-musthave-shortener/internal/handler"
	"antonovxx/go-musthave-shortener/internal/middleware"
	"antonovxx/go-musthave-shortener/internal/repository"
	"antonovxx/go-musthave-shortener/internal/service"
	"database/sql"
	"errors"
	"log"
	"net/http"

	_ "database/sql"

	"github.com/go-chi/chi/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
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

	repos, err := repository.NewURLRepository(cfg.FileStoragePath)
	if err != nil {

		logger.Fatal("failed to init repository", zap.Error(err))
	}
	defer repos.Close()

	var db *sql.DB
	var pinger handler.Pinger

	if cfg.DatabaseDSN != "" {
		db, err = sql.Open("pgx", cfg.DatabaseDSN)
		if err != nil {
			logger.Fatal("failed to open database", zap.Error(err))
		}
		pinger = db

		defer db.Close()
	}

	pingHandler := handler.NewPingHandler(pinger)
	shortener := service.NewShortenerService(repos, cfg.BaseURL)
	urlHandler := handler.NewURLHandler(shortener)

	r := chi.NewRouter()
	r.Use(middleware.Logger(logger))
	r.Use(middleware.Gzip)
	r.Get("/ping", pingHandler.HandlePing)
	r.Post("/", urlHandler.HandleShortenURL)
	r.Post("/api/shorten", urlHandler.HandleShortenURLJSON)
	r.Get("/{id}", urlHandler.HandleExpandURL)

	logger.Info("Starting server", zap.String("address", cfg.ServerAddress))

	if err := http.ListenAndServe(cfg.ServerAddress, r); !errors.Is(err, http.ErrServerClosed) {
		logger.Fatal("server failed", zap.Error(err))
	}
}
