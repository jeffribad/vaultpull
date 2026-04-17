package sync

import (
	"os"
	"testing"
)

func TestTTLConfigFromEnv_Defaults(t *testing.T) {
	os.Unsetenv("VAULTPULL_TTL_ENABLED")
	os.Unsetenv("VAULTPULL_TTL_MAX_AGE_DAYS")
	cfg := TTLConfigFromEnv()
	if cfg.Enabled {
		t.Error("expected Enabled=false by default")
	}
	if cfg.MaxAgeDays != 30 {
		t.Errorf("expected MaxAgeDays=30, got %d", cfg.MaxAgeDays)
	}
}

func TestTTLConfigFromEnv_Enabled(t *testing.T) {
	t.Setenv("VAULTPULL_TTL_ENABLED", "true")
	cfg := TTLConfigFromEnv()
	if !cfg.Enabled {
		t.Error("expected Enabled=true")
	}
}

func TestTTLConfigFromEnv_NumericEnabled(t *testing.T) {
	t.Setenv("VAULTPULL_TTL_ENABLED", "1")
	cfg := TTLConfigFromEnv()
	if !cfg.Enabled {
		t.Error("expected Enabled=true for value '1'")
	}
}

func TestTTLConfigFromEnv_CustomMaxAge(t *testing.T) {
	t.Setenv("VAULTPULL_TTL_MAX_AGE_DAYS", "90")
	cfg := TTLConfigFromEnv()
	if cfg.MaxAgeDays != 90 {
		t.Errorf("expected MaxAgeDays=90, got %d", cfg.MaxAgeDays)
	}
}

func TestTTLConfigFromEnv_InvalidMaxAge_FallsBackToDefault(t *testing.T) {
	t.Setenv("VAULTPULL_TTL_MAX_AGE_DAYS", "notanumber")
	cfg := TTLConfigFromEnv()
	if cfg.MaxAgeDays != 30 {
		t.Errorf("expected fallback MaxAgeDays=30, got %d", cfg.MaxAgeDays)
	}
}
