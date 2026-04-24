package sync

import (
	"testing"
)

func TestApplyWhitelist_Disabled_ReturnsOriginal(t *testing.T) {
	cfg := WhitelistConfig{Enabled: false, Keys: []string{"DB_HOST"}}
	secrets := map[string]string{"DB_HOST": "localhost", "API_KEY": "secret"}

	result := ApplyWhitelist(cfg, secrets)

	if len(result) != 2 {
		t.Errorf("expected 2 keys, got %d", len(result))
	}
}

func TestApplyWhitelist_EmptyKeys_ReturnsOriginal(t *testing.T) {
	cfg := WhitelistConfig{Enabled: true, Keys: nil}
	secrets := map[string]string{"DB_HOST": "localhost", "API_KEY": "secret"}

	result := ApplyWhitelist(cfg, secrets)

	if len(result) != 2 {
		t.Errorf("expected 2 keys, got %d", len(result))
	}
}

func TestApplyWhitelist_FiltersToAllowedKeys(t *testing.T) {
	cfg := WhitelistConfig{Enabled: true, Keys: []string{"DB_HOST", "DB_PORT"}}
	secrets := map[string]string{
		"DB_HOST": "localhost",
		"DB_PORT": "5432",
		"API_KEY": "secret",
		"SECRET_TOKEN": "tok",
	}

	result := ApplyWhitelist(cfg, secrets)

	if len(result) != 2 {
		t.Errorf("expected 2 keys, got %d", len(result))
	}
	if result["DB_HOST"] != "localhost" {
		t.Errorf("expected DB_HOST=localhost")
	}
	if result["DB_PORT"] != "5432" {
		t.Errorf("expected DB_PORT=5432")
	}
	if _, ok := result["API_KEY"]; ok {
		t.Error("expected API_KEY to be removed")
	}
}

func TestApplyWhitelist_CaseInsensitiveMatch(t *testing.T) {
	cfg := WhitelistConfig{Enabled: true, Keys: []string{"db_host"}}
	secrets := map[string]string{"DB_HOST": "localhost", "API_KEY": "secret"}

	result := ApplyWhitelist(cfg, secrets)

	if len(result) != 1 {
		t.Errorf("expected 1 key, got %d", len(result))
	}
	if result["DB_HOST"] != "localhost" {
		t.Errorf("expected DB_HOST to be present")
	}
}

func TestApplyWhitelist_DoesNotMutateInput(t *testing.T) {
	cfg := WhitelistConfig{Enabled: true, Keys: []string{"DB_HOST"}}
	secrets := map[string]string{"DB_HOST": "localhost", "API_KEY": "secret"}

	_ = ApplyWhitelist(cfg, secrets)

	if len(secrets) != 2 {
		t.Error("input map was mutated")
	}
}

func TestWhitelistConfigFromEnv_Defaults(t *testing.T) {
	t.Setenv("VAULTPULL_WHITELIST_ENABLED", "")
	t.Setenv("VAULTPULL_WHITELIST_KEYS", "")

	cfg := WhitelistConfigFromEnv()

	if cfg.Enabled {
		t.Error("expected Enabled=false by default")
	}
	if len(cfg.Keys) != 0 {
		t.Errorf("expected no keys by default, got %v", cfg.Keys)
	}
}

func TestWhitelistConfigFromEnv_ParsesKeys(t *testing.T) {
	t.Setenv("VAULTPULL_WHITELIST_ENABLED", "true")
	t.Setenv("VAULTPULL_WHITELIST_KEYS", "DB_HOST, DB_PORT , API_KEY")

	cfg := WhitelistConfigFromEnv()

	if !cfg.Enabled {
		t.Error("expected Enabled=true")
	}
	if len(cfg.Keys) != 3 {
		t.Errorf("expected 3 keys, got %d: %v", len(cfg.Keys), cfg.Keys)
	}
}
