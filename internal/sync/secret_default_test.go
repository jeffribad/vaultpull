package sync

import (
	"testing"
)

func TestDefaultConfigFromEnv_Defaults(t *testing.T) {
	t.Setenv("VAULTPULL_DEFAULTS_ENABLED", "")
	t.Setenv("VAULTPULL_DEFAULTS", "")

	cfg := DefaultConfigFromEnv()
	if cfg.Enabled {
		t.Error("expected Enabled=false")
	}
	if len(cfg.Defaults) != 0 {
		t.Errorf("expected empty defaults, got %v", cfg.Defaults)
	}
}

func TestDefaultConfigFromEnv_Enabled(t *testing.T) {
	t.Setenv("VAULTPULL_DEFAULTS_ENABLED", "true")
	t.Setenv("VAULTPULL_DEFAULTS", "REGION=us-east-1,LOG_LEVEL=info")

	cfg := DefaultConfigFromEnv()
	if !cfg.Enabled {
		t.Error("expected Enabled=true")
	}
	if cfg.Defaults["REGION"] != "us-east-1" {
		t.Errorf("expected REGION=us-east-1, got %q", cfg.Defaults["REGION"])
	}
	if cfg.Defaults["LOG_LEVEL"] != "info" {
		t.Errorf("expected LOG_LEVEL=info, got %q", cfg.Defaults["LOG_LEVEL"])
	}
}

func TestDefaultConfigFromEnv_NumericEnabled(t *testing.T) {
	t.Setenv("VAULTPULL_DEFAULTS_ENABLED", "1")
	t.Setenv("VAULTPULL_DEFAULTS", "")

	cfg := DefaultConfigFromEnv()
	if !cfg.Enabled {
		t.Error("expected Enabled=true for value '1'")
	}
}

func TestApplyDefaults_Disabled_ReturnsOriginal(t *testing.T) {
	cfg := DefaultConfig{Enabled: false, Defaults: map[string]string{"REGION": "us-east-1"}}
	secrets := map[string]string{"DB_HOST": "localhost"}

	result := ApplyDefaults(cfg, secrets)
	if _, ok := result["REGION"]; ok {
		t.Error("expected REGION to not be injected when disabled")
	}
}

func TestApplyDefaults_FillsMissingKeys(t *testing.T) {
	cfg := DefaultConfig{
		Enabled:  true,
		Defaults: map[string]string{"REGION": "us-east-1", "LOG_LEVEL": "info"},
	}
	secrets := map[string]string{"DB_HOST": "localhost"}

	result := ApplyDefaults(cfg, secrets)
	if result["REGION"] != "us-east-1" {
		t.Errorf("expected REGION=us-east-1, got %q", result["REGION"])
	}
	if result["LOG_LEVEL"] != "info" {
		t.Errorf("expected LOG_LEVEL=info, got %q", result["LOG_LEVEL"])
	}
	if result["DB_HOST"] != "localhost" {
		t.Errorf("expected DB_HOST=localhost, got %q", result["DB_HOST"])
	}
}

func TestApplyDefaults_DoesNotOverwriteExistingKeys(t *testing.T) {
	cfg := DefaultConfig{
		Enabled:  true,
		Defaults: map[string]string{"REGION": "eu-west-1"},
	}
	secrets := map[string]string{"REGION": "us-east-1"}

	result := ApplyDefaults(cfg, secrets)
	if result["REGION"] != "us-east-1" {
		t.Errorf("expected existing REGION to be preserved, got %q", result["REGION"])
	}
}

func TestApplyDefaults_DoesNotMutateInput(t *testing.T) {
	cfg := DefaultConfig{
		Enabled:  true,
		Defaults: map[string]string{"NEW_KEY": "value"},
	}
	secrets := map[string]string{"EXISTING": "yes"}

	_ = ApplyDefaults(cfg, secrets)
	if _, ok := secrets["NEW_KEY"]; ok {
		t.Error("ApplyDefaults must not mutate the input map")
	}
}
