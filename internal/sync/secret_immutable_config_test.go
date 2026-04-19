package sync

import (
	"os"
	"testing"
)

func TestImmutableConfigFromEnv_Defaults(t *testing.T) {
	os.Unsetenv("VAULTPULL_IMMUTABLE_ENABLED")
	os.Unsetenv("VAULTPULL_IMMUTABLE_KEYS")

	cfg := ImmutableConfigFromEnv()
	if cfg.Enabled {
		t.Error("expected Enabled=false by default")
	}
	if len(cfg.Keys) != 0 {
		t.Errorf("expected no keys, got %v", cfg.Keys)
	}
}

func TestImmutableConfigFromEnv_Enabled(t *testing.T) {
	os.Setenv("VAULTPULL_IMMUTABLE_ENABLED", "true")
	defer os.Unsetenv("VAULTPULL_IMMUTABLE_ENABLED")

	cfg := ImmutableConfigFromEnv()
	if !cfg.Enabled {
		t.Error("expected Enabled=true")
	}
}

func TestImmutableConfigFromEnv_NumericEnabled(t *testing.T) {
	os.Setenv("VAULTPULL_IMMUTABLE_ENABLED", "1")
	defer os.Unsetenv("VAULTPULL_IMMUTABLE_ENABLED")

	cfg := ImmutableConfigFromEnv()
	if !cfg.Enabled {
		t.Error("expected Enabled=true for value '1'")
	}
}

func TestImmutableConfigFromEnv_ParsesKeys(t *testing.T) {
	os.Setenv("VAULTPULL_IMMUTABLE_ENABLED", "true")
	os.Setenv("VAULTPULL_IMMUTABLE_KEYS", "DB_PASS, API_KEY , SECRET")
	defer os.Unsetenv("VAULTPULL_IMMUTABLE_ENABLED")
	defer os.Unsetenv("VAULTPULL_IMMUTABLE_KEYS")

	cfg := ImmutableConfigFromEnv()
	if len(cfg.Keys) != 3 {
		t.Fatalf("expected 3 keys, got %d", len(cfg.Keys))
	}
	if cfg.Keys[1] != "API_KEY" {
		t.Errorf("expected trimmed key 'API_KEY', got %q", cfg.Keys[1])
	}
}
