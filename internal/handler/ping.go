package handler

import (
	"context"
	"net/http"
	"time"
)

type Pinger interface {
	PingContext(ctx context.Context) error
}

type PingHandler struct {
	db Pinger
}

func NewPingHandler(db Pinger) *PingHandler {
	return &PingHandler{db: db}
}

func (h *PingHandler) HandlePing(w http.ResponseWriter, r *http.Request) {
	if h.db == nil {
		http.Error(w, "database is not configured", http.StatusInternalServerError)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	if err := h.db.PingContext(ctx); err != nil {
		http.Error(w, "database is unavailable", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
