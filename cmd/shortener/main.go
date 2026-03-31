package main

import (
	"antonovxx/go-musthave-shortener/internal/handler"
	"antonovxx/go-musthave-shortener/internal/repository"
	"antonovxx/go-musthave-shortener/internal/service"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func main() {
	repos := repository.NewURLRepository()
	shortener := service.NewShortenerService(repos)
	urlHandler := handler.NewURLHandler(shortener, "http://localhost:8080")

	r := chi.NewRouter()
	r.Post("/", urlHandler.HandlePost)
	r.Get("/{id}", urlHandler.HandleGet)

	log.Println("Starting server on :8080")

	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}
