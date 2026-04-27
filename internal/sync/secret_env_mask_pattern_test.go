package sync

import (
	"testing"
)

func TestMaskPatternConfigFromEnv_Defaults(t *testing.T) {
	cfg := MaskPatternConfigFromEnv()
	if cfg.Enabled {
		t.Error("expected Enabled=false by default")
	}
	if cfg.Mask != "***" {
		t.Errorf("expected default mask '***', got %q", cfg.Mask)
	}
	if cfg.Pattern != "" {
		t.Errorf("expected empty pattern by default, got %q", cfg.Pattern)
	}
}

func TestMaskPatternConfigFromEnv_Enabled(t *testing.T) {
	t.Setenv("VAULTPULL_MASK_PATTERN_ENABLED", "true")
	t.Setenv("VAULTPULL_MASK_PATTERN_REGEX", "secret")
	t.Setenv("VAULTPULL_MASK_PATTERN_MASK", "[REDACTED]")

	cfg := MaskPatternConfigFromEnv()
	if !cfg.Enabled {
		t.Error("expected Enabled=true")
	}
	if cfg.Pattern != "secret" {
		t.Errorf("unexpected pattern: %q", cfg.Pattern)
	}
	if cfg.Mask != "[REDACTED]" {
		t.Errorf("unexpected mask: %q", cfg.Mask)
	}
	if cfg.CompiledRe == nil {
		t.Error("expected compiled regex to be set")
	}
}

func TestApplyMaskPattern_Disabled_ReturnsOriginal(t *testing.T) {
	cfg := MaskPatternConfig{Enabled: false}
	secrets := map[string]string{"SECRET_KEY": "sensitive"}
	out := ApplyMaskPattern(cfg, secrets)
	if out["SECRET_KEY"] != "sensitive" {
		t.Errorf("expected original value, got %q", out["SECRET_KEY"])
	}
}

func TestApplyMaskPattern_MasksMatchingKeys(t *testing.T) {
	t.Setenv("VAULTPULL_MASK_PATTERN_ENABLED", "1")
	t.Setenv("VAULTPULL_MASK_PATTERN_REGEX", "password|secret")

	cfg := MaskPatternConfigFromEnv()
	secrets := map[string]string{
		"DB_PASSWORD": "hunter2",
		"API_SECRET":  "abc123",
		"APP_NAME":    "vaultpull",
	}

	out := ApplyMaskPattern(cfg, secrets)

	if out["DB_PASSWORD"] != "***" {
		t.Errorf("DB_PASSWORD should be masked, got %q", out["DB_PASSWORD"])
	}
	if out["API_SECRET"] != "***" {
		t.Errorf("API_SECRET should be masked, got %q", out["API_SECRET"])
	}
	if out["APP_NAME"] != "vaultpull" {
		t.Errorf("APP_NAME should be unchanged, got %q", out["APP_NAME"])
	}
}

func TestApplyMaskPattern_DoesNotMutateInput(t *testing.T) {
	t.Setenv("VAULTPULL_MASK_PATTERN_ENABLED", "true")
	t.Setenv("VAULTPULL_MASK_PATTERN_REGEX", "token")

	cfg := MaskPatternConfigFromEnv()
	original := map[string]string{"API_TOKEN": "secret-value"}
	ApplyMaskPattern(cfg, original)

	if original["API_TOKEN"] != "secret-value" {
		t.Error("input map was mutated")
	}
}

func TestApplyMaskPattern_InvalidRegex_ReturnsOriginal(t *testing.T) {
	cfg := MaskPatternConfig{
		Enabled:    true,
		Pattern:    "[invalid",
		Mask:       "***",
		CompiledRe: nil, // invalid regex won't compile
	}
	secrets := map[string]string{"KEY": "value"}
	out := ApplyMaskPattern(cfg, secrets)
	if out["KEY"] != "value" {
		t.Errorf("expected original value when regex is nil, got %q", out["KEY"])
	}
}
