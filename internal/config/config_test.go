package config

import (
	"log/slog"
	"strings"
	"testing"
	"time"
)

func env(m map[string]string) func(string) (string, bool) {
	return func(k string) (string, bool) { v, ok := m[k]; return v, ok }
}

func TestLoadDefaults(t *testing.T) {
	cfg, err := Load(env(nil))
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if cfg.Env != "dev" || cfg.Addr != ":8080" || cfg.LogLevel != slog.LevelInfo || cfg.ShutdownTimeout != 10*time.Second {
		t.Fatalf("unexpected defaults: %+v", cfg)
	}
	if !cfg.IsDev() {
		t.Fatal("expected IsDev")
	}
}

func TestLoadOverrides(t *testing.T) {
	cfg, err := Load(env(map[string]string{
		"AEGIS_ENV":              "PROD",
		"AEGIS_ADDR":             "127.0.0.1:9000",
		"DATABASE_URL":           "postgres://x",
		"REDIS_URL":              "redis://y",
		"LOG_LEVEL":              "debug",
		"AEGIS_SHUTDOWN_TIMEOUT": "1m",
	}))
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if cfg.Env != "prod" || cfg.Addr != "127.0.0.1:9000" || cfg.DatabaseURL != "postgres://x" ||
		cfg.RedisURL != "redis://y" || cfg.LogLevel != slog.LevelDebug || cfg.ShutdownTimeout != time.Minute {
		t.Fatalf("overrides not applied: %+v", cfg)
	}
}

func TestLoadErrors(t *testing.T) {
	tests := []struct {
		name string
		env  map[string]string
		want string
	}{
		{"bad env", map[string]string{"AEGIS_ENV": "staging"}, "AEGIS_ENV"},
		{"bad level", map[string]string{"LOG_LEVEL": "loud"}, "LOG_LEVEL"},
		{"bad timeout", map[string]string{"AEGIS_SHUTDOWN_TIMEOUT": "soon"}, "AEGIS_SHUTDOWN_TIMEOUT"},
		{"zero timeout", map[string]string{"AEGIS_SHUTDOWN_TIMEOUT": "0s"}, "must be positive"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := Load(env(tc.env))
			if err == nil || !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("want error containing %q, got %v", tc.want, err)
			}
		})
	}
}
