package sync

import (
	"os"
	"testing"
)

func TestExpiryConfigFromEnv_Defaults(t *testing.T) {
	os.Unsetenv("VAULTPULL_EXPIRY_WARN_ENABLED")
	os.Unsetenv("VAULTPULL_EXPIRY_WARN_DAYS")

	cfg := ExpiryConfigFromEnv()
	if cfg.Enabled {
		t.Error("expected Enabled to be false by default")
	}
	if cfg.WarnWithinDays != 7 {
		t.Errorf("expected WarnWithinDays=7, got %d", cfg.WarnWithinDays)
	}
}

func TestExpiryConfigFromEnv_Enabled(t *testing.T) {
	t.Setenv("VAULTPULL_EXPIRY_WARN_ENABLED", "true")
	os.Unsetenv("VAULTPULL_EXPIRY_WARN_DAYS")

	cfg := ExpiryConfigFromEnv()
	if !cfg.Enabled {
		t.Error("expected Enabled to be true")
	}
}

func TestExpiryConfigFromEnv_NumericEnabled(t *testing.T) {
	t.Setenv("VAULTPULL_EXPIRY_WARN_ENABLED", "1")
	cfg := ExpiryConfigFromEnv()
	if !cfg.Enabled {
		t.Error("expected Enabled=true for value '1'")
	}
}

func TestExpiryConfigFromEnv_CustomDays(t *testing.T) {
	t.Setenv("VAULTPULL_EXPIRY_WARN_ENABLED", "true")
	t.Setenv("VAULTPULL_EXPIRY_WARN_DAYS", "14")

	cfg := ExpiryConfigFromEnv()
	if cfg.WarnWithinDays != 14 {
		t.Errorf("expected WarnWithinDays=14, got %d", cfg.WarnWithinDays)
	}
}

func TestExpiryConfigFromEnv_InvalidDays_FallsBackToDefault(t *testing.T) {
	t.Setenv("VAULTPULL_EXPIRY_WARN_ENABLED", "true")
	t.Setenv("VAULTPULL_EXPIRY_WARN_DAYS", "bad")

	cfg := ExpiryConfigFromEnv()
	if cfg.WarnWithinDays != 7 {
		t.Errorf("expected fallback WarnWithinDays=7, got %d", cfg.WarnWithinDays)
	}
}
