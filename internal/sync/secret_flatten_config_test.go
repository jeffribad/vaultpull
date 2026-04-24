package sync

import (
	"os"
	"testing"
)

func TestFlattenConfigFromEnv_Defaults(t *testing.T) {
	os.Unsetenv("VAULTPULL_FLATTEN_ENABLED")
	os.Unsetenv("VAULTPULL_FLATTEN_SEPARATOR")
	os.Unsetenv("VAULTPULL_FLATTEN_MAX_DEPTH")

	cfg := FlattenConfigFromEnv()
	if cfg.Enabled {
		t.Error("expected Enabled=false by default")
	}
	if cfg.Separator != "__" {
		t.Errorf("expected separator=__, got %q", cfg.Separator)
	}
	if cfg.MaxDepth != 5 {
		t.Errorf("expected MaxDepth=5, got %d", cfg.MaxDepth)
	}
}

func TestFlattenConfigFromEnv_Enabled(t *testing.T) {
	os.Setenv("VAULTPULL_FLATTEN_ENABLED", "true")
	defer os.Unsetenv("VAULTPULL_FLATTEN_ENABLED")

	cfg := FlattenConfigFromEnv()
	if !cfg.Enabled {
		t.Error("expected Enabled=true")
	}
}

func TestFlattenConfigFromEnv_NumericEnabled(t *testing.T) {
	os.Setenv("VAULTPULL_FLATTEN_ENABLED", "1")
	defer os.Unsetenv("VAULTPULL_FLATTEN_ENABLED")

	cfg := FlattenConfigFromEnv()
	if !cfg.Enabled {
		t.Error("expected Enabled=true for value '1'")
	}
}

func TestFlattenConfigFromEnv_CustomSeparator(t *testing.T) {
	os.Setenv("VAULTPULL_FLATTEN_SEPARATOR", "-")
	defer os.Unsetenv("VAULTPULL_FLATTEN_SEPARATOR")

	cfg := FlattenConfigFromEnv()
	if cfg.Separator != "-" {
		t.Errorf("expected separator=-, got %q", cfg.Separator)
	}
}

func TestFlattenConfigFromEnv_CustomMaxDepth(t *testing.T) {
	os.Setenv("VAULTPULL_FLATTEN_MAX_DEPTH", "3")
	defer os.Unsetenv("VAULTPULL_FLATTEN_MAX_DEPTH")

	cfg := FlattenConfigFromEnv()
	if cfg.MaxDepth != 3 {
		t.Errorf("expected MaxDepth=3, got %d", cfg.MaxDepth)
	}
}

func TestFlattenConfigFromEnv_InvalidMaxDepth_FallsBackToDefault(t *testing.T) {
	os.Setenv("VAULTPULL_FLATTEN_MAX_DEPTH", "notanumber")
	defer os.Unsetenv("VAULTPULL_FLATTEN_MAX_DEPTH")

	cfg := FlattenConfigFromEnv()
	if cfg.MaxDepth != 5 {
		t.Errorf("expected MaxDepth=5 on invalid input, got %d", cfg.MaxDepth)
	}
}
