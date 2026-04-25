package sync

import (
	"testing"
)

func TestPrefixAdd_Integration_ConfigDriven(t *testing.T) {
	t.Setenv("VAULTPULL_PREFIX_ADD_ENABLED", "true")
	t.Setenv("VAULTPULL_PREFIX_ADD_VALUE", "SVC_")
	t.Setenv("VAULTPULL_PREFIX_ADD_KEYS", "")

	cfg := PrefixAddConfigFromEnv()
	secrets := map[string]string{
		"DB_HOST": "db.internal",
		"DB_PORT": "5432",
		"API_KEY": "secret",
	}

	result := AddKeyPrefix(cfg, secrets)

	expected := []string{"SVC_DB_HOST", "SVC_DB_PORT", "SVC_API_KEY"}
	for _, key := range expected {
		if _, ok := result[key]; !ok {
			t.Errorf("expected key %q in result", key)
		}
	}
	if len(result) != 3 {
		t.Errorf("expected 3 keys, got %d", len(result))
	}
}

func TestPrefixAdd_Integration_DisabledPassthrough(t *testing.T) {
	t.Setenv("VAULTPULL_PREFIX_ADD_ENABLED", "false")
	t.Setenv("VAULTPULL_PREFIX_ADD_VALUE", "SVC_")

	cfg := PrefixAddConfigFromEnv()
	secrets := map[string]string{
		"HOST": "localhost",
	}

	result := AddKeyPrefix(cfg, secrets)

	if _, ok := result["HOST"]; !ok {
		t.Error("expected original key HOST to be present when disabled")
	}
	if _, ok := result["SVC_HOST"]; ok {
		t.Error("expected no prefixed key when disabled")
	}
}

func TestPrefixAdd_Integration_SubsetPipeline(t *testing.T) {
	t.Setenv("VAULTPULL_PREFIX_ADD_ENABLED", "true")
	t.Setenv("VAULTPULL_PREFIX_ADD_VALUE", "DB_")
	t.Setenv("VAULTPULL_PREFIX_ADD_KEYS", "HOST,PORT")

	cfg := PrefixAddConfigFromEnv()
	secrets := map[string]string{
		"HOST": "db.internal",
		"PORT": "5432",
		"TOKEN": "abc123",
	}

	result := AddKeyPrefix(cfg, secrets)

	if result["DB_HOST"] != "db.internal" {
		t.Errorf("expected DB_HOST=db.internal, got %q", result["DB_HOST"])
	}
	if result["DB_PORT"] != "5432" {
		t.Errorf("expected DB_PORT=5432, got %q", result["DB_PORT"])
	}
	if result["TOKEN"] != "abc123" {
		t.Errorf("expected TOKEN to remain unprefixed, got %q", result["TOKEN"])
	}
}
