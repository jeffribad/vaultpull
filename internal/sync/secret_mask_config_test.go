package sync

import (
	"testing"
)

func TestMaskConfigFromEnv_Defaults(t *testing.T) {
	t.Setenv("VAULTPULL_MASK_ENABLED", "")
	t.Setenv("VAULTPULL_MASK_KEYS", "")
	t.Setenv("VAULTPULL_MASK_REVEAL_LENGTH", "")

	cfg := MaskConfigFromEnv()
	if !cfg.Enabled {
		t.Error("expected Enabled=true by default")
	}
	if len(cfg.CustomKeys) != 0 {
		t.Errorf("expected no custom keys, got %v", cfg.CustomKeys)
	}
	if cfg.RevealLength != 4 {
		t.Errorf("expected RevealLength=4, got %d", cfg.RevealLength)
	}
}

func TestMaskConfigFromEnv_Disabled(t *testing.T) {
	t.Setenv("VAULTPULL_MASK_ENABLED", "false")
	cfg := MaskConfigFromEnv()
	if cfg.Enabled {
		t.Error("expected Enabled=false")
	}
}

func TestMaskConfigFromEnv_NumericEnabled(t *testing.T) {
	t.Setenv("VAULTPULL_MASK_ENABLED", "1")
	cfg := MaskConfigFromEnv()
	if !cfg.Enabled {
		t.Error("expected Enabled=true for value '1'")
	}
}

func TestMaskConfigFromEnv_CustomKeys(t *testing.T) {
	t.Setenv("VAULTPULL_MASK_KEYS", "my_token, api_secret , custom_key")
	cfg := MaskConfigFromEnv()
	if len(cfg.CustomKeys) != 3 {
		t.Fatalf("expected 3 custom keys, got %d", len(cfg.CustomKeys))
	}
	if cfg.CustomKeys[0] != "MY_TOKEN" {
		t.Errorf("expected MY_TOKEN, got %s", cfg.CustomKeys[0])
	}
}

func TestMaskConfigFromEnv_CustomRevealLength(t *testing.T) {
	t.Setenv("VAULTPULL_MASK_REVEAL_LENGTH", "6")
	cfg := MaskConfigFromEnv()
	if cfg.RevealLength != 6 {
		t.Errorf("expected RevealLength=6, got %d", cfg.RevealLength)
	}
}

func TestMaskConfigFromEnv_InvalidRevealLength_FallsBackToDefault(t *testing.T) {
	t.Setenv("VAULTPULL_MASK_REVEAL_LENGTH", "notanumber")
	cfg := MaskConfigFromEnv()
	if cfg.RevealLength != 4 {
		t.Errorf("expected fallback RevealLength=4, got %d", cfg.RevealLength)
	}
}
