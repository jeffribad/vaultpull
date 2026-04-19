package sync

import (
	"testing"
)

func TestEnvMap_Integration_ConfigDriven(t *testing.T) {
	t.Setenv("VAULTPULL_ENVMAP_ENABLED", "true")
	t.Setenv("VAULTPULL_ENVMAP_KEYS", "db_password:DATABASE_PASSWORD,api_token:SERVICE_TOKEN")

	cfg := EnvMapConfigFromEnv()
	secrets := map[string]string{
		"db_password": "supersecret",
		"api_token":   "tok_abc123",
		"unrelated":   "keep",
	}

	result := ApplyEnvMap(cfg, secrets)

	if result["DATABASE_PASSWORD"] != "supersecret" {
		t.Errorf("expected DATABASE_PASSWORD remapped, got %v", result)
	}
	if result["SERVICE_TOKEN"] != "tok_abc123" {
		t.Errorf("expected SERVICE_TOKEN remapped, got %v", result)
	}
	if result["unrelated"] != "keep" {
		t.Errorf("expected unrelated preserved, got %v", result)
	}
	if _, ok := result["db_password"]; ok {
		t.Error("old key db_password should be removed after remapping")
	}
}

func TestEnvMap_Integration_DisabledPassthrough(t *testing.T) {
	t.Setenv("VAULTPULL_ENVMAP_ENABLED", "false")
	t.Setenv("VAULTPULL_ENVMAP_KEYS", "db_password:DATABASE_PASSWORD")

	cfg := EnvMapConfigFromEnv()
	secrets := map[string]string{"db_password": "val"}

	result := ApplyEnvMap(cfg, secrets)
	if result["db_password"] != "val" {
		t.Error("disabled envmap should not rename keys")
	}
}
