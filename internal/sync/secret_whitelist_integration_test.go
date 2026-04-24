package sync

import (
	"testing"
)

func TestWhitelist_Integration_ConfigDriven(t *testing.T) {
	t.Setenv("VAULTPULL_WHITELIST_ENABLED", "true")
	t.Setenv("VAULTPULL_WHITELIST_KEYS", "DB_HOST,DB_PORT")

	cfg := WhitelistConfigFromEnv()

	secrets := map[string]string{
		"DB_HOST":      "localhost",
		"DB_PORT":      "5432",
		"API_KEY":      "supersecret",
		"ADMIN_TOKEN":  "tok123",
	}

	result := ApplyWhitelist(cfg, secrets)

	if len(result) != 2 {
		t.Errorf("expected 2 keys after whitelist, got %d", len(result))
	}
	for _, key := range []string{"DB_HOST", "DB_PORT"} {
		if _, ok := result[key]; !ok {
			t.Errorf("expected key %q to be present", key)
		}
	}
	for _, key := range []string{"API_KEY", "ADMIN_TOKEN"} {
		if _, ok := result[key]; ok {
			t.Errorf("expected key %q to be removed", key)
		}
	}
}

func TestWhitelist_Integration_DisabledPassthrough(t *testing.T) {
	t.Setenv("VAULTPULL_WHITELIST_ENABLED", "false")
	t.Setenv("VAULTPULL_WHITELIST_KEYS", "DB_HOST")

	cfg := WhitelistConfigFromEnv()

	secrets := map[string]string{
		"DB_HOST": "localhost",
		"API_KEY": "supersecret",
	}

	result := ApplyWhitelist(cfg, secrets)

	if len(result) != 2 {
		t.Errorf("expected all 2 keys to pass through when disabled, got %d", len(result))
	}
}
