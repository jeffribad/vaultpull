package sync

import (
	"testing"
	"time"
)

func TestWatchConfigFromEnv_Defaults(t *testing.T) {
	t.Setenv("VAULTPULL_WATCH_ENABLED", "")
	t.Setenv("VAULTPULL_WATCH_INTERVAL_SECONDS", "")
	t.Setenv("VAULTPULL_WATCH_KEYS", "")

	cfg := WatchConfigFromEnv()
	if cfg.Enabled {
		t.Error("expected Enabled=false by default")
	}
	if cfg.Interval != 60*time.Second {
		t.Errorf("expected 60s interval, got %v", cfg.Interval)
	}
	if len(cfg.Keys) != 0 {
		t.Errorf("expected no keys, got %v", cfg.Keys)
	}
}

func TestWatchConfigFromEnv_Enabled(t *testing.T) {
	t.Setenv("VAULTPULL_WATCH_ENABLED", "true")
	cfg := WatchConfigFromEnv()
	if !cfg.Enabled {
		t.Error("expected Enabled=true")
	}
}

func TestWatchConfigFromEnv_NumericEnabled(t *testing.T) {
	t.Setenv("VAULTPULL_WATCH_ENABLED", "1")
	cfg := WatchConfigFromEnv()
	if !cfg.Enabled {
		t.Error("expected Enabled=true for value '1'")
	}
}

func TestWatchConfigFromEnv_CustomInterval(t *testing.T) {
	t.Setenv("VAULTPULL_WATCH_ENABLED", "true")
	t.Setenv("VAULTPULL_WATCH_INTERVAL_SECONDS", "30")
	cfg := WatchConfigFromEnv()
	if cfg.Interval != 30*time.Second {
		t.Errorf("expected 30s, got %v", cfg.Interval)
	}
}

func TestWatchConfigFromEnv_InvalidInterval_FallsBackToDefault(t *testing.T) {
	t.Setenv("VAULTPULL_WATCH_INTERVAL_SECONDS", "notanumber")
	cfg := WatchConfigFromEnv()
	if cfg.Interval != 60*time.Second {
		t.Errorf("expected fallback to 60s, got %v", cfg.Interval)
	}
}

func TestWatchConfigFromEnv_ParsesKeys(t *testing.T) {
	t.Setenv("VAULTPULL_WATCH_KEYS", "DB_PASS, API_KEY , SECRET")
	cfg := WatchConfigFromEnv()
	if len(cfg.Keys) != 3 {
		t.Fatalf("expected 3 keys, got %d", len(cfg.Keys))
	}
	if cfg.Keys[1] != "API_KEY" {
		t.Errorf("expected trimmed key, got %q", cfg.Keys[1])
	}
}
