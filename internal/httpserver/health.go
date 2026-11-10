package httpserver

import (
	"context"
	"net/http"
	"sync"
	"time"
)

// Checker reports whether a dependency is usable. Implementations must be
// cheap: readiness is polled by orchestrators.
type Checker interface {
	Name() string
	Check(ctx context.Context) error
}

// CheckerFunc adapts a function to Checker.
type CheckerFunc struct {
	CheckName string
	Fn        func(ctx context.Context) error
}

// Name implements Checker.
func (c CheckerFunc) Name() string { return c.CheckName }

// Check implements Checker.
func (c CheckerFunc) Check(ctx context.Context) error { return c.Fn(ctx) }

// Health serves liveness and readiness.
type Health struct {
	mu       sync.RWMutex
	checkers []Checker
	timeout  time.Duration
}

// NewHealth returns an empty registry: /readyz is ready until a checker fails.
func NewHealth() *Health { return &Health{timeout: 2 * time.Second} }

// Register adds a dependency check.
func (h *Health) Register(c Checker) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.checkers = append(h.checkers, c)
}

// Live always answers ok: the process is up and can serve HTTP.
func (h *Health) Live(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// Ready runs every registered check and answers 503 if any fails.
func (h *Health) Ready(w http.ResponseWriter, r *http.Request) {
	h.mu.RLock()
	checkers := append([]Checker(nil), h.checkers...)
	h.mu.RUnlock()

	ctx, cancel := context.WithTimeout(r.Context(), h.timeout)
	defer cancel()

	status := http.StatusOK
	checks := make(map[string]string, len(checkers))
	for _, c := range checkers {
		if err := c.Check(ctx); err != nil {
			status = http.StatusServiceUnavailable
			checks[c.Name()] = err.Error()
			continue
		}
		checks[c.Name()] = "ok"
	}

	body := map[string]any{"status": "ok", "checks": checks}
	if status != http.StatusOK {
		body["status"] = "unavailable"
	}
	writeJSON(w, status, body)
}
