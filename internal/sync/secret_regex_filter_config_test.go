package sync

import (
	"os"
	"testing"
)

func TestRegexFilterConfigFromEnv_Defaults(t *testing.T) {
	os.Unsetenv("VAULTPULL_REGEX_FILTER_ENABLED")
	os.Unsetenv("VAULTPULL_REGEX_FILTER_ALLOW")
	os.Unsetenv("VAULTPULL_REGEX_FILTER_DENY")

	cfg := RegexFilterConfigFromEnv()
	if cfg.Enabled {
		t.Error("expected Enabled=false by default")
	}
	if cfg.AllowPattern != "" {
		t.Errorf("expected empty AllowPattern, got %q", cfg.AllowPattern)
	}
	if cfg.DenyPattern != "" {
		t.Errorf("expected empty DenyPattern, got %q", cfg.DenyPattern)
	}
}

func TestRegexFilterConfigFromEnv_Enabled(t *testing.T) {
	t.Setenv("VAULTPULL_REGEX_FILTER_ENABLED", "true")
	cfg := RegexFilterConfigFromEnv()
	if !cfg.Enabled {
		t.Error("expected Enabled=true")
	}
}

func TestRegexFilterConfigFromEnv_NumericEnabled(t *testing.T) {
	t.Setenv("VAULTPULL_REGEX_FILTER_ENABLED", "1")
	cfg := RegexFilterConfigFromEnv()
	if !cfg.Enabled {
		t.Error("expected Enabled=true for value '1'")
	}
}

func TestRegexFilterConfigFromEnv_ParsesPatterns(t *testing.T) {
	t.Setenv("VAULTPULL_REGEX_FILTER_ENABLED", "true")
	t.Setenv("VAULTPULL_REGEX_FILTER_ALLOW", "^DB_")
	t.Setenv("VAULTPULL_REGEX_FILTER_DENY", "_SECRET$")

	cfg := RegexFilterConfigFromEnv()
	if cfg.AllowPattern != "^DB_" {
		t.Errorf("expected AllowPattern=^DB_, got %q", cfg.AllowPattern)
	}
	if cfg.DenyPattern != "_SECRET$" {
		t.Errorf("expected DenyPattern=_SECRET$, got %q", cfg.DenyPattern)
	}
}
