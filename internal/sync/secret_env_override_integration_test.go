package sync

import (
	"testing"
)

func TestEnvOverride_Integration_ConfigDriven(t *testing.T) {
	t.Setenv("VAULTPULL_ENV_OVERRIDE_ENABLED", "1")
	t.Setenv("VAULTPULL_ENV_OVERRIDE_PREFIX", "LOCAL_")
	t.Setenv("VAULTPULL_ENV_OVERRIDE_PRIORITY", "env")
	t.Setenv("LOCAL_API_KEY", "local-api-key")

	cfg := EnvOverrideConfigFromEnv()
	secrets := map[string]string{
		"API_KEY": "vault-api-key",
		"DB_PASS": "vault-db-pass",
	}

	result := ApplyEnvOverrides(cfg, secrets)

	if result["API_KEY"] != "local-api-key" {
		t.Errorf("expected local-api-key, got %s", result["API_KEY"])
	}
	if result["DB_PASS"] != "vault-db-pass" {
		t.Errorf("expected vault-db-pass, got %s", result["DB_PASS"])
	}
}

func TestEnvOverride_Integration_DisabledPassthrough(t *testing.T) {
	t.Setenv("VAULTPULL_ENV_OVERRIDE_ENABLED", "false")
	t.Setenv("VAULTPULL_ENV_OVERRIDE_PREFIX", "LOCAL_")
	t.Setenv("LOCAL_API_KEY", "local-api-key")

	cfg := EnvOverrideConfigFromEnv()
	secrets := map[string]string{
		"API_KEY": "vault-api-key",
	}

	result := ApplyEnvOverrides(cfg, secrets)

	if result["API_KEY"] != "vault-api-key" {
		t.Errorf("expected vault-api-key, got %s", result["API_KEY"])
	}
}
