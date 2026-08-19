package middleware

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func jsonHandler(w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(r.Body)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}

func TestGzip_SendsGzip(t *testing.T) {
	handler := Gzip(http.HandlerFunc(jsonHandler))
	srv := httptest.NewServer(handler)
	defer srv.Close()

	requestBody := `{"url":"https://practicum.yandex.ru"}`

	buf := bytes.NewBuffer(nil)
	zw := gzip.NewWriter(buf)
	zw.Write([]byte(requestBody))
	zw.Close()

	req := httptest.NewRequest(http.MethodPost, srv.URL, buf)
	req.RequestURI = ""
	req.Header.Set("Content-Encoding", "gzip")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	b, _ := io.ReadAll(resp.Body)
	if string(b) != requestBody {
		t.Errorf("got %v, want %v", string(b), requestBody)
	}
}

func TestGzip_AcceptsGzip(t *testing.T) {
	handler := Gzip(http.HandlerFunc(jsonHandler))
	srv := httptest.NewServer(handler)
	defer srv.Close()

	requestBody := `{"url":"https://practicum.yandex.ru"}`

	req := httptest.NewRequest(http.MethodPost, srv.URL, strings.NewReader(requestBody))
	req.RequestURI = ""
	req.Header.Set("Accept-Encoding", "gzip")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.Header.Get("Content-Encoding") != "gzip" {
		t.Fatalf("expected gzip Content-Encoding, got %q", resp.Header.Get("Content-Encoding"))
	}

	zr, err := gzip.NewReader(resp.Body)
	if err != nil {
		t.Fatalf("failed to create gzip reader: %v", err)
	}

	b, _ := io.ReadAll(zr)
	if string(b) != requestBody {
		t.Errorf("got %v, want %v", string(b), requestBody)
	}
}

func TestGzip_NonCompressibleContentType_NotCompressed(t *testing.T) {
	plainHandler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("http://localhost:8080/abc"))
	}

	handler := Gzip(http.HandlerFunc(plainHandler))
	srv := httptest.NewServer(handler)
	defer srv.Close()

	req := httptest.NewRequest(http.MethodGet, srv.URL, nil)
	req.RequestURI = ""
	req.Header.Set("Accept-Encoding", "gzip")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.Header.Get("Content-Encoding") == "gzip" {
		t.Errorf("expected no gzip compression for text/plain content type")
	}
}

func TestGzip_ImplicitWriteHeader_CompressesJSON(t *testing.T) {
	implicitHandler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"result":"http://localhost:8080/abc"}`))
	}

	handler := Gzip(http.HandlerFunc(implicitHandler))
	srv := httptest.NewServer(handler)
	defer srv.Close()

	expectedBody := `{"result":"http://localhost:8080/abc"}`

	req := httptest.NewRequest(http.MethodGet, srv.URL, nil)
	req.RequestURI = ""
	req.Header.Set("Accept-Encoding", "gzip")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("wrong status code: got %v want %v", resp.StatusCode, http.StatusOK)
	}

	if resp.Header.Get("Content-Encoding") != "gzip" {
		t.Fatalf("expected gzip Content-Encoding, got %q", resp.Header.Get("Content-Encoding"))
	}

	zr, err := gzip.NewReader(resp.Body)
	if err != nil {
		t.Fatalf("failed to create gzip reader: %v", err)
	}
	defer zr.Close()

	b, err := io.ReadAll(zr)
	if err != nil {
		t.Fatalf("failed to read decompressed body: %v", err)
	}

	if string(b) != expectedBody {
		t.Errorf("got %v, want %v", string(b), expectedBody)
	}
}
