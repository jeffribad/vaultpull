package sync

import (
	"os"
	"testing"
)

func TestPrefixStripConfigFromEnv_Defaults(t *testing.T) {
	os.Unsetenv("VAULTPULL_PREFIX_STRIP_ENABLED")
	os.Unsetenv("VAULTPULL_PREFIX_STRIP_PREFIX")
	cfg := PrefixStripConfigFromEnv()
	if cfg.Enabled {
		t.Error("expected Enabled=false by default")
	}
	if cfg.Prefix != "" {
		t.Errorf("expected empty Prefix, got %q", cfg.Prefix)
	}
}

func TestPrefixStripConfigFromEnv_Enabled(t *testing.T) {
	t.Setenv("VAULTPULL_PREFIX_STRIP_ENABLED", "true")
	t.Setenv("VAULTPULL_PREFIX_STRIP_PREFIX", "APP_")
	cfg := PrefixStripConfigFromEnv()
	if !cfg.Enabled {
		t.Error("expected Enabled=true")
	}
	if cfg.Prefix != "APP_" {
		t.Errorf("expected Prefix=APP_, got %q", cfg.Prefix)
	}
}

func TestStripKeyPrefix_Disabled_ReturnsOriginal(t *testing.T) {
	cfg := PrefixStripConfig{Enabled: false, Prefix: "APP_"}
	input := map[string]string{"APP_FOO": "bar"}
	out := StripKeyPrefix(cfg, input)
	if _, ok := out["APP_FOO"]; !ok {
		t.Error("expected APP_FOO to remain unchanged")
	}
}

func TestStripKeyPrefix_StripsMatchingKeys(t *testing.T) {
	cfg := PrefixStripConfig{Enabled: true, Prefix: "APP_"}
	input := map[string]string{"APP_FOO": "1", "APP_BAR": "2", "OTHER": "3"}
	out := StripKeyPrefix(cfg, input)
	if _, ok := out["FOO"]; !ok {
		t.Error("expected FOO after stripping APP_")
	}
	if _, ok := out["BAR"]; !ok {
		t.Error("expected BAR after stripping APP_")
	}
	if _, ok := out["OTHER"]; !ok {
		t.Error("expected OTHER to remain")
	}
	if _, ok := out["APP_FOO"]; ok {
		t.Error("did not expect APP_FOO in output")
	}
}

func TestStripKeyPrefix_EmptyPrefix_ReturnsOriginal(t *testing.T) {
	cfg := PrefixStripConfig{Enabled: true, Prefix: ""}
	input := map[string]string{"FOO": "bar"}
	out := StripKeyPrefix(cfg, input)
	if out["FOO"] != "bar" {
		t.Error("expected original map returned for empty prefix")
	}
}

func TestStripKeyPrefix_DoesNotMutateInput(t *testing.T) {
	cfg := PrefixStripConfig{Enabled: true, Prefix: "X_"}
	input := map[string]string{"X_KEY": "val"}
	_ = StripKeyPrefix(cfg, input)
	if _, ok := input["X_KEY"]; !ok {
		t.Error("input map was mutated")
	}
}

func TestStripKeyPrefix_PrefixOnlyKey_KeepsOriginal(t *testing.T) {
	cfg := PrefixStripConfig{Enabled: true, Prefix: "APP_"}
	input := map[string]string{"APP_": "lonely"}
	out := StripKeyPrefix(cfg, input)
	if out["APP_"] != "lonely" {
		t.Error("expected key-only-prefix entry to be kept with original key")
	}
}
