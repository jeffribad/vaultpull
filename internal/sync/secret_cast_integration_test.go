package sync

import (
	"testing"
)

func TestCastSecrets_Integration_ConfigDriven(t *testing.T) {
	t.Setenv("VAULTPULL_CAST_ENABLED", "true")
	t.Setenv("VAULTPULL_CAST_BOOL_KEYS", "DEBUG,FEATURE_X")
	t.Setenv("VAULTPULL_CAST_INT_KEYS", "PORT,WORKERS")

	cfg := CastConfigFromEnv()
	secrets := map[string]string{
		"DEBUG":     "yes",
		"FEATURE_X": "0",
		"PORT":      "443xyz",
		"WORKERS":   "04",
		"OTHER":     "unchanged",
	}

	out := CastSecrets(secrets, cfg)

	if out["DEBUG"] != "true" {
		t.Errorf("DEBUG: expected \"true\", got %q", out["DEBUG"])
	}
	if out["FEATURE_X"] != "false" {
		t.Errorf("FEATURE_X: expected \"false\", got %q", out["FEATURE_X"])
	}
	if out["PORT"] != "443" {
		t.Errorf("PORT: expected \"443\", got %q", out["PORT"])
	}
	if out["WORKERS"] != "4" {
		t.Errorf("WORKERS: expected \"4\", got %q", out["WORKERS"])
	}
	if out["OTHER"] != "unchanged" {
		t.Errorf("OTHER: expected \"unchanged\", got %q", out["OTHER"])
	}
}

func TestCastSecrets_Integration_DisabledPassthrough(t *testing.T) {
	t.Setenv("VAULTPULL_CAST_ENABLED", "false")
	t.Setenv("VAULTPULL_CAST_BOOL_KEYS", "DEBUG")

	cfg := CastConfigFromEnv()
	secrets := map[string]string{"DEBUG": "yes"}
	out := CastSecrets(secrets, cfg)

	if out["DEBUG"] != "yes" {
		t.Errorf("expected passthrough \"yes\", got %q", out["DEBUG"])
	}
}
