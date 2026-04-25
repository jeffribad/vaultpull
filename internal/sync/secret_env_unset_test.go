package sync

import (
	"os"
	"testing"
)

func TestApplyUnset_Disabled_ReturnsOriginal(t *testing.T) {
	cfg := UnsetConfig{Enabled: false, Keys: []string{"OLD_KEY"}}
	secrets := map[string]string{"OLD_KEY": "v", "KEEP": "yes"}
	got := ApplyUnset(cfg, secrets)
	if len(got) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(got))
	}
}

func TestApplyUnset_RemovesExplicitKeys(t *testing.T) {
	cfg := UnsetConfig{Enabled: true, Keys: []string{"OLD_KEY", "DEPRECATED"}}
	secrets := map[string]string{"OLD_KEY": "v", "DEPRECATED": "d", "KEEP": "yes"}
	got := ApplyUnset(cfg, secrets)
	if _, ok := got["OLD_KEY"]; ok {
		t.Error("expected OLD_KEY to be removed")
	}
	if _, ok := got["DEPRECATED"]; ok {
		t.Error("expected DEPRECATED to be removed")
	}
	if got["KEEP"] != "yes" {
		t.Error("expected KEEP to be preserved")
	}
}

func TestApplyUnset_NoKeys_ReturnsOriginal(t *testing.T) {
	cfg := UnsetConfig{Enabled: true, Keys: nil}
	secrets := map[string]string{"A": "1", "B": "2"}
	got := ApplyUnset(cfg, secrets)
	if len(got) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(got))
	}
}

func TestApplyUnset_DoesNotMutateInput(t *testing.T) {
	cfg := UnsetConfig{Enabled: true, Keys: []string{"X"}}
	secrets := map[string]string{"X": "val", "Y": "keep"}
	_ = ApplyUnset(cfg, secrets)
	if _, ok := secrets["X"]; !ok {
		t.Error("ApplyUnset mutated the input map")
	}
}

func TestApplyUnset_SyncWithEnv_RemovesEnvOnlyKeys(t *testing.T) {
	// Set an env var that is NOT in vault secrets.
	os.Setenv("STALE_ENV_KEY", "stale")
	defer os.Unsetenv("STALE_ENV_KEY")

	cfg := UnsetConfig{Enabled: true, SyncWithEnv: true}
	// STALE_ENV_KEY is not in secrets, so it should not appear in output.
	secrets := map[string]string{"VAULT_KEY": "v"}
	got := ApplyUnset(cfg, secrets)
	if _, ok := got["STALE_ENV_KEY"]; ok {
		t.Error("expected STALE_ENV_KEY to be absent from output")
	}
	if got["VAULT_KEY"] != "v" {
		t.Error("expected VAULT_KEY to be preserved")
	}
}

func TestUnsetConfigFromEnv_Defaults(t *testing.T) {
	os.Unsetenv("VAULTPULL_UNSET_ENABLED")
	os.Unsetenv("VAULTPULL_UNSET_KEYS")
	os.Unsetenv("VAULTPULL_UNSET_SYNC_WITH_ENV")
	cfg := UnsetConfigFromEnv()
	if cfg.Enabled {
		t.Error("expected Enabled=false by default")
	}
	if len(cfg.Keys) != 0 {
		t.Error("expected no keys by default")
	}
	if cfg.SyncWithEnv {
		t.Error("expected SyncWithEnv=false by default")
	}
}

func TestUnsetConfigFromEnv_ParsesKeys(t *testing.T) {
	os.Setenv("VAULTPULL_UNSET_ENABLED", "true")
	os.Setenv("VAULTPULL_UNSET_KEYS", "FOO, BAR")
	defer os.Unsetenv("VAULTPULL_UNSET_ENABLED")
	defer os.Unsetenv("VAULTPULL_UNSET_KEYS")
	cfg := UnsetConfigFromEnv()
	if !cfg.Enabled {
		t.Error("expected Enabled=true")
	}
	if len(cfg.Keys) != 2 || cfg.Keys[0] != "FOO" || cfg.Keys[1] != "BAR" {
		t.Errorf("unexpected keys: %v", cfg.Keys)
	}
}
