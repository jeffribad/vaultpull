package sync

import (
	"os"
	"testing"
)

func TestDeprecateConfigFromEnv_Defaults(t *testing.T) {
	os.Unsetenv("VAULTPULL_DEPRECATE_ENABLED")
	os.Unsetenv("VAULTPULL_DEPRECATED_KEYS")
	os.Unsetenv("VAULTPULL_DEPRECATE_FAIL")

	cfg := DeprecateConfigFromEnv()
	if cfg.Enabled {
		t.Error("expected Enabled=false by default")
	}
	if cfg.FailOnUsage {
		t.Error("expected FailOnUsage=false by default")
	}
	if len(cfg.Deprecated) != 0 {
		t.Errorf("expected empty deprecated map, got %v", cfg.Deprecated)
	}
}

func TestDeprecateConfigFromEnv_Enabled(t *testing.T) {
	t.Setenv("VAULTPULL_DEPRECATE_ENABLED", "1")
	t.Setenv("VAULTPULL_DEPRECATE_FAIL", "true")

	cfg := DeprecateConfigFromEnv()
	if !cfg.Enabled {
		t.Error("expected Enabled=true")
	}
	if !cfg.FailOnUsage {
		t.Error("expected FailOnUsage=true")
	}
}

func TestDeprecateConfigFromEnv_ParsesKeys(t *testing.T) {
	t.Setenv("VAULTPULL_DEPRECATE_ENABLED", "1")
	t.Setenv("VAULTPULL_DEPRECATED_KEYS", "OLD_DB_URL:DATABASE_URL, LEGACY_TOKEN:")

	cfg := DeprecateConfigFromEnv()
	if len(cfg.Deprecated) != 2 {
		t.Fatalf("expected 2 deprecated keys, got %d", len(cfg.Deprecated))
	}
	if cfg.Deprecated["OLD_DB_URL"] != "DATABASE_URL" {
		t.Errorf("unexpected replacement: %q", cfg.Deprecated["OLD_DB_URL"])
	}
	if _, ok := cfg.Deprecated["LEGACY_TOKEN"]; !ok {
		t.Error("expected LEGACY_TOKEN to be present with empty replacement")
	}
}

func TestDeprecateConfigFromEnv_SkipsMalformed(t *testing.T) {
	t.Setenv("VAULTPULL_DEPRECATE_ENABLED", "1")
	t.Setenv("VAULTPULL_DEPRECATED_KEYS", "VALID_OLD:VALID_NEW,MALFORMED")

	cfg := DeprecateConfigFromEnv()
	if len(cfg.Deprecated) != 1 {
		t.Errorf("expected 1 valid key, got %d: %v", len(cfg.Deprecated), cfg.Deprecated)
	}
}
