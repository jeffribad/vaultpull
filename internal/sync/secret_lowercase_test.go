package sync

import (
	"testing"
)

func TestLowercaseConfigFromEnv_Defaults(t *testing.T) {
	t.Setenv("VAULTPULL_LOWERCASE_KEYS", "")
	cfg := LowercaseConfigFromEnv()
	if cfg.Enabled {
		t.Error("expected Enabled to be false by default")
	}
}

func TestLowercaseConfigFromEnv_Enabled(t *testing.T) {
	t.Setenv("VAULTPULL_LOWERCASE_KEYS", "true")
	cfg := LowercaseConfigFromEnv()
	if !cfg.Enabled {
		t.Error("expected Enabled to be true")
	}
}

func TestLowercaseConfigFromEnv_NumericEnabled(t *testing.T) {
	t.Setenv("VAULTPULL_LOWERCASE_KEYS", "1")
	cfg := LowercaseConfigFromEnv()
	if !cfg.Enabled {
		t.Error("expected Enabled to be true for value '1'")
	}
}

func TestLowercaseKeys_Disabled_ReturnsOriginal(t *testing.T) {
	input := map[string]string{"DB_HOST": "localhost", "API_KEY": "secret"}
	result := LowercaseKeys(LowercaseConfig{Enabled: false}, input)
	if result["DB_HOST"] != "localhost" || result["API_KEY"] != "secret" {
		t.Error("expected original keys to be preserved when disabled")
	}
}

func TestLowercaseKeys_ConvertsKeys(t *testing.T) {
	input := map[string]string{"DB_HOST": "localhost", "API_KEY": "secret"}
	result := LowercaseKeys(LowercaseConfig{Enabled: true}, input)
	if result["db_host"] != "localhost" {
		t.Errorf("expected db_host=localhost, got %q", result["db_host"])
	}
	if result["api_key"] != "secret" {
		t.Errorf("expected api_key=secret, got %q", result["api_key"])
	}
	if _, ok := result["DB_HOST"]; ok {
		t.Error("expected original uppercase key to be absent")
	}
}

func TestLowercaseKeys_DoesNotMutateInput(t *testing.T) {
	input := map[string]string{"MY_VAR": "value"}
	_ = LowercaseKeys(LowercaseConfig{Enabled: true}, input)
	if _, ok := input["MY_VAR"]; !ok {
		t.Error("expected original map to remain unchanged")
	}
}

func TestLowercaseKeys_EmptyMap(t *testing.T) {
	result := LowercaseKeys(LowercaseConfig{Enabled: true}, map[string]string{})
	if len(result) != 0 {
		t.Errorf("expected empty result, got %d entries", len(result))
	}
}
