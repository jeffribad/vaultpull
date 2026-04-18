package sync

import (
	"testing"
)

func TestRedactConfigFromEnv_Defaults(t *testing.T) {
	t.Setenv("VAULTPULL_REDACT_ENABLED", "")
	t.Setenv("VAULTPULL_REDACT_KEYS", "")

	cfg := RedactConfigFromEnv()
	if cfg.Enabled {
		t.Error("expected Enabled to be false by default")
	}
	if len(cfg.Keys) != 0 {
		t.Errorf("expected no keys, got %v", cfg.Keys)
	}
}

func TestRedactConfigFromEnv_Enabled(t *testing.T) {
	t.Setenv("VAULTPULL_REDACT_ENABLED", "true")
	t.Setenv("VAULTPULL_REDACT_KEYS", "API_KEY, SECRET_TOKEN")

	cfg := RedactConfigFromEnv()
	if !cfg.Enabled {
		t.Error("expected Enabled to be true")
	}
	if len(cfg.Keys) != 2 {
		t.Errorf("expected 2 keys, got %d", len(cfg.Keys))
	}
}

func TestRedactConfigFromEnv_NumericEnabled(t *testing.T) {
	t.Setenv("VAULTPULL_REDACT_ENABLED", "1")
	t.Setenv("VAULTPULL_REDACT_KEYS", "")

	cfg := RedactConfigFromEnv()
	if !cfg.Enabled {
		t.Error("expected Enabled to be true for value '1'")
	}
}
