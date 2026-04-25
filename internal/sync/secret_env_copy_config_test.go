package sync

import (
	"os"
	"testing"
)

func TestCopyConfigFromEnv_Defaults(t *testing.T) {
	os.Unsetenv("VAULTPULL_ENV_COPY_ENABLED")
	os.Unsetenv("VAULTPULL_ENV_COPY_PAIRS")

	cfg := CopyConfigFromEnv()
	if cfg.Enabled {
		t.Error("expected Enabled=false by default")
	}
	if len(cfg.Pairs) != 0 {
		t.Errorf("expected empty Pairs, got %v", cfg.Pairs)
	}
}

func TestCopyConfigFromEnv_Enabled(t *testing.T) {
	os.Setenv("VAULTPULL_ENV_COPY_ENABLED", "true")
	defer os.Unsetenv("VAULTPULL_ENV_COPY_ENABLED")

	cfg := CopyConfigFromEnv()
	if !cfg.Enabled {
		t.Error("expected Enabled=true")
	}
}

func TestCopyConfigFromEnv_NumericEnabled(t *testing.T) {
	os.Setenv("VAULTPULL_ENV_COPY_ENABLED", "1")
	defer os.Unsetenv("VAULTPULL_ENV_COPY_ENABLED")

	cfg := CopyConfigFromEnv()
	if !cfg.Enabled {
		t.Error("expected Enabled=true for value '1'")
	}
}

func TestCopyConfigFromEnv_ParsesPairs(t *testing.T) {
	os.Setenv("VAULTPULL_ENV_COPY_ENABLED", "true")
	os.Setenv("VAULTPULL_ENV_COPY_PAIRS", "DB_PASS:DATABASE_PASSWORD, API_KEY:SERVICE_API_KEY")
	defer os.Unsetenv("VAULTPULL_ENV_COPY_ENABLED")
	defer os.Unsetenv("VAULTPULL_ENV_COPY_PAIRS")

	cfg := CopyConfigFromEnv()
	if cfg.Pairs["DB_PASS"] != "DATABASE_PASSWORD" {
		t.Errorf("expected DB_PASS->DATABASE_PASSWORD, got %v", cfg.Pairs)
	}
	if cfg.Pairs["API_KEY"] != "SERVICE_API_KEY" {
		t.Errorf("expected API_KEY->SERVICE_API_KEY, got %v", cfg.Pairs)
	}
}

func TestCopyConfigFromEnv_SkipsMalformedPairs(t *testing.T) {
	os.Setenv("VAULTPULL_ENV_COPY_ENABLED", "true")
	os.Setenv("VAULTPULL_ENV_COPY_PAIRS", "BADENTRY, GOOD:KEY")
	defer os.Unsetenv("VAULTPULL_ENV_COPY_ENABLED")
	defer os.Unsetenv("VAULTPULL_ENV_COPY_PAIRS")

	cfg := CopyConfigFromEnv()
	if len(cfg.Pairs) != 1 {
		t.Errorf("expected 1 valid pair, got %d: %v", len(cfg.Pairs), cfg.Pairs)
	}
	if cfg.Pairs["GOOD"] != "KEY" {
		t.Errorf("expected GOOD->KEY, got %v", cfg.Pairs)
	}
}
