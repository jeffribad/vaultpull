package sync

import (
	"testing"
)

func TestValidateConfigFromEnv_Defaults(t *testing.T) {
	t.Setenv("VAULTPULL_VALIDATE", "")
	t.Setenv("VAULTPULL_REQUIRED_KEYS", "")
	t.Setenv("VAULTPULL_NONEMPTY_KEYS", "")

	cfg := ValidateConfigFromEnv()
	if cfg.Enabled {
		t.Error("expected Enabled=false by default")
	}
	if len(cfg.RequiredKeys) != 0 {
		t.Errorf("expected no required keys, got %v", cfg.RequiredKeys)
	}
}

func TestValidateConfigFromEnv_Enabled(t *testing.T) {
	t.Setenv("VAULTPULL_VALIDATE", "true")
	t.Setenv("VAULTPULL_REQUIRED_KEYS", "DB_HOST, DB_PORT")
	t.Setenv("VAULTPULL_NONEMPTY_KEYS", "API_KEY")

	cfg := ValidateConfigFromEnv()
	if !cfg.Enabled {
		t.Error("expected Enabled=true")
	}
	if len(cfg.RequiredKeys) != 2 {
		t.Errorf("expected 2 required keys, got %v", cfg.RequiredKeys)
	}
	if len(cfg.NonemptyKeys) != 1 || cfg.NonemptyKeys[0] != "API_KEY" {
		t.Errorf("unexpected nonempty keys: %v", cfg.NonemptyKeys)
	}
}

func TestValidateConfigFromEnv_NumericEnabled(t *testing.T) {
	t.Setenv("VAULTPULL_VALIDATE", "1")
	cfg := ValidateConfigFromEnv()
	if !cfg.Enabled {
		t.Error("expected Enabled=true for value '1'")
	}
}
