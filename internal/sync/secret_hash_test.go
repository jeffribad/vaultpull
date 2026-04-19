package sync

import (
	"crypto/sha256"
	"fmt"
	"os"
	"testing"
)

func TestHashSecrets_Disabled_ReturnsOriginal(t *testing.T) {
	cfg := HashConfig{Enabled: false, Keys: []string{"SECRET"}}
	in := map[string]string{"SECRET": "mysecret"}
	out := HashSecrets(cfg, in)
	if out["SECRET"] != "mysecret" {
		t.Errorf("expected original value, got %q", out["SECRET"])
	}
}

func TestHashSecrets_NoKeys_ReturnsOriginal(t *testing.T) {
	cfg := HashConfig{Enabled: true, Keys: []string{}}
	in := map[string]string{"SECRET": "mysecret"}
	out := HashSecrets(cfg, in)
	if out["SECRET"] != "mysecret" {
		t.Errorf("expected original value, got %q", out["SECRET"])
	}
}

func TestHashSecrets_HashesConfiguredKey(t *testing.T) {
	cfg := HashConfig{Enabled: true, Keys: []string{"API_SECRET"}}
	in := map[string]string{"API_SECRET": "hunter2", "DB_HOST": "localhost"}
	out := HashSecrets(cfg, in)

	sum := sha256.Sum256([]byte("hunter2"))
	expected := fmt.Sprintf("%x", sum)
	if out["API_SECRET"] != expected {
		t.Errorf("expected hash %q, got %q", expected, out["API_SECRET"])
	}
	if out["DB_HOST"] != "localhost" {
		t.Errorf("expected DB_HOST unchanged, got %q", out["DB_HOST"])
	}
}

func TestHashSecrets_CaseInsensitiveKeyMatch(t *testing.T) {
	cfg := HashConfig{Enabled: true, Keys: []string{"api_secret"}}
	in := map[string]string{"API_SECRET": "value"}
	out := HashSecrets(cfg, in)

	sum := sha256.Sum256([]byte("value"))
	expected := fmt.Sprintf("%x", sum)
	if out["API_SECRET"] != expected {
		t.Errorf("expected hashed value, got %q", out["API_SECRET"])
	}
}

func TestHashSecrets_DoesNotMutateInput(t *testing.T) {
	cfg := HashConfig{Enabled: true, Keys: []string{"TOKEN"}}
	in := map[string]string{"TOKEN": "original"}
	HashSecrets(cfg, in)
	if in["TOKEN"] != "original" {
		t.Error("input map was mutated")
	}
}

func TestHashConfigFromEnv_Defaults(t *testing.T) {
	os.Unsetenv("VAULTPULL_HASH_ENABLED")
	os.Unsetenv("VAULTPULL_HASH_KEYS")
	cfg := HashConfigFromEnv()
	if cfg.Enabled {
		t.Error("expected Enabled=false by default")
	}
	if len(cfg.Keys) != 0 {
		t.Errorf("expected no keys by default, got %v", cfg.Keys)
	}
}

func TestHashConfigFromEnv_Enabled(t *testing.T) {
	t.Setenv("VAULTPULL_HASH_ENABLED", "true")
	t.Setenv("VAULTPULL_HASH_KEYS", "SECRET, TOKEN")
	cfg := HashConfigFromEnv()
	if !cfg.Enabled {
		t.Error("expected Enabled=true")
	}
	if len(cfg.Keys) != 2 || cfg.Keys[0] != "SECRET" || cfg.Keys[1] != "TOKEN" {
		t.Errorf("unexpected keys: %v", cfg.Keys)
	}
}
