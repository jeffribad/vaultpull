package sync

import (
	"testing"
)

func TestSuffixStripConfigFromEnv_Defaults(t *testing.T) {
	cfg := SuffixStripConfigFromEnv()
	if cfg.Enabled {
		t.Error("expected Enabled=false by default")
	}
	if cfg.Suffix != "" {
		t.Errorf("expected empty Suffix, got %q", cfg.Suffix)
	}
}

func TestSuffixStripConfigFromEnv_Enabled(t *testing.T) {
	t.Setenv("VAULTPULL_SUFFIX_STRIP_ENABLED", "true")
	t.Setenv("VAULTPULL_SUFFIX_STRIP_SUFFIX", "_SECRET")
	cfg := SuffixStripConfigFromEnv()
	if !cfg.Enabled {
		t.Error("expected Enabled=true")
	}
	if cfg.Suffix != "_SECRET" {
		t.Errorf("expected Suffix=_SECRET, got %q", cfg.Suffix)
	}
}

func TestStripKeySuffix_Disabled_ReturnsOriginal(t *testing.T) {
	cfg := SuffixStripConfig{Enabled: false, Suffix: "_SECRET"}
	input := map[string]string{"DB_PASSWORD_SECRET": "abc"}
	out := StripKeySuffix(cfg, input)
	if _, ok := out["DB_PASSWORD_SECRET"]; !ok {
		t.Error("expected original key to be preserved when disabled")
	}
}

func TestStripKeySuffix_EmptySuffix_ReturnsOriginal(t *testing.T) {
	cfg := SuffixStripConfig{Enabled: true, Suffix: ""}
	input := map[string]string{"DB_PASSWORD": "abc"}
	out := StripKeySuffix(cfg, input)
	if _, ok := out["DB_PASSWORD"]; !ok {
		t.Error("expected original key when suffix is empty")
	}
}

func TestStripKeySuffix_StripsMatchingKeys(t *testing.T) {
	cfg := SuffixStripConfig{Enabled: true, Suffix: "_SECRET"}
	input := map[string]string{
		"DB_PASSWORD_SECRET": "hunter2",
		"API_KEY_SECRET":     "xyz",
		"PLAIN_KEY":         "plain",
	}
	out := StripKeySuffix(cfg, input)

	if v, ok := out["DB_PASSWORD"]; !ok || v != "hunter2" {
		t.Errorf("expected DB_PASSWORD=hunter2, got %v (present=%v)", v, ok)
	}
	if v, ok := out["API_KEY"]; !ok || v != "xyz" {
		t.Errorf("expected API_KEY=xyz, got %v (present=%v)", v, ok)
	}
	if v, ok := out["PLAIN_KEY"]; !ok || v != "plain" {
		t.Errorf("expected PLAIN_KEY=plain, got %v (present=%v)", v, ok)
	}
}

func TestStripKeySuffix_EmptyResultKey_KeepsOriginal(t *testing.T) {
	cfg := SuffixStripConfig{Enabled: true, Suffix: "_SECRET"}
	input := map[string]string{"_SECRET": "val"}
	out := StripKeySuffix(cfg, input)
	if _, ok := out["_SECRET"]; !ok {
		t.Error("expected original key to be kept when stripping would produce empty key")
	}
}

func TestStripKeySuffix_DoesNotMutateInput(t *testing.T) {
	cfg := SuffixStripConfig{Enabled: true, Suffix: "_SECRET"}
	input := map[string]string{"TOKEN_SECRET": "abc"}
	_ = StripKeySuffix(cfg, input)
	if _, ok := input["TOKEN_SECRET"]; !ok {
		t.Error("original map was mutated")
	}
}
