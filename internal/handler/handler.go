package handler

import (
	"fmt"
	"io"
	"net/http"
)

type Service interface {
	Shorten(originalUrl string) string
	Resolve(id string) (string, bool)
}

type URLHandler struct {
	service Service
	baseUrl string
}

func NewURLHandler(service Service, baseUrl string) *URLHandler {
	return &URLHandler{
		service: service,
		baseUrl: baseUrl,
	}
}

func (h *URLHandler) HandlePost(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	id := h.service.Shorten(string(body))

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "%s/%s", h.baseUrl, id)
}

func (h *URLHandler) HandleGet(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

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
