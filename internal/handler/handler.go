package handler

import (
	"fmt"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Service interface {
	Shorten(originalURL string) string
	Resolve(id string) (string, bool)
}

type URLHandler struct {
	service Service
	baseURL string
}

func NewURLHandler(service Service, baseURL string) *URLHandler {
	return &URLHandler{
		service: service,
		baseURL: baseURL,
	}
}

func (h *URLHandler) HandlePost(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)

	if len(body) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	id := h.service.Shorten(string(body))

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	_, err = fmt.Fprintf(w, "%s/%s", h.baseURL, id)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

func (h *URLHandler) HandleGet(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	originalURL, ok := h.service.Resolve(id)

	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Location", originalURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
