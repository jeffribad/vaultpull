package sync

import (
	"testing"
)

func TestEnvPromote_Integration_ConfigDriven(t *testing.T) {
	t.Setenv("VAULTPULL_PROMOTE_ENABLED", "true")
	t.Setenv("VAULTPULL_PROMOTE_FROM_PREFIX", "DEV_")
	t.Setenv("VAULTPULL_PROMOTE_TO_PREFIX", "STAGING_")
	t.Setenv("VAULTPULL_PROMOTE_OVERWRITE", "true")

	cfg := PromoteConfigFromEnv()
	secrets := map[string]string{
		"DEV_API_KEY": "abc123",
		"DEV_DB_URL":  "postgres://dev",
		"SHARED":      "common",
	}

	out := ApplyPromotion(secrets, cfg)

	if out["STAGING_API_KEY"] != "abc123" {
		t.Errorf("expected STAGING_API_KEY=abc123, got %q", out["STAGING_API_KEY"])
	}
	if out["STAGING_DB_URL"] != "postgres://dev" {
		t.Errorf("expected STAGING_DB_URL=postgres://dev, got %q", out["STAGING_DB_URL"])
	}
	if out["SHARED"] != "common" {
		t.Error("expected SHARED key to be preserved")
	}
	if out["DEV_API_KEY"] != "abc123" {
		t.Error("expected original DEV_ keys to be preserved")
	}
}

func TestEnvPromote_Integration_DisabledPassthrough(t *testing.T) {
	t.Setenv("VAULTPULL_PROMOTE_ENABLED", "false")
	t.Setenv("VAULTPULL_PROMOTE_FROM_PREFIX", "DEV_")
	t.Setenv("VAULTPULL_PROMOTE_TO_PREFIX", "PROD_")

	cfg := PromoteConfigFromEnv()
	secrets := map[string]string{"DEV_KEY": "value"}
	out := ApplyPromotion(secrets, cfg)

	if len(out) != 1 {
		t.Errorf("expected 1 key when disabled, got %d", len(out))
	}
	if _, ok := out["PROD_KEY"]; ok {
		t.Error("expected no promoted keys when disabled")
	}
}
