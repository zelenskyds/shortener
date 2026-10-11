package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestMakeRegistrUrl(t *testing.T) {
	tests := []struct {
		name       string
		method     string
		body       string
		wantStatus int
	}{
		{
			name:       "create short URL",
			method:     http.MethodPost,
			body:       "https://example.com/path?q=hello",
			wantStatus: http.StatusCreated,
		},
		{
			name:       "GET is not allowed",
			method:     http.MethodGet,
			wantStatus: http.StatusMethodNotAllowed,
		},
		{
			name:       "PUT is not allowed",
			method:     http.MethodPut,
			wantStatus: http.StatusMethodNotAllowed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := make(DB)
			req := httptest.NewRequest(tt.method, "/", strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "text/plain")
			res := httptest.NewRecorder()

			newRouter("localhost:8080", db).ServeHTTP(res, req)
			result := res.Result()
			defer result.Body.Close()

			if result.StatusCode != tt.wantStatus {
				t.Fatalf("status = %d, want %d", result.StatusCode, tt.wantStatus)
			}

			if tt.wantStatus != http.StatusCreated {
				if len(db) != 0 {
					t.Fatal("rejected request changed the database")
				}
				return
			}

			if got := result.Header.Get("Content-Type"); got != "text/plain" {
				t.Errorf("Content-Type = %q, want text/plain", got)
			}

			shortURL := res.Body.String()
			const prefix = "http://localhost:8080/"
			if !strings.HasPrefix(shortURL, prefix) {
				t.Fatalf("body = %q, want URL starting with %q", shortURL, prefix)
			}
			id := strings.TrimPrefix(shortURL, prefix)
			if id == "" || strings.ContainsAny(id, "/?\x23 \t\r\n") {
				t.Fatalf("invalid ID in response: %q", id)
			}
			if got := db[id]; got != tt.body {
				t.Errorf("saved URL = %q, want %q", got, tt.body)
			}
		})
	}
}

func TestMakeGetUrl(t *testing.T) {
	tests := []struct {
		name         string
		method       string
		path         string
		wantStatus   int
		wantLocation string
	}{
		{
			name:         "redirect to original URL",
			method:       http.MethodGet,
			path:         "/abc123",
			wantStatus:   http.StatusTemporaryRedirect,
			wantLocation: "https://example.com/path?q=hello",
		},
		{
			name:       "unknown ID",
			method:     http.MethodGet,
			path:       "/missing",
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "POST is not allowed",
			method:     http.MethodPost,
			path:       "/abc123",
			wantStatus: http.StatusMethodNotAllowed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := DB{"abc123": "https://example.com/path?q=hello"}
			mux := newRouter("localhost:8080", db)
			req := httptest.NewRequest(tt.method, tt.path, nil)
			res := httptest.NewRecorder()

			mux.ServeHTTP(res, req)
			result := res.Result()
			defer result.Body.Close()

			if result.StatusCode != tt.wantStatus {
				t.Errorf("status = %d, want %d", result.StatusCode, tt.wantStatus)
			}
			if got := result.Header.Get("Location"); got != tt.wantLocation {
				t.Errorf("Location = %q, want %q", got, tt.wantLocation)
			}
		})
	}
}
