package sync

import (
	"testing"
)

func TestSubstituteSecrets_Integration_ConfigDriven(t *testing.T) {
	t.Setenv("VAULTPULL_SUBSTITUTE_ENABLED", "true")
	t.Setenv("VAULTPULL_SUBSTITUTE_ALLOW_EMPTY", "false")
	t.Setenv("VAULTPULL_SUBSTITUTE_PREFIX", "")

	cfg := SubstituteConfigFromEnv()

	secrets := map[string]string{
		"DB_HOST": "db.internal",
		"DB_PORT": "5432",
		"DB_DSN":  "postgres://${DB_HOST}:${DB_PORT}/mydb",
	}

	out, err := SubstituteSecrets(cfg, secrets)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := "postgres://db.internal:5432/mydb"
	if out["DB_DSN"] != want {
		t.Errorf("DB_DSN: want %q, got %q", want, out["DB_DSN"])
	}
	if out["DB_HOST"] != "db.internal" {
		t.Errorf("DB_HOST should be unchanged")
	}
}

func TestSubstituteSecrets_Integration_DisabledPassthrough(t *testing.T) {
	t.Setenv("VAULTPULL_SUBSTITUTE_ENABLED", "false")

	cfg := SubstituteConfigFromEnv()

	secrets := map[string]string{
		"BASE": "http://example.com",
		"URL":  "${BASE}/path",
	}

	out, err := SubstituteSecrets(cfg, secrets)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if out["URL"] != "${BASE}/path" {
		t.Errorf("expected passthrough when disabled, got %q", out["URL"])
	}
}
