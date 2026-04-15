package sync

import (
	"testing"
	"time"
)

func TestRetryConfigFromEnv_Defaults(t *testing.T) {
	unsetRetryEnv(t)
	cfg := RetryConfigFromEnv()
	def := DefaultRetryConfig()
	if cfg.MaxAttempts != def.MaxAttempts {
		t.Errorf("MaxAttempts: got %d, want %d", cfg.MaxAttempts, def.MaxAttempts)
	}
	if cfg.Delay != def.Delay {
		t.Errorf("Delay: got %v, want %v", cfg.Delay, def.Delay)
	}
	if cfg.Multiplier != def.Multiplier {
		t.Errorf("Multiplier: got %f, want %f", cfg.Multiplier, def.Multiplier)
	}
}

func TestRetryConfigFromEnv_CustomValues(t *testing.T) {
	unsetRetryEnv(t)
	t.Setenv("VAULTPULL_RETRY_ATTEMPTS", "5")
	t.Setenv("VAULTPULL_RETRY_DELAY_MS", "200")
	t.Setenv("VAULTPULL_RETRY_MULTIPLIER", "1.5")

	cfg := RetryConfigFromEnv()
	if cfg.MaxAttempts != 5 {
		t.Errorf("MaxAttempts: got %d, want 5", cfg.MaxAttempts)
	}
	if cfg.Delay != 200*time.Millisecond {
		t.Errorf("Delay: got %v, want 200ms", cfg.Delay)
	}
	if cfg.Multiplier != 1.5 {
		t.Errorf("Multiplier: got %f, want 1.5", cfg.Multiplier)
	}
}

func TestRetryConfigFromEnv_InvalidValues_FallsBackToDefaults(t *testing.T) {
	unsetRetryEnv(t)
	t.Setenv("VAULTPULL_RETRY_ATTEMPTS", "notanint")
	t.Setenv("VAULTPULL_RETRY_DELAY_MS", "abc")
	t.Setenv("VAULTPULL_RETRY_MULTIPLIER", "bad")

	cfg := RetryConfigFromEnv()
	def := DefaultRetryConfig()
	if cfg.MaxAttempts != def.MaxAttempts {
		t.Errorf("MaxAttempts should fall back to default, got %d", cfg.MaxAttempts)
	}
	if cfg.Delay != def.Delay {
		t.Errorf("Delay should fall back to default, got %v", cfg.Delay)
	}
}

func TestRetryConfigFromEnv_MultiplierBelowOne_FallsBack(t *testing.T) {
	unsetRetryEnv(t)
	t.Setenv("VAULTPULL_RETRY_MULTIPLIER", "0.5")
	cfg := RetryConfigFromEnv()
	if cfg.Multiplier != DefaultRetryConfig().Multiplier {
		t.Errorf("expected default multiplier for value < 1.0, got %f", cfg.Multiplier)
	}
}

func unsetRetryEnv(t *testing.T) {
	t.Helper()
	for _, k := range []string{
		"VAULTPULL_RETRY_ATTEMPTS",
		"VAULTPULL_RETRY_DELAY_MS",
		"VAULTPULL_RETRY_MULTIPLIER",
	} {
		t.Setenv(k, "")
	}
}
