// Package config loads runtime configuration from the environment.
//
// Every setting has a documented default so a fresh checkout runs with no
// .env file. Secrets are read here once and never logged.
package config

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"
)

// Config is the fully resolved runtime configuration for a process.
type Config struct {
	// Env is "dev", "test" or "prod". It changes logging format and cookie flags.
	Env string
	// Addr is the listen address for the HTTP server.
	Addr string
	// DatabaseURL is a pgx connection string. Empty means "not configured yet".
	DatabaseURL string
	// RedisURL is a redis:// URL. Empty means "not configured yet".
	RedisURL string
	// LogLevel is one of debug, info, warn, error.
	LogLevel slog.Level
	// ShutdownTimeout bounds graceful shutdown.
	ShutdownTimeout time.Duration
}

// Load reads configuration from the given lookup function (usually os.LookupEnv).
func Load(lookup func(string) (string, bool)) (Config, error) {
	get := func(key, def string) string {
		if v, ok := lookup(key); ok && v != "" {
			return v
		}
		return def
	}

	cfg := Config{
		Env:         strings.ToLower(get("AEGIS_ENV", "dev")),
		Addr:        get("AEGIS_ADDR", ":8080"),
		DatabaseURL: get("DATABASE_URL", ""),
		RedisURL:    get("REDIS_URL", ""),
	}

	level, err := parseLevel(get("LOG_LEVEL", "info"))
	if err != nil {
		return Config{}, err
	}
	cfg.LogLevel = level

	timeout, err := time.ParseDuration(get("AEGIS_SHUTDOWN_TIMEOUT", "10s"))
	if err != nil {
		return Config{}, fmt.Errorf("AEGIS_SHUTDOWN_TIMEOUT: %w", err)
	}
	cfg.ShutdownTimeout = timeout

	if err := cfg.validate(); err != nil {
		return Config{}, err
	}
	return cfg, nil
}

// FromEnv is Load with the process environment.
func FromEnv() (Config, error) { return Load(os.LookupEnv) }

// IsDev reports whether the process runs in the development environment.
func (c Config) IsDev() bool { return c.Env == "dev" }

func (c Config) validate() error {
	var errs []error
	switch c.Env {
	case "dev", "test", "prod":
	default:
		errs = append(errs, fmt.Errorf("AEGIS_ENV: unknown environment %q (want dev, test or prod)", c.Env))
	}
	if c.Addr == "" {
		errs = append(errs, errors.New("AEGIS_ADDR: must not be empty"))
	}
	if c.ShutdownTimeout <= 0 {
		errs = append(errs, errors.New("AEGIS_SHUTDOWN_TIMEOUT: must be positive"))
	}
	return errors.Join(errs...)
}

func parseLevel(s string) (slog.Level, error) {
	switch strings.ToLower(s) {
	case "debug":
		return slog.LevelDebug, nil
	case "info":
		return slog.LevelInfo, nil
	case "warn", "warning":
		return slog.LevelWarn, nil
	case "error":
		return slog.LevelError, nil
	}
	return 0, fmt.Errorf("LOG_LEVEL: unknown level %q", s)
}
