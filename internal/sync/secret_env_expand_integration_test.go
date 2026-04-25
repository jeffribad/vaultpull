package sync_test

import (
	"testing"

	"github.com/yourusername/vaultpull/internal/sync"
)

func TestExpandSecrets_Integration_ConfigDriven(t *testing.T) {
	t.Setenv("VAULTPULL_EXPAND_ENABLED", "true")
	t.Setenv("VAULTPULL_EXPAND_ALLOW_ENV", "false")

	secrets := map[string]string{
		"BASE_URL":  "https://api.example.com",
		"FULL_URL":  "${BASE_URL}/v1/resource",
		"APP_NAME":  "myapp",
		"LOG_LABEL": "[${APP_NAME}] ready",
	}

	cfg := sync.ExpandConfigFromEnv()
	out, err := sync.ExpandSecrets(secrets, cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if out["FULL_URL"] != "https://api.example.com/v1/resource" {
		t.Errorf("FULL_URL: got %q", out["FULL_URL"])
	}
	if out["LOG_LABEL"] != "[myapp] ready" {
		t.Errorf("LOG_LABEL: got %q", out["LOG_LABEL"])
	}
}

func TestExpandSecrets_Integration_DisabledPassthrough(t *testing.T) {
	t.Setenv("VAULTPULL_EXPAND_ENABLED", "false")

	secrets := map[string]string{
		"KEY": "${SHOULD_NOT_EXPAND}",
	}

	cfg := sync.ExpandConfigFromEnv()
	out, err := sync.ExpandSecrets(secrets, cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["KEY"] != "${SHOULD_NOT_EXPAND}" {
		t.Errorf("expected raw value, got %q", out["KEY"])
	}
}
