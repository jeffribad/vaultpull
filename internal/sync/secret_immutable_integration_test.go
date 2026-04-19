package sync

import (
	"os"
	"testing"
)

func TestImmutable_Integration_ConfigDriven(t *testing.T) {
	os.Setenv("VAULTPULL_IMMUTABLE_ENABLED", "true")
	os.Setenv("VAULTPULL_IMMUTABLE_KEYS", "ROOT_TOKEN")
	defer os.Unsetenv("VAULTPULL_IMMUTABLE_ENABLED")
	defer os.Unsetenv("VAULTPULL_IMMUTABLE_KEYS")

	cfg := ImmutableConfigFromEnv()

	current := map[string]string{"ROOT_TOKEN": "original"}
	incoming := map[string]string{"ROOT_TOKEN": "replaced"}

	if err := EnforceImmutable(cfg, current, incoming); err == nil {
		t.Fatal("expected immutable violation error")
	}
}

func TestImmutable_Integration_DisabledPassthrough(t *testing.T) {
	os.Setenv("VAULTPULL_IMMUTABLE_ENABLED", "false")
	os.Setenv("VAULTPULL_IMMUTABLE_KEYS", "ROOT_TOKEN")
	defer os.Unsetenv("VAULTPULL_IMMUTABLE_ENABLED")
	defer os.Unsetenv("VAULTPULL_IMMUTABLE_KEYS")

	cfg := ImmutableConfigFromEnv()

	current := map[string]string{"ROOT_TOKEN": "original"}
	incoming := map[string]string{"ROOT_TOKEN": "replaced"}

	if err := EnforceImmutable(cfg, current, incoming); err != nil {
		t.Fatalf("expected nil when disabled, got %v", err)
	}
}
