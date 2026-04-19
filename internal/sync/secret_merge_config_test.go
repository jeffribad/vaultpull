package sync

import (
	"os"
	"testing"
)

func TestMergeConfigFromEnv_Defaults(t *testing.T) {
	os.Unsetenv("VAULTPULL_MERGE_ENABLED")
	os.Unsetenv("VAULTPULL_MERGE_STRATEGY")
	cfg := MergeConfigFromEnv()
	if cfg.Enabled {
		t.Error("expected disabled by default")
	}
	if cfg.Strategy != "vault-wins" {
		t.Errorf("expected vault-wins default, got %s", cfg.Strategy)
	}
}

func TestMergeConfigFromEnv_Enabled(t *testing.T) {
	os.Setenv("VAULTPULL_MERGE_ENABLED", "true")
	defer os.Unsetenv("VAULTPULL_MERGE_ENABLED")
	cfg := MergeConfigFromEnv()
	if !cfg.Enabled {
		t.Error("expected enabled")
	}
}

func TestMergeConfigFromEnv_CustomStrategy(t *testing.T) {
	os.Setenv("VAULTPULL_MERGE_STRATEGY", "local-wins")
	defer os.Unsetenv("VAULTPULL_MERGE_STRATEGY")
	cfg := MergeConfigFromEnv()
	if cfg.Strategy != "local-wins" {
		t.Errorf("expected local-wins, got %s", cfg.Strategy)
	}
}

func TestMergeConfigFromEnv_NumericEnabled(t *testing.T) {
	os.Setenv("VAULTPULL_MERGE_ENABLED", "1")
	defer os.Unsetenv("VAULTPULL_MERGE_ENABLED")
	cfg := MergeConfigFromEnv()
	if !cfg.Enabled {
		t.Error("expected enabled via numeric 1")
	}
}
