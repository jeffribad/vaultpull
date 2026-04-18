package sync

import (
	"testing"
)

func TestTruncateConfigFromEnv_Defaults(t *testing.T) {
	t.Setenv("VAULTPULL_TRUNCATE_ENABLED", "")
	cfg := TruncateConfigFromEnv()
	if cfg.Enabled {
		t.Error("expected Enabled=false by default")
	}
	if cfg.MaxLength != 256 {
		t.Errorf("expected MaxLength=256, got %d", cfg.MaxLength)
	}
	if cfg.Suffix != "..." {
		t.Errorf("expected Suffix='...', got %q", cfg.Suffix)
	}
}

func TestTruncateConfigFromEnv_Enabled(t *testing.T) {
	t.Setenv("VAULTPULL_TRUNCATE_ENABLED", "true")
	t.Setenv("VAULTPULL_TRUNCATE_MAX_LENGTH", "64")
	t.Setenv("VAULTPULL_TRUNCATE_SUFFIX", "~~")
	cfg := TruncateConfigFromEnv()
	if !cfg.Enabled {
		t.Error("expected Enabled=true")
	}
	if cfg.MaxLength != 64 {
		t.Errorf("expected MaxLength=64, got %d", cfg.MaxLength)
	}
	if cfg.Suffix != "~~" {
		t.Errorf("expected Suffix='~~', got %q", cfg.Suffix)
	}
}

func TestTruncateSecrets_Disabled_ReturnsOriginal(t *testing.T) {
	secrets := map[string]string{"KEY": "a very long value that should not be truncated"}
	cfg := TruncateConfig{Enabled: false, MaxLength: 5, Suffix: "..."}
	out, truncated := TruncateSecrets(secrets, cfg)
	if out["KEY"] != secrets["KEY"] {
		t.Error("expected original value when disabled")
	}
	if len(truncated) != 0 {
		t.Error("expected no truncated keys when disabled")
	}
}

func TestTruncateSecrets_TruncatesLongValues(t *testing.T) {
	secrets := map[string]string{"TOKEN": "abcdefghij"}
	cfg := TruncateConfig{Enabled: true, MaxLength: 5, Suffix: "..."}
	out, truncated := TruncateSecrets(secrets, cfg)
	if out["TOKEN"] != "abcde..." {
		t.Errorf("expected 'abcde...', got %q", out["TOKEN"])
	}
	if len(truncated) != 1 {
		t.Errorf("expected 1 truncated key, got %d", len(truncated))
	}
}

func TestTruncateSecrets_ShortValues_Unchanged(t *testing.T) {
	secrets := map[string]string{"KEY": "hi"}
	cfg := TruncateConfig{Enabled: true, MaxLength: 10, Suffix: "..."}
	out, truncated := TruncateSecrets(secrets, cfg)
	if out["KEY"] != "hi" {
		t.Errorf("expected 'hi', got %q", out["KEY"])
	}
	if len(truncated) != 0 {
		t.Error("expected no truncated keys for short value")
	}
}

func TestTruncateSecrets_DoesNotMutateInput(t *testing.T) {
	secrets := map[string]string{"KEY": "abcdefghij"}
	cfg := TruncateConfig{Enabled: true, MaxLength: 3, Suffix: "!"}
	TruncateSecrets(secrets, cfg)
	if secrets["KEY"] != "abcdefghij" {
		t.Error("original map was mutated")
	}
}
