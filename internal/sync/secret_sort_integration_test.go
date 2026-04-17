package sync

import (
	"testing"
)

func TestApplySort_Integration_ConfigDriven(t *testing.T) {
	t.Setenv("VAULTPULL_SORT_ENABLED", "true")
	t.Setenv("VAULTPULL_SORT_FIELD", "key")
	t.Setenv("VAULTPULL_SORT_DIRECTION", "asc")

	cfg := SortConfigFromEnv()
	secrets := map[string]string{
		"DB_HOST":     "localhost",
		"API_KEY":     "secret",
		"APP_VERSION": "1.0",
	}

	keys := ApplySort(secrets, cfg)
	if len(keys) != 3 {
		t.Fatalf("expected 3 keys, got %d", len(keys))
	}
	if keys[0] != "API_KEY" {
		t.Errorf("expected API_KEY first, got %s", keys[0])
	}
	if keys[1] != "APP_VERSION" {
		t.Errorf("expected APP_VERSION second, got %s", keys[1])
	}
	if keys[2] != "DB_HOST" {
		t.Errorf("expected DB_HOST third, got %s", keys[2])
	}
}

func TestApplySort_Integration_DisabledFallback(t *testing.T) {
	t.Setenv("VAULTPULL_SORT_ENABLED", "false")

	cfg := SortConfigFromEnv()
	secrets := map[string]string{
		"Z_KEY": "last",
		"A_KEY": "first",
	}

	keys := ApplySort(secrets, cfg)
	if len(keys) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(keys))
	}
	if keys[0] != "A_KEY" {
		t.Errorf("expected A_KEY first even when disabled, got %s", keys[0])
	}
}
