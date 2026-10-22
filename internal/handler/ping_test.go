package handler

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

type fakePinger struct {
	err error
}

func (f fakePinger) PingContext(context.Context) error {
	return f.err
}

func TestHandlePing(t *testing.T) {
	tests := []struct {
		name       string
		pinger     Pinger
		wantStatus int
	}{
		{
			name:       "database available",
			pinger:     fakePinger{},
			wantStatus: http.StatusOK,
		},
		{
			name:       "database unavailable",
			pinger:     fakePinger{err: errors.New("connection refused")},
			wantStatus: http.StatusInternalServerError,
		},
		{
			name:       "database not configured",
			pinger:     nil,
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewPingHandler(tt.pinger)
			req := httptest.NewRequest(http.MethodGet, "/ping", nil)
			writer := httptest.NewRecorder()

			h.HandlePing(writer, req)

			result := writer.Result()
			defer result.Body.Close()

			if result.StatusCode != tt.wantStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", result.StatusCode, tt.wantStatus)
			}
		})
	}
}
