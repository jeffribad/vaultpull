package sync

import (
	"testing"
	"time"
)

func TestSecretExpiry_Integration_ConfigDriven(t *testing.T) {
	t.Setenv("VAULTPULL_EXPIRY_WARN_ENABLED", "true")
	t.Setenv("VAULTPULL_EXPIRY_WARN_DAYS", "3")

	cfg := ExpiryConfigFromEnv()

	secrets := map[string]string{
		"CERT_EXPIRES":  time.Now().Add(1 * 24 * time.Hour).Format(time.RFC3339),
		"TOKEN_EXPIRES": time.Now().Add(10 * 24 * time.Hour).Format(time.RFC3339),
		"DB_PASSWORD":   "not-a-date",
	}

	warnings := CheckSecretExpiry(cfg, secrets)
	if len(warnings) != 1 {
		t.Fatalf("expected 1 warning, got %d", len(warnings))
	}
	if warnings[0].Key != "CERT_EXPIRES" {
		t.Errorf("expected warning for CERT_EXPIRES, got %s", warnings[0].Key)
	}
}

func TestSecretExpiry_Integration_DisabledSkipsAll(t *testing.T) {
	t.Setenv("VAULTPULL_EXPIRY_WARN_ENABLED", "false")

	cfg := ExpiryConfigFromEnv()

	secrets := map[string]string{
		"EXPIRED_KEY": time.Now().Add(-24 * time.Hour).Format(time.RFC3339),
	}

	warnings := CheckSecretExpiry(cfg, secrets)
	if len(warnings) != 0 {
		t.Errorf("expected no warnings when disabled, got %d", len(warnings))
	}
}
