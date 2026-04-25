package sync

import (
	"os"
	"testing"
)

func TestExpandSecrets_Disabled_ReturnsOriginal(t *testing.T) {
	secrets := map[string]string{
		"KEY": "${OTHER}",
		"OTHER": "world",
	}
	cfg := ExpandConfig{Enabled: false}
	out, err := ExpandSecrets(secrets, cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["KEY"] != "${OTHER}" {
		t.Errorf("expected unexpanded value, got %q", out["KEY"])
	}
}

func TestExpandSecrets_ResolvesFromMap(t *testing.T) {
	secrets := map[string]string{
		"GREETING": "Hello, ${NAME}!",
		"NAME":     "World",
	}
	cfg := ExpandConfig{Enabled: true}
	out, err := ExpandSecrets(secrets, cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["GREETING"] != "Hello, World!" {
		t.Errorf("expected 'Hello, World!', got %q", out["GREETING"])
	}
}

func TestExpandSecrets_UnresolvedRef_EmptyString(t *testing.T) {
	secrets := map[string]string{
		"KEY": "prefix_${MISSING}_suffix",
	}
	cfg := ExpandConfig{Enabled: true, AllowEnv: false}
	out, err := ExpandSecrets(secrets, cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["KEY"] != "prefix__suffix" {
		t.Errorf("expected 'prefix__suffix', got %q", out["KEY"])
	}
}

func TestExpandSecrets_AllowEnv_ResolvesFromOS(t *testing.T) {
	t.Setenv("MY_OS_VAR", "from-os")
	secrets := map[string]string{
		"KEY": "${MY_OS_VAR}",
	}
	cfg := ExpandConfig{Enabled: true, AllowEnv: true}
	out, err := ExpandSecrets(secrets, cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["KEY"] != "from-os" {
		t.Errorf("expected 'from-os', got %q", out["KEY"])
	}
}

func TestExpandSecrets_DoesNotMutateInput(t *testing.T) {
	secrets := map[string]string{
		"A": "${B}",
		"B": "beta",
	}
	cfg := ExpandConfig{Enabled: true}
	_, _ = ExpandSecrets(secrets, cfg)
	if secrets["A"] != "${B}" {
		t.Error("input map was mutated")
	}
}

func TestExpandConfigFromEnv_Defaults(t *testing.T) {
	os.Unsetenv("VAULTPULL_EXPAND_ENABLED")
	os.Unsetenv("VAULTPULL_EXPAND_ALLOW_ENV")
	cfg := ExpandConfigFromEnv()
	if cfg.Enabled {
		t.Error("expected Enabled=false by default")
	}
	if cfg.AllowEnv {
		t.Error("expected AllowEnv=false by default")
	}
}

func TestExpandConfigFromEnv_Enabled(t *testing.T) {
	t.Setenv("VAULTPULL_EXPAND_ENABLED", "true")
	t.Setenv("VAULTPULL_EXPAND_ALLOW_ENV", "1")
	cfg := ExpandConfigFromEnv()
	if !cfg.Enabled {
		t.Error("expected Enabled=true")
	}
	if !cfg.AllowEnv {
		t.Error("expected AllowEnv=true")
	}
}
