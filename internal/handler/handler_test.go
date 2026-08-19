package handler

import (
	"antonovxx/go-musthave-shortener/internal/handler/mocks"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"go.uber.org/mock/gomock"
)

const mockID = "mockID"
const mockURL = "mockURL"
const mockShortURL = "http://localhost:8080/" + mockID

func TestHandleShortenURL_ValidValues_ExpectedSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockService := mocks.NewMockService(ctrl)
	mockService.EXPECT().Shorten(mockURL).Return(mockShortURL, nil)

	h := NewURLHandler(mockService)
	body := strings.NewReader(mockURL)
	req := httptest.NewRequest(http.MethodPost, "/", body)
	writer := httptest.NewRecorder()

	h.HandleShortenURL(writer, req)

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

func TestHandleShortenURL_BadRequest_ExpectedEmptyBody(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockService := mocks.NewMockService(ctrl)

	h := NewURLHandler(mockService)
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	writer := httptest.NewRecorder()

	h.HandleShortenURL(writer, req)

	result := writer.Result()
	defer result.Body.Close()

	if result.StatusCode != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", result.StatusCode, http.StatusBadRequest)
	}
}

func TestHandleExpandURL_IncorrectID_ExpectedBadRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockService := mocks.NewMockService(ctrl)
	mockService.EXPECT().Resolve("unknown").Return("", errors.New("url not found"))

	h := NewURLHandler(mockService)
	req := httptest.NewRequest(http.MethodGet, "/unknown", nil)
	req = withChiParam(req, "id", "unknown")
	writer := httptest.NewRecorder()

	h.HandleExpandURL(writer, req)

	result := writer.Result()
	defer result.Body.Close()

	if result.StatusCode != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", result.StatusCode, http.StatusBadRequest)
	}
}

func TestHandleExpandURL_IncorrectURL_ExpectedBadRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockService := mocks.NewMockService(ctrl)

	h := NewURLHandler(mockService)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	writer := httptest.NewRecorder()

	h.HandleExpandURL(writer, req)

	result := writer.Result()
	defer result.Body.Close()

	if result.StatusCode != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", result.StatusCode, http.StatusBadRequest)
	}
}

func TestHandleShortenURLJSON_ValidValues_ExpectedSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockService := mocks.NewMockService(ctrl)
	mockService.EXPECT().Shorten(mockURL).Return(mockShortURL, nil)

	h := NewURLHandler(mockService)
	body := strings.NewReader(`{"url":"mockURL"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/shorten", body)
	writer := httptest.NewRecorder()

	h.HandleShortenURLJSON(writer, req)

	result := writer.Result()
	defer result.Body.Close()

	if result.StatusCode != http.StatusCreated {
		t.Errorf("wrong status code: got %v want %v", result.StatusCode, http.StatusCreated)
	}

	if ct := result.Header.Get("Content-Type"); ct != "application/json" {
		t.Errorf("wrong content type: got %v want %v", ct, "application/json")
	}

	var resp shortenResponse
	if err := json.NewDecoder(result.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.Result != mockShortURL {
		t.Errorf("wrong result: got %v want %v", resp.Result, mockShortURL)
	}
}

func TestHandleShortenURLJSON_InvalidJSON_ExpectedBadRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockService := mocks.NewMockService(ctrl)

	h := NewURLHandler(mockService)
	body := strings.NewReader(`not a json`)
	req := httptest.NewRequest(http.MethodPost, "/api/shorten", body)
	writer := httptest.NewRecorder()

	h.HandleShortenURLJSON(writer, req)

	result := writer.Result()
	defer result.Body.Close()

	if result.StatusCode != http.StatusBadRequest {
		t.Errorf("wrong status code: got %v want %v", result.StatusCode, http.StatusBadRequest)
	}

	assertJSONError(t, result, http.StatusBadRequest, "invalid request body")
}

func TestHandleShortenURLJSON_EmptyURL_ExpectedBadRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockService := mocks.NewMockService(ctrl)
	h := NewURLHandler(mockService)
	body := strings.NewReader(`{"url":""}`)
	req := httptest.NewRequest(http.MethodPost, "/api/shorten", body)
	writer := httptest.NewRecorder()
	h.HandleShortenURLJSON(writer, req)
	result := writer.Result()
	defer result.Body.Close()

	assertJSONError(t, result, http.StatusBadRequest, "url is required")
}

func TestHandleShortenURLJSON_ShortenError_ExpectedInternalServerError(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockService := mocks.NewMockService(ctrl)

	mockService.EXPECT().Shorten(mockURL).Return("", errors.New("failed to generate unique id"))

	h := NewURLHandler(mockService)
	body := strings.NewReader(`{"url":"mockURL"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/shorten", body)
	writer := httptest.NewRecorder()
	h.HandleShortenURLJSON(writer, req)
	result := writer.Result()
	defer result.Body.Close()

	assertJSONError(t, result, http.StatusInternalServerError, "failed to shorten url")
}

func assertJSONError(t *testing.T, result *http.Response, wantStatus int, wantMessage string) {
	t.Helper()

	if result.StatusCode != wantStatus {
		t.Errorf("wrong status code: got %v want %v", result.StatusCode, wantStatus)
	}

	if ct := result.Header.Get("Content-Type"); ct != "application/json" {
		t.Errorf("wrong content type: got %v want %v", ct, "application/json")
	}

	var resp errorResponse

	if err := json.NewDecoder(result.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode error response: %v", err)
	}

	if resp.Error != wantMessage {
		t.Errorf("wrong error message: got %v want %v", resp.Error, wantMessage)
	}
}

func withChiParam(r *http.Request, key, value string) *http.Request {
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add(key, value)

	return r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
}
