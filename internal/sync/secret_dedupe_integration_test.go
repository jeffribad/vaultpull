package sync

import (
	"testing"
)

func TestDedupeSecrets_Integration_ConfigDriven(t *testing.T) {
	t.Setenv("VAULTPULL_DEDUPE_ENABLED", "true")
	t.Setenv("VAULTPULL_DEDUPE_CASE_SENSITIVE", "false")

	cfg := DedupeConfigFromEnv()
	secrets := map[string]string{
		"DB_HOST":  "localhost",
		"db_host":  "prodhost",
		"API_KEY":  "abc123",
		"api_key":  "xyz789",
		"UNIQUE":   "only",
	}

	out := DedupeSecrets(secrets, cfg)

	// Should collapse DB_HOST/db_host and API_KEY/api_key -> 3 unique keys
	if len(out) != 3 {
		t.Errorf("expected 3 keys after dedup, got %d: %v", len(out), out)
	}
}

func TestDedupeSecrets_Integration_DisabledPreservesAll(t *testing.T) {
	t.Setenv("VAULTPULL_DEDUPE_ENABLED", "false")

	cfg := DedupeConfigFromEnv()
	secrets := map[string]string{
		"DB_HOST": "localhost",
		"db_host": "prodhost",
	}

	out := DedupeSecrets(secrets, cfg)
	if len(out) != 2 {
		t.Errorf("expected 2 keys when disabled, got %d", len(out))
	}
}
