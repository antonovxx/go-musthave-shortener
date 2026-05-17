package main

import (
	"antonovxx/go-musthave-shortener/internal/config"
	"antonovxx/go-musthave-shortener/internal/handler"
	"antonovxx/go-musthave-shortener/internal/repository"
	"antonovxx/go-musthave-shortener/internal/service"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func main() {
	cfg := config.NewConfig()

	repos := repository.NewURLRepository()
	shortener := service.NewShortenerService(repos)
	urlHandler := handler.NewURLHandler(shortener, cfg.BaseURL)

	r := chi.NewRouter()
	r.Post("/", urlHandler.HandlePost)
	r.Get("/{id}", urlHandler.HandleGet)

	log.Printf("Starting server on %s\n", cfg.ServerAddress)

	if err := http.ListenAndServe(cfg.ServerAddress, r); err != nil {
		log.Fatal(err)
	}
}
