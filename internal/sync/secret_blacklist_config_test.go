package sync

import (
	"os"
	"testing"
)

func TestBlacklistConfigFromEnv_Defaults(t *testing.T) {
	os.Unsetenv("VAULTPULL_BLACKLIST_ENABLED")
	os.Unsetenv("VAULTPULL_BLACKLIST_KEYS")
	cfg := BlacklistConfigFromEnv()
	if cfg.Enabled {
		t.Error("expected Enabled=false by default")
	}
	if len(cfg.Keys) != 0 {
		t.Errorf("expected no keys by default, got %v", cfg.Keys)
	}
}

func TestBlacklistConfigFromEnv_Enabled(t *testing.T) {
	t.Setenv("VAULTPULL_BLACKLIST_ENABLED", "true")
	t.Setenv("VAULTPULL_BLACKLIST_KEYS", "SECRET,TOKEN")
	cfg := BlacklistConfigFromEnv()
	if !cfg.Enabled {
		t.Error("expected Enabled=true")
	}
	if len(cfg.Keys) != 2 {
		t.Errorf("expected 2 keys, got %d", len(cfg.Keys))
	}
}

func TestBlacklistConfigFromEnv_NumericEnabled(t *testing.T) {
	t.Setenv("VAULTPULL_BLACKLIST_ENABLED", "1")
	cfg := BlacklistConfigFromEnv()
	if !cfg.Enabled {
		t.Error("expected Enabled=true for value '1'")
	}
}

func TestBlacklistConfigFromEnv_ParsesKeys(t *testing.T) {
	t.Setenv("VAULTPULL_BLACKLIST_ENABLED", "true")
	t.Setenv("VAULTPULL_BLACKLIST_KEYS", " DEBUG_TOKEN , INTERNAL_SECRET , ADMIN_PASS ")
	cfg := BlacklistConfigFromEnv()
	if len(cfg.Keys) != 3 {
		t.Fatalf("expected 3 keys, got %d: %v", len(cfg.Keys), cfg.Keys)
	}
	if cfg.Keys[0] != "DEBUG_TOKEN" {
		t.Errorf("expected trimmed key DEBUG_TOKEN, got %q", cfg.Keys[0])
	}
}
