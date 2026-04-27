package sync

import (
	"os"
	"testing"
)

func TestSplitConfigFromEnv_Defaults(t *testing.T) {
	os.Unsetenv("VAULTPULL_SPLIT_ENABLED")
	os.Unsetenv("VAULTPULL_SPLIT_SOURCE")
	os.Unsetenv("VAULTPULL_SPLIT_DELIMITER")
	os.Unsetenv("VAULTPULL_SPLIT_SEPARATOR")

	cfg := SplitConfigFromEnv()
	if cfg.Enabled {
		t.Error("expected Enabled=false by default")
	}
	if cfg.SourceKey != "" {
		t.Errorf("expected empty SourceKey, got %q", cfg.SourceKey)
	}
	if cfg.Delimiter != "," {
		t.Errorf("expected default delimiter ',', got %q", cfg.Delimiter)
	}
	if cfg.Separator != "=" {
		t.Errorf("expected default separator '=', got %q", cfg.Separator)
	}
}

func TestSplitConfigFromEnv_Enabled(t *testing.T) {
	t.Setenv("VAULTPULL_SPLIT_ENABLED", "true")
	t.Setenv("VAULTPULL_SPLIT_SOURCE", "MY_MULTI")

	cfg := SplitConfigFromEnv()
	if !cfg.Enabled {
		t.Error("expected Enabled=true")
	}
	if cfg.SourceKey != "MY_MULTI" {
		t.Errorf("expected SourceKey=MY_MULTI, got %q", cfg.SourceKey)
	}
}

func TestSplitConfigFromEnv_CustomDelimiterAndSeparator(t *testing.T) {
	t.Setenv("VAULTPULL_SPLIT_DELIMITER", "|")
	t.Setenv("VAULTPULL_SPLIT_SEPARATOR", ":")

	cfg := SplitConfigFromEnv()
	if cfg.Delimiter != "|" {
		t.Errorf("expected delimiter '|', got %q", cfg.Delimiter)
	}
	if cfg.Separator != ":" {
		t.Errorf("expected separator ':', got %q", cfg.Separator)
	}
}

func TestSplitConfigFromEnv_NumericEnabled(t *testing.T) {
	t.Setenv("VAULTPULL_SPLIT_ENABLED", "1")

	cfg := SplitConfigFromEnv()
	if !cfg.Enabled {
		t.Error("expected Enabled=true for numeric '1'")
	}
}
