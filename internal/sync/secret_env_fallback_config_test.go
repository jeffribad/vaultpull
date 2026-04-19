package sync

import (
	"testing"
)

func TestFallbackConfigFromEnv_Defaults(t *testing.T) {
	t.Setenv("VAULTPULL_FALLBACK_ENABLED", "")
	t.Setenv("VAULTPULL_FALLBACK_KEYS", "")
	t.Setenv("VAULTPULL_FALLBACK_PREFIX", "")
	cfg := FallbackConfigFromEnv()
	if cfg.Enabled {
		t.Error("expected Enabled=false by default")
	}
	if len(cfg.Keys) != 0 {
		t.Errorf("expected no keys, got %v", cfg.Keys)
	}
	if cfg.Prefix != "" {
		t.Errorf("expected empty prefix, got %s", cfg.Prefix)
	}
}

func TestFallbackConfigFromEnv_Enabled(t *testing.T) {
	t.Setenv("VAULTPULL_FALLBACK_ENABLED", "true")
	cfg := FallbackConfigFromEnv()
	if !cfg.Enabled {
		t.Error("expected Enabled=true")
	}
}

func TestFallbackConfigFromEnv_ParsesKeys(t *testing.T) {
	t.Setenv("VAULTPULL_FALLBACK_ENABLED", "1")
	t.Setenv("VAULTPULL_FALLBACK_KEYS", "DB_URL, API_KEY , TOKEN")
	cfg := FallbackConfigFromEnv()
	if len(cfg.Keys) != 3 {
		t.Fatalf("expected 3 keys, got %d", len(cfg.Keys))
	}
	if cfg.Keys[1] != "API_KEY" {
		t.Errorf("expected API_KEY, got %s", cfg.Keys[1])
	}
}

func TestFallbackConfigFromEnv_ParsesPrefix(t *testing.T) {
	t.Setenv("VAULTPULL_FALLBACK_PREFIX", "LOCAL_")
	cfg := FallbackConfigFromEnv()
	if cfg.Prefix != "LOCAL_" {
		t.Errorf("expected LOCAL_, got %s", cfg.Prefix)
	}
}
