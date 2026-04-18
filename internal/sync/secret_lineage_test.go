package sync

import (
	"strings"
	"testing"
	"os"
)

func TestLineageConfigFromEnv_Defaults(t *testing.T) {
	os.Unsetenv("VAULTPULL_LINEAGE_ENABLED")
	os.Unsetenv("VAULTPULL_LINEAGE_ANNOTATE")
	cfg := LineageConfigFromEnv()
	if cfg.Enabled {
		t.Error("expected Enabled=false by default")
	}
	if cfg.Annotate {
		t.Error("expected Annotate=false by default")
	}
}

func TestLineageConfigFromEnv_Enabled(t *testing.T) {
	t.Setenv("VAULTPULL_LINEAGE_ENABLED", "true")
	t.Setenv("VAULTPULL_LINEAGE_ANNOTATE", "1")
	t.Setenv("VAULT_ADDR", "https://vault.example.com")
	t.Setenv("VAULTPULL_SECRET_PATH", "secret/myapp")
	cfg := LineageConfigFromEnv()
	if !cfg.Enabled {
		t.Error("expected Enabled=true")
	}
	if !cfg.Annotate {
		t.Error("expected Annotate=true")
	}
	if cfg.VaultAddr != "https://vault.example.com" {
		t.Errorf("unexpected VaultAddr: %s", cfg.VaultAddr)
	}
}

func TestLineageConfigFromEnv_NumericEnabled(t *testing.T) {
	t.Setenv("VAULTPULL_LINEAGE_ENABLED", "1")
	cfg := LineageConfigFromEnv()
	if !cfg.Enabled {
		t.Error("expected Enabled=true for value '1'")
	}
}

func TestInjectLineage_Disabled_ReturnsOriginal(t *testing.T) {
	cfg := LineageConfig{Enabled: false, Annotate: false}
	lines := []string{"KEY=val"}
	out := InjectLineage(cfg, map[string]string{"KEY": "val"}, lines)
	if len(out) != 1 || out[0] != "KEY=val" {
		t.Errorf("expected original lines, got %v", out)
	}
}

func TestInjectLineage_InjectsHeader(t *testing.T) {
	cfg := LineageConfig{
		Enabled:    true,
		Annotate:   true,
		VaultAddr:  "https://vault.local",
		SecretPath: "secret/app",
	}
	secrets := map[string]string{"FOO": "bar", "BAZ": "qux"}
	lines := []string{"FOO=bar", "BAZ=qux"}
	out := InjectLineage(cfg, secrets, lines)
	if len(out) < 4 {
		t.Fatalf("expected at least 4 lines, got %d", len(out))
	}
	if !strings.HasPrefix(out[0], "# vaultpull lineage:") {
		t.Errorf("expected lineage comment, got: %s", out[0])
	}
	if !strings.Contains(out[1], "vault.local") {
		t.Errorf("expected vault addr in source line, got: %s", out[1])
	}
	if out[3] != "" {
		t.Errorf("expected blank separator line, got: %q", out[3])
	}
}

func TestInjectLineage_EmptySecrets_ReturnsOriginal(t *testing.T) {
	cfg := LineageConfig{Enabled: true, Annotate: true}
	out := InjectLineage(cfg, map[string]string{}, []string{"X=1"})
	if len(out) != 1 {
		t.Errorf("expected 1 line for empty secrets, got %d", len(out))
	}
}

func TestLineageRecord_String(t *testing.T) {
	r := LineageRecord{
		KeyCount:   5,
		SecretPath: "secret/app",
		VaultAddr:  "https://vault.io",
	}
	s := r.String()
	if !strings.Contains(s, "5 keys") {
		t.Errorf("expected key count in string, got: %s", s)
	}
}
