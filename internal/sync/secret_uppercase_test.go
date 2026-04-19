package sync

import (
	"testing"
)

func TestUppercaseConfigFromEnv_Defaults(t *testing.T) {
	t.Setenv("VAULTPULL_UPPERCASE_KEYS", "")
	cfg := UppercaseConfigFromEnv()
	if cfg.Enabled {
		t.Error("expected Enabled=false by default")
	}
}

func TestUppercaseConfigFromEnv_Enabled(t *testing.T) {
	t.Setenv("VAULTPULL_UPPERCASE_KEYS", "true")
	cfg := UppercaseConfigFromEnv()
	if !cfg.Enabled {
		t.Error("expected Enabled=true")
	}
}

func TestUppercaseKeys_Disabled_ReturnsOriginal(t *testing.T) {
	input := map[string]string{"db_host": "localhost"}
	out := UppercaseKeys(UppercaseConfig{Enabled: false}, input)
	if out["db_host"] != "localhost" {
		t.Errorf("expected original key preserved, got %v", out)
	}
	if _, ok := out["DB_HOST"]; ok {
		t.Error("did not expect uppercased key when disabled")
	}
}

func TestUppercaseKeys_ConvertsKeys(t *testing.T) {
	input := map[string]string{
		"db_host": "localhost",
		"api_key": "secret",
	}
	out := UppercaseKeys(UppercaseConfig{Enabled: true}, input)
	if out["DB_HOST"] != "localhost" {
		t.Errorf("expected DB_HOST=localhost, got %v", out["DB_HOST"])
	}
	if out["API_KEY"] != "secret" {
		t.Errorf("expected API_KEY=secret, got %v", out["API_KEY"])
	}
	if _, ok := out["db_host"]; ok {
		t.Error("original lowercase key should not be present")
	}
}

func TestUppercaseKeys_DoesNotMutateInput(t *testing.T) {
	input := map[string]string{"foo": "bar"}
	UppercaseKeys(UppercaseConfig{Enabled: true}, input)
	if _, ok := input["FOO"]; ok {
		t.Error("input map should not be mutated")
	}
}

func TestUppercaseKeys_EmptyMap(t *testing.T) {
	out := UppercaseKeys(UppercaseConfig{Enabled: true}, map[string]string{})
	if len(out) != 0 {
		t.Errorf("expected empty map, got %v", out)
	}
}
