package sync

import (
	"testing"
)

func TestEnvProfile_Integration_ConfigDriven(t *testing.T) {
	t.Setenv("VAULTPULL_PROFILE_ENABLED", "true")
	t.Setenv("VAULTPULL_PROFILE", "staging")
	t.Setenv("VAULTPULL_PROFILES", "dev,staging,prod")

	cfg := EnvProfileConfigFromEnv()
	secrets := map[string]string{
		"API_KEY__staging": "stg-secret",
		"API_KEY__prod":    "prod-secret",
		"LOG_LEVEL":        "info",
	}

	out := ApplyEnvProfile(cfg, secrets)

	if v, ok := out["API_KEY"]; !ok || v != "stg-secret" {
		t.Errorf("expected API_KEY=stg-secret, got %v", out)
	}
	if _, ok := out["API_KEY__prod"]; ok {
		t.Error("prod key should be excluded in staging profile")
	}
	if _, ok := out["LOG_LEVEL"]; !ok {
		t.Error("non-profiled key should be included")
	}
}

func TestEnvProfile_Integration_DisabledPassthrough(t *testing.T) {
	t.Setenv("VAULTPULL_PROFILE_ENABLED", "false")
	t.Setenv("VAULTPULL_PROFILE", "prod")
	t.Setenv("VAULTPULL_PROFILES", "dev,prod")

	cfg := EnvProfileConfigFromEnv()
	secrets := map[string]string{
		"KEY__prod": "v1",
		"KEY__dev":  "v2",
	}

	out := ApplyEnvProfile(cfg, secrets)
	if len(out) != len(secrets) {
		t.Errorf("expected all %d keys passthrough, got %d", len(secrets), len(out))
	}
}
