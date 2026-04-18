package sync

import (
	"os"
	"testing"
)

func TestAliasConfigFromEnv_Defaults(t *testing.T) {
	os.Unsetenv("VAULTPULL_ALIAS_ENABLED")
	os.Unsetenv("VAULTPULL_ALIASES")
	cfg := AliasConfigFromEnv()
	if cfg.Enabled {
		t.Fatal("expected disabled by default")
	}
	if len(cfg.Aliases) != 0 {
		t.Fatalf("expected no aliases by default, got %v", cfg.Aliases)
	}
}

func TestAliasConfigFromEnv_Enabled(t *testing.T) {
	os.Setenv("VAULTPULL_ALIAS_ENABLED", "1")
	os.Setenv("VAULTPULL_ALIASES", "DB_HOST:DATABASE_HOST")
	defer os.Unsetenv("VAULTPULL_ALIAS_ENABLED")
	defer os.Unsetenv("VAULTPULL_ALIASES")
	cfg := AliasConfigFromEnv()
	if !cfg.Enabled {
		t.Fatal("expected enabled")
	}
	if len(cfg.Aliases["DB_HOST"]) != 1 {
		t.Fatalf("expected 1 alias for DB_HOST, got %v", cfg.Aliases)
	}
}

func TestAliasConfigFromEnv_TrueString(t *testing.T) {
	os.Setenv("VAULTPULL_ALIAS_ENABLED", "true")
	defer os.Unsetenv("VAULTPULL_ALIAS_ENABLED")
	cfg := AliasConfigFromEnv()
	if !cfg.Enabled {
		t.Fatal("expected enabled with 'true' string")
	}
}
