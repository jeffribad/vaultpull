package sync

import (
	"os"
	"testing"
)

func TestMergeSecrets_Integration_ConfigDriven(t *testing.T) {
	os.Setenv("VAULTPULL_MERGE_ENABLED", "true")
	os.Setenv("VAULTPULL_MERGE_STRATEGY", "vault-wins")
	defer os.Unsetenv("VAULTPULL_MERGE_ENABLED")
	defer os.Unsetenv("VAULTPULL_MERGE_STRATEGY")

	cfg := MergeConfigFromEnv()
	vault := map[string]string{"DB_HOST": "vault-host", "API_KEY": "vault-key"}
	base := map[string]string{"DB_HOST": "local-host", "LOG_LEVEL": "debug"}

	out := MergeSecrets(cfg, vault, base)

	if out["DB_HOST"] != "vault-host" {
		t.Errorf("vault-wins: expected vault-host, got %s", out["DB_HOST"])
	}
	if out["LOG_LEVEL"] != "debug" {
		t.Errorf("expected base-only key preserved, got %s", out["LOG_LEVEL"])
	}
	if out["API_KEY"] != "vault-key" {
		t.Errorf("expected vault-only key present, got %s", out["API_KEY"])
	}
}

func TestMergeSecrets_Integration_DisabledPassthrough(t *testing.T) {
	os.Setenv("VAULTPULL_MERGE_ENABLED", "false")
	defer os.Unsetenv("VAULTPULL_MERGE_ENABLED")

	cfg := MergeConfigFromEnv()
	vault := map[string]string{"X": "vault-x"}
	base := map[string]string{"X": "local-x", "Y": "local-y"}

	out := MergeSecrets(cfg, vault, base)

	if out["X"] != "vault-x" {
		t.Errorf("expected vault-x, got %s", out["X"])
	}
	if _, ok := out["Y"]; ok {
		t.Error("expected Y absent when merge disabled")
	}
}
