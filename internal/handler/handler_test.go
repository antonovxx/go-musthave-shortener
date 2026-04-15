package handler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
)

type mockService struct {
}

const mockID = "mockID"
const mockURL = "mockURL"
const baseURL = "http://localhost:8080"

func (m *mockService) Shorten(url string) string {
	return mockID
}

func (m *mockService) Resolve(id string) (string, bool) {
	if id == mockID {
		return mockURL, true
	}

	return "", false
}

func TestHandlePost_ValidValues_ExpectedSuccess(t *testing.T) {
	h := NewURLHandler(&mockService{}, baseURL)

	body := strings.NewReader(mockURL)
	req := httptest.NewRequest(http.MethodPost, "/", body)
	writer := httptest.NewRecorder()

	h.HandlePost(writer, req)

	result := writer.Result()
	defer result.Body.Close()

	if result.StatusCode != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v", result.StatusCode, http.StatusCreated)
	}

	contentType := result.Header.Get("Content-Type")

	const contentTextPlain = "text/plain"

	if contentType != contentTextPlain {
		t.Errorf("handler returned wrong content type: got %v want %v", contentType, contentTextPlain)
	}

	responseBody := writer.Body.String()
	if responseBody != "http://localhost:8080/mockID" {
		t.Errorf("handler returned wrong response body: got %v want %v", responseBody, `http://localhost:8080/mockURL`)
	}
}

func TestHandlePost_BadRequest_ExpectedEmptyBody(t *testing.T) {
	h := NewURLHandler(&mockService{}, baseURL)
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	writer := httptest.NewRecorder()

	h.HandlePost(writer, req)

	result := writer.Result()
	defer result.Body.Close()

	if result.StatusCode != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", result.StatusCode, http.StatusBadRequest)
	}
}

func TestHandleGet_IncorrectID_ExpectedBadRequest(t *testing.T) {
	h := NewURLHandler(&mockService{}, baseURL)
	req := httptest.NewRequest(http.MethodGet, "/unknown", nil)
	req = withChiParam(req, "id", "unknown")
	writer := httptest.NewRecorder()

	h.HandleGet(writer, req)

	result := writer.Result()
	defer result.Body.Close()

	if result.StatusCode != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", result.StatusCode, http.StatusBadRequest)
	}
}

func TestHandleGet_IncorrectURL_ExpectedBadRequest(t *testing.T) {
	h := NewURLHandler(&mockService{}, baseURL)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	writer := httptest.NewRecorder()

	h.HandleGet(writer, req)

	result := writer.Result()
	defer result.Body.Close()

	if result.StatusCode != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", result.StatusCode, http.StatusBadRequest)
	}
}

func withChiParam(r *http.Request, key, value string) *http.Request {
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add(key, value)

	return r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
}
