package sync

import (
	"os"
	"testing"
)

func TestCastConfigFromEnv_Defaults(t *testing.T) {
	os.Unsetenv("VAULTPULL_CAST_ENABLED")
	os.Unsetenv("VAULTPULL_CAST_BOOL_KEYS")
	os.Unsetenv("VAULTPULL_CAST_INT_KEYS")
	cfg := CastConfigFromEnv()
	if cfg.Enabled {
		t.Error("expected Enabled=false by default")
	}
	if len(cfg.BoolKeys) != 0 {
		t.Errorf("expected empty BoolKeys, got %v", cfg.BoolKeys)
	}
	if len(cfg.IntKeys) != 0 {
		t.Errorf("expected empty IntKeys, got %v", cfg.IntKeys)
	}
}

func TestCastConfigFromEnv_Enabled(t *testing.T) {
	t.Setenv("VAULTPULL_CAST_ENABLED", "true")
	cfg := CastConfigFromEnv()
	if !cfg.Enabled {
		t.Error("expected Enabled=true")
	}
}

func TestCastConfigFromEnv_NumericEnabled(t *testing.T) {
	t.Setenv("VAULTPULL_CAST_ENABLED", "1")
	cfg := CastConfigFromEnv()
	if !cfg.Enabled {
		t.Error("expected Enabled=true for value '1'")
	}
}

func TestCastConfigFromEnv_ParsesKeys(t *testing.T) {
	t.Setenv("VAULTPULL_CAST_ENABLED", "true")
	t.Setenv("VAULTPULL_CAST_BOOL_KEYS", "FEATURE_FLAG, DEBUG")
	t.Setenv("VAULTPULL_CAST_INT_KEYS", "PORT , TIMEOUT")
	cfg := CastConfigFromEnv()
	if len(cfg.BoolKeys) != 2 || cfg.BoolKeys[0] != "FEATURE_FLAG" || cfg.BoolKeys[1] != "DEBUG" {
		t.Errorf("unexpected BoolKeys: %v", cfg.BoolKeys)
	}
	if len(cfg.IntKeys) != 2 || cfg.IntKeys[0] != "PORT" || cfg.IntKeys[1] != "TIMEOUT" {
		t.Errorf("unexpected IntKeys: %v", cfg.IntKeys)
	}
}
