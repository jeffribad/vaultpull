package sync

import (
	"os"
	"testing"
)

func TestPrefixAddConfigFromEnv_Defaults(t *testing.T) {
	os.Unsetenv("VAULTPULL_PREFIX_ADD_ENABLED")
	os.Unsetenv("VAULTPULL_PREFIX_ADD_VALUE")
	os.Unsetenv("VAULTPULL_PREFIX_ADD_KEYS")
	cfg := PrefixAddConfigFromEnv()
	if cfg.Enabled {
		t.Error("expected Enabled=false by default")
	}
	if cfg.Prefix != "" {
		t.Errorf("expected empty prefix, got %q", cfg.Prefix)
	}
	if len(cfg.Keys) != 0 {
		t.Errorf("expected no keys, got %v", cfg.Keys)
	}
}

func TestPrefixAddConfigFromEnv_Enabled(t *testing.T) {
	t.Setenv("VAULTPULL_PREFIX_ADD_ENABLED", "true")
	t.Setenv("VAULTPULL_PREFIX_ADD_VALUE", "APP_")
	cfg := PrefixAddConfigFromEnv()
	if !cfg.Enabled {
		t.Error("expected Enabled=true")
	}
	if cfg.Prefix != "APP_" {
		t.Errorf("expected prefix APP_, got %q", cfg.Prefix)
	}
}

func TestPrefixAddConfigFromEnv_NumericEnabled(t *testing.T) {
	t.Setenv("VAULTPULL_PREFIX_ADD_ENABLED", "1")
	cfg := PrefixAddConfigFromEnv()
	if !cfg.Enabled {
		t.Error("expected Enabled=true for value '1'")
	}
}

func TestPrefixAddConfigFromEnv_ParsesKeys(t *testing.T) {
	t.Setenv("VAULTPULL_PREFIX_ADD_ENABLED", "true")
	t.Setenv("VAULTPULL_PREFIX_ADD_VALUE", "SVC_")
	t.Setenv("VAULTPULL_PREFIX_ADD_KEYS", "DB_HOST, DB_PORT , DB_NAME")
	cfg := PrefixAddConfigFromEnv()
	if len(cfg.Keys) != 3 {
		t.Fatalf("expected 3 keys, got %d: %v", len(cfg.Keys), cfg.Keys)
	}
	if cfg.Keys[1] != "DB_PORT" {
		t.Errorf("expected DB_PORT, got %q", cfg.Keys[1])
	}
}
