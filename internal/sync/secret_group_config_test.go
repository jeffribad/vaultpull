package sync

import (
	"os"
	"testing"
)

func TestGroupConfigFromEnv_Defaults(t *testing.T) {
	os.Unsetenv("VAULTPULL_GROUP_ENABLED")
	os.Unsetenv("VAULTPULL_GROUP_KEY")
	os.Unsetenv("VAULTPULL_GROUP_OUT_DIR")

	cfg := GroupConfigFromEnv()
	if cfg.Enabled {
		t.Error("expected Enabled=false by default")
	}
	if cfg.GroupKey != "group" {
		t.Errorf("expected GroupKey=group, got %s", cfg.GroupKey)
	}
	if cfg.OutDir != "." {
		t.Errorf("expected OutDir=., got %s", cfg.OutDir)
	}
}

func TestGroupConfigFromEnv_Enabled(t *testing.T) {
	t.Setenv("VAULTPULL_GROUP_ENABLED", "true")
	cfg := GroupConfigFromEnv()
	if !cfg.Enabled {
		t.Error("expected Enabled=true")
	}
}

func TestGroupConfigFromEnv_NumericEnabled(t *testing.T) {
	t.Setenv("VAULTPULL_GROUP_ENABLED", "1")
	cfg := GroupConfigFromEnv()
	if !cfg.Enabled {
		t.Error("expected Enabled=true for value '1'")
	}
}

func TestGroupConfigFromEnv_CustomValues(t *testing.T) {
	t.Setenv("VAULTPULL_GROUP_KEY", "team")
	t.Setenv("VAULTPULL_GROUP_OUT_DIR", "/tmp/envs")
	cfg := GroupConfigFromEnv()
	if cfg.GroupKey != "team" {
		t.Errorf("expected GroupKey=team, got %s", cfg.GroupKey)
	}
	if cfg.OutDir != "/tmp/envs" {
		t.Errorf("expected OutDir=/tmp/envs, got %s", cfg.OutDir)
	}
}
