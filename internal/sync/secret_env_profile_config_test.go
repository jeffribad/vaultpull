package sync

import (
	"os"
	"testing"
)

func TestEnvProfileConfigFromEnv_Defaults(t *testing.T) {
	os.Unsetenv("VAULTPULL_PROFILE_ENABLED")
	os.Unsetenv("VAULTPULL_PROFILE")
	os.Unsetenv("VAULTPULL_PROFILES")

	cfg := EnvProfileConfigFromEnv()
	if cfg.Enabled {
		t.Error("expected Enabled=false by default")
	}
	if cfg.Profile != "" {
		t.Errorf("expected empty Profile, got %q", cfg.Profile)
	}
	if len(cfg.Profiles) != 0 {
		t.Errorf("expected empty Profiles, got %v", cfg.Profiles)
	}
}

func TestEnvProfileConfigFromEnv_Enabled(t *testing.T) {
	t.Setenv("VAULTPULL_PROFILE_ENABLED", "true")
	t.Setenv("VAULTPULL_PROFILE", "prod")
	t.Setenv("VAULTPULL_PROFILES", "dev,staging,prod")

	cfg := EnvProfileConfigFromEnv()
	if !cfg.Enabled {
		t.Error("expected Enabled=true")
	}
	if cfg.Profile != "prod" {
		t.Errorf("expected Profile=prod, got %q", cfg.Profile)
	}
	if len(cfg.Profiles) != 3 {
		t.Errorf("expected 3 profiles, got %d", len(cfg.Profiles))
	}
}

func TestEnvProfileConfigFromEnv_NumericEnabled(t *testing.T) {
	t.Setenv("VAULTPULL_PROFILE_ENABLED", "1")
	t.Setenv("VAULTPULL_PROFILE", "staging")
	os.Unsetenv("VAULTPULL_PROFILES")

	cfg := EnvProfileConfigFromEnv()
	if !cfg.Enabled {
		t.Error("expected Enabled=true with value '1'")
	}
}
