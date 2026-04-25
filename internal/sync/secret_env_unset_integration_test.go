package sync_test

import (
	"os"
	"testing"

	"github.com/your-org/vaultpull/internal/sync"
)

func TestUnset_Integration_ConfigDriven(t *testing.T) {
	os.Setenv("VAULTPULL_UNSET_ENABLED", "true")
	os.Setenv("VAULTPULL_UNSET_KEYS", "LEGACY_TOKEN,OLD_SECRET")
	defer os.Unsetenv("VAULTPULL_UNSET_ENABLED")
	defer os.Unsetenv("VAULTPULL_UNSET_KEYS")

	secrets := map[string]string{
		"LEGACY_TOKEN": "abc",
		"OLD_SECRET":   "xyz",
		"ACTIVE_KEY":   "keep",
	}

	cfg := sync.UnsetConfigFromEnv()
	result := sync.ApplyUnset(cfg, secrets)

	if _, ok := result["LEGACY_TOKEN"]; ok {
		t.Error("LEGACY_TOKEN should have been removed")
	}
	if _, ok := result["OLD_SECRET"]; ok {
		t.Error("OLD_SECRET should have been removed")
	}
	if result["ACTIVE_KEY"] != "keep" {
		t.Error("ACTIVE_KEY should be preserved")
	}
}

func TestUnset_Integration_DisabledPassthrough(t *testing.T) {
	os.Unsetenv("VAULTPULL_UNSET_ENABLED")

	secrets := map[string]string{
		"LEGACY_TOKEN": "abc",
		"ACTIVE_KEY":   "keep",
	}

	cfg := sync.UnsetConfigFromEnv()
	result := sync.ApplyUnset(cfg, secrets)

	if len(result) != len(secrets) {
		t.Errorf("expected passthrough, got %d keys", len(result))
	}
}
