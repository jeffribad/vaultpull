package sync

import (
	"crypto/sha256"
	"fmt"
	"testing"
)

func TestHashSecrets_Integration_ConfigDriven(t *testing.T) {
	t.Setenv("VAULTPULL_HASH_ENABLED", "true")
	t.Setenv("VAULTPULL_HASH_KEYS", "DB_PASSWORD,API_KEY")

	cfg := HashConfigFromEnv()
	secrets := map[string]string{
		"DB_PASSWORD": "s3cr3t",
		"API_KEY":     "myapikey",
		"APP_NAME":    "vaultpull",
	}

	out := HashSecrets(cfg, secrets)

	for _, key := range []string{"DB_PASSWORD", "API_KEY"} {
		original := secrets[key]
		sum := sha256.Sum256([]byte(original))
		expected := fmt.Sprintf("%x", sum)
		if out[key] != expected {
			t.Errorf("key %s: expected hash %q, got %q", key, expected, out[key])
		}
	}

	if out["APP_NAME"] != "vaultpull" {
		t.Errorf("APP_NAME should be unchanged, got %q", out["APP_NAME"])
	}
}

func TestHashSecrets_Integration_DisabledPassthrough(t *testing.T) {
	t.Setenv("VAULTPULL_HASH_ENABLED", "false")
	t.Setenv("VAULTPULL_HASH_KEYS", "DB_PASSWORD")

	cfg := HashConfigFromEnv()
	secrets := map[string]string{"DB_PASSWORD": "plaintext"}
	out := HashSecrets(cfg, secrets)

	if out["DB_PASSWORD"] != "plaintext" {
		t.Errorf("expected passthrough when disabled, got %q", out["DB_PASSWORD"])
	}
}
