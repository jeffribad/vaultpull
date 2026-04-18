package sync

import (
	"os"
	"testing"
)

func TestDedupeConfigFromEnv_Defaults(t *testing.T) {
	os.Unsetenv("VAULTPULL_DEDUPE_ENABLED")
	os.Unsetenv("VAULTPULL_DEDUPE_CASE_SENSITIVE")
	cfg := DedupeConfigFromEnv()
	if cfg.Enabled {
		t.Error("expected Enabled=false by default")
	}
	if cfg.CaseSensitive {
		t.Error("expected CaseSensitive=false by default")
	}
}

func TestDedupeConfigFromEnv_Enabled(t *testing.T) {
	t.Setenv("VAULTPULL_DEDUPE_ENABLED", "true")
	cfg := DedupeConfigFromEnv()
	if !cfg.Enabled {
		t.Error("expected Enabled=true")
	}
}

func TestDedupeConfigFromEnv_NumericEnabled(t *testing.T) {
	t.Setenv("VAULTPULL_DEDUPE_ENABLED", "1")
	cfg := DedupeConfigFromEnv()
	if !cfg.Enabled {
		t.Error("expected Enabled=true for value '1'")
	}
}

func TestDedupeConfigFromEnv_CaseSensitive(t *testing.T) {
	t.Setenv("VAULTPULL_DEDUPE_ENABLED", "true")
	t.Setenv("VAULTPULL_DEDUPE_CASE_SENSITIVE", "true")
	cfg := DedupeConfigFromEnv()
	if !cfg.CaseSensitive {
		t.Error("expected CaseSensitive=true")
	}
}
