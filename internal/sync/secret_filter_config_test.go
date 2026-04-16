package sync

import (
	"testing"
)

func TestSecretFilterConfigFromEnv_Defaults(t *testing.T) {
	t.Setenv("VAULTPULL_INCLUDE_KEYS", "")
	t.Setenv("VAULTPULL_EXCLUDE_KEYS", "")

	cfg := SecretFilterConfigFromEnv()
	if len(cfg.IncludeKeys) != 0 {
		t.Errorf("expected no include keys, got %v", cfg.IncludeKeys)
	}
	if len(cfg.ExcludeKeys) != 0 {
		t.Errorf("expected no exclude keys, got %v", cfg.ExcludeKeys)
	}
}

func TestSecretFilterConfigFromEnv_ParsesKeys(t *testing.T) {
	t.Setenv("VAULTPULL_INCLUDE_KEYS", "DB_HOST, DB_PORT , API_KEY")
	t.Setenv("VAULTPULL_EXCLUDE_KEYS", "SECRET_TOKEN")

	cfg := SecretFilterConfigFromEnv()
	if len(cfg.IncludeKeys) != 3 {
		t.Fatalf("expected 3 include keys, got %d", len(cfg.IncludeKeys))
	}
	if cfg.IncludeKeys[1] != "DB_PORT" {
		t.Errorf("expected DB_PORT, got %s", cfg.IncludeKeys[1])
	}
	if len(cfg.ExcludeKeys) != 1 || cfg.ExcludeKeys[0] != "SECRET_TOKEN" {
		t.Errorf("unexpected exclude keys: %v", cfg.ExcludeKeys)
	}
}
