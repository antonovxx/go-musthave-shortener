package main

import (
	"antonovxx/go-musthave-shortener/internal/handler"
	"antonovxx/go-musthave-shortener/internal/repository"
	"antonovxx/go-musthave-shortener/internal/service"
	"log"
	"net/http"
)

func main() {
	repos := repository.NewURLRepository()
	shortener := service.NewShortenerService(repos)
	urlHandler := handler.NewUrlHandler(shortener, "http://localhost:8080")

	mux := http.NewServeMux()
	mux.HandleFunc("POST /", urlHandler.HandlePost)
	mux.HandleFunc("GET /{id}", urlHandler.HandleGet)

	log.Println("Starting server on :8080")

	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
