package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

//go:generate mockgen -destination=mocks/mock_service.go -package=mocks . Service
type Service interface {
	Shorten(originalURL string) (string, error)
	Resolve(id string) (string, error)
}

type URLHandler struct {
	service Service
}

func NewURLHandler(service Service) *URLHandler {
	return &URLHandler{
		service: service,
	}
}

func (h *URLHandler) HandleShortenURLJSON(w http.ResponseWriter, r *http.Request) {
	var req shortenRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.URL == "" {
		writeJSONError(w, http.StatusBadRequest, "url is required")
		return
	}

	shortenURL, err := h.service.Shorten(req.URL)

	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "failed to shorten url")
		return
	}

	resp := shortenResponse{Result: shortenURL}
	respBody, err := json.Marshal(resp)

	if err != nil {
		log.Printf("failed to marshal shorten response: %v", err)
		writeJSONError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if _, err := w.Write(respBody); err != nil {
		log.Printf("failed to write response: %v", err)
	}
}

func (h *URLHandler) HandleShortenURL(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if len(body) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	shortURL, err := h.service.Shorten(string(body))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)

	if _, err = fmt.Fprint(w, shortURL); err != nil {
		log.Printf("failed to write response: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *URLHandler) HandleExpandURL(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	originalURL, err := h.service.Resolve(id)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Location", originalURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func writeJSONError(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	body, err := json.Marshal(errorResponse{Error: message})

	if err != nil {
		log.Printf("failed to marshal error response: %v", err)
		return
	}

	if _, err := w.Write(body); err != nil {
		log.Printf("failed to write error response: %v", err)
	}
}
