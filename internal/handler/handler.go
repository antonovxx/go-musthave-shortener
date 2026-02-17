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

type UrlHandler struct {
	service Service
	baseUrl string
}

func NewUrlHandler(service Service, baseUrl string) *UrlHandler {
	return &UrlHandler{
		service: service,
		baseUrl: baseUrl,
	}
}

func (h *UrlHandler) HandlePost(w http.ResponseWriter, r *http.Request) {
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

func (h *UrlHandler) HandleGet(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	originalUrl, ok := h.service.Resolve(id)

	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Location", originalUrl)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
