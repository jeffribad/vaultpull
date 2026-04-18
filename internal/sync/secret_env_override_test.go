package sync

import (
	"testing"
)

func TestApplyEnvOverrides_Disabled_ReturnsOriginal(t *testing.T) {
	cfg := EnvOverrideConfig{Enabled: false, Prefix: "OVERRIDE_", Priority: "env"}
	secrets := map[string]string{"DB_HOST": "vault-host"}
	t.Setenv("OVERRIDE_DB_HOST", "local-host")

	result := ApplyEnvOverrides(cfg, secrets)
	if result["DB_HOST"] != "vault-host" {
		t.Errorf("expected vault-host, got %s", result["DB_HOST"])
	}
}

func TestApplyEnvOverrides_EnvPriority_OverridesVault(t *testing.T) {
	cfg := EnvOverrideConfig{Enabled: true, Prefix: "OVERRIDE_", Priority: "env"}
	secrets := map[string]string{"DB_HOST": "vault-host"}
	t.Setenv("OVERRIDE_DB_HOST", "local-host")

	result := ApplyEnvOverrides(cfg, secrets)
	if result["DB_HOST"] != "local-host" {
		t.Errorf("expected local-host, got %s", result["DB_HOST"])
	}
}

func TestApplyEnvOverrides_VaultPriority_KeepsVault(t *testing.T) {
	cfg := EnvOverrideConfig{Enabled: true, Prefix: "OVERRIDE_", Priority: "vault"}
	secrets := map[string]string{"DB_HOST": "vault-host"}
	t.Setenv("OVERRIDE_DB_HOST", "local-host")

	result := ApplyEnvOverrides(cfg, secrets)
	if result["DB_HOST"] != "vault-host" {
		t.Errorf("expected vault-host, got %s", result["DB_HOST"])
	}
}

func TestApplyEnvOverrides_AddsNewKey(t *testing.T) {
	cfg := EnvOverrideConfig{Enabled: true, Prefix: "OVERRIDE_", Priority: "env"}
	secrets := map[string]string{"DB_HOST": "vault-host"}
	t.Setenv("OVERRIDE_NEW_KEY", "new-value")

	result := ApplyEnvOverrides(cfg, secrets)
	if result["NEW_KEY"] != "new-value" {
		t.Errorf("expected new-value, got %s", result["NEW_KEY"])
	}
}

func TestApplyEnvOverrides_DoesNotMutateInput(t *testing.T) {
	cfg := EnvOverrideConfig{Enabled: true, Prefix: "OVERRIDE_", Priority: "env"}
	secrets := map[string]string{"KEY": "original"}
	t.Setenv("OVERRIDE_KEY", "changed")

	ApplyEnvOverrides(cfg, secrets)
	if secrets["KEY"] != "original" {
		t.Error("input map was mutated")
	}
}

func TestEnvOverrideConfigFromEnv_Defaults(t *testing.T) {
	t.Setenv("VAULTPULL_ENV_OVERRIDE_ENABLED", "")
	t.Setenv("VAULTPULL_ENV_OVERRIDE_PREFIX", "")
	t.Setenv("VAULTPULL_ENV_OVERRIDE_PRIORITY", "")

	cfg := EnvOverrideConfigFromEnv()
	if cfg.Enabled {
		t.Error("expected disabled by default")
	}
	if cfg.Prefix != "VAULTPULL_OVERRIDE_" {
		t.Errorf("unexpected default prefix: %s", cfg.Prefix)
	}
	if cfg.Priority != "env" {
		t.Errorf("unexpected default priority: %s", cfg.Priority)
	}
}

func TestEnvOverrideConfigFromEnv_Enabled(t *testing.T) {
	t.Setenv("VAULTPULL_ENV_OVERRIDE_ENABLED", "true")
	cfg := EnvOverrideConfigFromEnv()
	if !cfg.Enabled {
		t.Error("expected enabled")
	}
}
