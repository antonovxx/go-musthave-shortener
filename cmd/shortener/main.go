package main

import (
	"antonovxx/go-musthave-shortener/internal/config"
	"antonovxx/go-musthave-shortener/internal/handler"
	"antonovxx/go-musthave-shortener/internal/repository"
	"antonovxx/go-musthave-shortener/internal/service"
	"errors"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func main() {
	cfg := config.NewConfig()

	repos := repository.NewURLRepository()
	shortener := service.NewShortenerService(repos, cfg.BaseURL)
	urlHandler := handler.NewURLHandler(shortener)

	r := chi.NewRouter()
	r.Post("/", urlHandler.HandleShortenURL)
	r.Get("/{id}", urlHandler.HandleExpandURL)

	log.Printf("Starting server on %s\n", cfg.ServerAddress)

	if err := http.ListenAndServe(cfg.ServerAddress, r); !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}
}
