package sync

import (
	"os"
	"testing"
)

func TestApplyAliases_Integration_ConfigDriven(t *testing.T) {
	os.Setenv("VAULTPULL_ALIAS_ENABLED", "1")
	os.Setenv("VAULTPULL_ALIASES", "DB_HOST:DATABASE_HOST,DB_HOST:PG_HOST")
	defer os.Unsetenv("VAULTPULL_ALIAS_ENABLED")
	defer os.Unsetenv("VAULTPULL_ALIASES")

	cfg := AliasConfigFromEnv()
	secrets := map[string]string{
		"DB_HOST": "db.internal",
		"APP_KEY": "secret123",
	}
	out := ApplyAliases(secrets, cfg)

	if out["DATABASE_HOST"] != "db.internal" {
		t.Errorf("expected DATABASE_HOST=db.internal, got %q", out["DATABASE_HOST"])
	}
	if out["PG_HOST"] != "db.internal" {
		t.Errorf("expected PG_HOST=db.internal, got %q", out["PG_HOST"])
	}
	if out["DB_HOST"] != "db.internal" {
		t.Error("original DB_HOST should be preserved")
	}
	if out["APP_KEY"] != "secret123" {
		t.Error("unrelated key APP_KEY should be preserved")
	}
}

func TestApplyAliases_Integration_DisabledPassthrough(t *testing.T) {
	os.Setenv("VAULTPULL_ALIAS_ENABLED", "0")
	os.Setenv("VAULTPULL_ALIASES", "DB_HOST:DATABASE_HOST")
	defer os.Unsetenv("VAULTPULL_ALIAS_ENABLED")
	defer os.Unsetenv("VAULTPULL_ALIASES")

	cfg := AliasConfigFromEnv()
	secrets := map[string]string{"DB_HOST": "db.internal"}
	out := ApplyAliases(secrets, cfg)

	if _, ok := out["DATABASE_HOST"]; ok {
		t.Error("alias should not be injected when disabled")
	}
	if len(out) != 1 {
		t.Errorf("expected 1 key, got %d", len(out))
	}
}
