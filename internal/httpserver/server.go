// Package httpserver wires the chi router, middleware and lifecycle for the API.
package httpserver

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/rj-karan/AegisFlow/internal/config"
)

// Server is the HTTP server for the control plane API.
type Server struct {
	cfg    config.Config
	log    *slog.Logger
	http   *http.Server
	router chi.Router
	health *Health
}

// New builds a Server with routes registered. Dependency checks (database,
// queue) are registered on the returned Server's Health before Run.
func New(cfg config.Config, log *slog.Logger) *Server {
	s := &Server{
		cfg:    cfg,
		log:    log,
		router: chi.NewRouter(),
		health: NewHealth(),
	}
	s.routes()
	s.http = &http.Server{
		Addr:              cfg.Addr,
		Handler:           s.router,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
		MaxHeaderBytes:    1 << 16,
	}
	return s
}

// Health exposes the readiness registry so main can add dependency checks.
func (s *Server) Health() *Health { return s.health }

// Handler returns the root handler, mainly for tests.
func (s *Server) Handler() http.Handler { return s.router }

func (s *Server) routes() {
	r := s.router
	r.Use(middleware.RequestID)
	r.Use(requestLogger(s.log))
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(30 * time.Second))

	r.Get("/healthz", s.health.Live)
	r.Get("/readyz", s.health.Ready)

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		writeError(w, http.StatusNotFound, "not found")
	})
	r.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
	})
}

// Run serves until ctx is cancelled, then shuts down gracefully within the
// configured timeout.
func (s *Server) Run(ctx context.Context) error {
	errc := make(chan error, 1)
	go func() {
		s.log.Info("http server listening", "addr", s.cfg.Addr, "env", s.cfg.Env)
		if err := s.http.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errc <- fmt.Errorf("listen on %s: %w", s.cfg.Addr, err)
			return
		}
		errc <- nil
	}()

	select {
	case err := <-errc:
		return err
	case <-ctx.Done():
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), s.cfg.ShutdownTimeout)
	defer cancel()
	s.log.Info("http server shutting down", "timeout", s.cfg.ShutdownTimeout)
	if err := s.http.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("shutdown: %w", err)
	}
	return <-errc
}
