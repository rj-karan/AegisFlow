package httpserver

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rj-karan/AegisFlow/internal/config"
)

func newTestServer(t *testing.T) *Server {
	t.Helper()
	cfg, err := config.Load(func(string) (string, bool) { return "", false })
	if err != nil {
		t.Fatalf("config: %v", err)
	}
	return New(cfg, slog.New(slog.NewTextHandler(io.Discard, nil)))
}

func get(t *testing.T, h http.Handler, path string) (*httptest.ResponseRecorder, map[string]any) {
	t.Helper()
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, path, nil))
	var body map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode %s: %v (body %q)", path, err, rec.Body.String())
	}
	return rec, body
}

func TestHealthz(t *testing.T) {
	s := newTestServer(t)
	rec, body := get(t, s.Handler(), "/healthz")
	if rec.Code != http.StatusOK || body["status"] != "ok" {
		t.Fatalf("got %d %v", rec.Code, body)
	}
	if ct := rec.Header().Get("Content-Type"); ct != "application/json; charset=utf-8" {
		t.Fatalf("content-type %q", ct)
	}
}

func TestReadyzNoCheckers(t *testing.T) {
	s := newTestServer(t)
	rec, body := get(t, s.Handler(), "/readyz")
	if rec.Code != http.StatusOK || body["status"] != "ok" {
		t.Fatalf("got %d %v", rec.Code, body)
	}
}

func TestReadyzFailingChecker(t *testing.T) {
	s := newTestServer(t)
	s.Health().Register(CheckerFunc{CheckName: "db", Fn: func(context.Context) error { return nil }})
	s.Health().Register(CheckerFunc{CheckName: "redis", Fn: func(context.Context) error { return errors.New("dial tcp: refused") }})

	rec, body := get(t, s.Handler(), "/readyz")
	if rec.Code != http.StatusServiceUnavailable || body["status"] != "unavailable" {
		t.Fatalf("got %d %v", rec.Code, body)
	}
	checks := body["checks"].(map[string]any)
	if checks["db"] != "ok" || checks["redis"] != "dial tcp: refused" {
		t.Fatalf("checks %v", checks)
	}
}

func TestNotFoundIsJSON(t *testing.T) {
	s := newTestServer(t)
	rec, body := get(t, s.Handler(), "/nope")
	if rec.Code != http.StatusNotFound {
		t.Fatalf("got %d", rec.Code)
	}
	if body["error"].(map[string]any)["message"] != "not found" {
		t.Fatalf("body %v", body)
	}
}

func TestRequestIDHeader(t *testing.T) {
	s := newTestServer(t)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	req.Header.Set("X-Request-Id", "abc-123")
	s.Handler().ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("got %d", rec.Code)
	}
}
