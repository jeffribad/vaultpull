package sync

import (
	"testing"
)

func TestBlacklist_Integration_ConfigDriven(t *testing.T) {
	t.Setenv("VAULTPULL_BLACKLIST_ENABLED", "true")
	t.Setenv("VAULTPULL_BLACKLIST_KEYS", "INTERNAL_TOKEN,DEBUG_MODE")

	cfg := BlacklistConfigFromEnv()

	secrets := map[string]string{
		"APP_NAME":       "myapp",
		"INTERNAL_TOKEN": "tok_secret",
		"DEBUG_MODE":     "true",
		"DB_HOST":        "localhost",
	}

	out := ApplyBlacklist(cfg, secrets)

	if _, ok := out["INTERNAL_TOKEN"]; ok {
		t.Error("INTERNAL_TOKEN should have been blacklisted")
	}
	if _, ok := out["DEBUG_MODE"]; ok {
		t.Error("DEBUG_MODE should have been blacklisted")
	}
	if out["APP_NAME"] != "myapp" {
		t.Errorf("APP_NAME should be preserved, got %q", out["APP_NAME"])
	}
	if out["DB_HOST"] != "localhost" {
		t.Errorf("DB_HOST should be preserved, got %q", out["DB_HOST"])
	}
}

func TestBlacklist_Integration_DisabledPassthrough(t *testing.T) {
	t.Setenv("VAULTPULL_BLACKLIST_ENABLED", "false")
	t.Setenv("VAULTPULL_BLACKLIST_KEYS", "INTERNAL_TOKEN")

	cfg := BlacklistConfigFromEnv()

	secrets := map[string]string{
		"INTERNAL_TOKEN": "should-survive",
		"OTHER":          "value",
	}

	out := ApplyBlacklist(cfg, secrets)

	if len(out) != 2 {
		t.Errorf("expected all 2 keys when disabled, got %d", len(out))
	}
	if out["INTERNAL_TOKEN"] != "should-survive" {
		t.Error("expected INTERNAL_TOKEN to survive when blacklist is disabled")
	}
}
