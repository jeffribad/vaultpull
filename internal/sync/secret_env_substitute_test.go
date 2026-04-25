package sync

import (
	"os"
	"testing"
)

func TestSubstituteSecrets_Disabled_ReturnsOriginal(t *testing.T) {
	cfg := SubstituteConfig{Enabled: false}
	input := map[string]string{"KEY": "${OTHER}"}
	out, err := SubstituteSecrets(cfg, input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["KEY"] != "${OTHER}" {
		t.Errorf("expected original value, got %q", out["KEY"])
	}
}

func TestSubstituteSecrets_ResolvesFromMap(t *testing.T) {
	cfg := SubstituteConfig{Enabled: true}
	input := map[string]string{
		"BASE_URL": "https://example.com",
		"API_URL":  "${BASE_URL}/api",
	}
	out, err := SubstituteSecrets(cfg, input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["API_URL"] != "https://example.com/api" {
		t.Errorf("expected resolved URL, got %q", out["API_URL"])
	}
}

func TestSubstituteSecrets_UnresolvedRef_ReturnsError(t *testing.T) {
	cfg := SubstituteConfig{Enabled: true, AllowEmpty: false}
	input := map[string]string{"KEY": "${MISSING_VAR_XYZ}"}
	_, err := SubstituteSecrets(cfg, input)
	if err == nil {
		t.Fatal("expected error for unresolved variable")
	}
}

func TestSubstituteSecrets_AllowEmpty_ReturnsEmptyString(t *testing.T) {
	cfg := SubstituteConfig{Enabled: true, AllowEmpty: true}
	input := map[string]string{"KEY": "prefix-${MISSING_VAR_XYZ}-suffix"}
	out, err := SubstituteSecrets(cfg, input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["KEY"] != "prefix--suffix" {
		t.Errorf("expected empty substitution, got %q", out["KEY"])
	}
}

func TestSubstituteSecrets_ResolvesFromOS(t *testing.T) {
	t.Setenv("VAULTPULL_TEST_OSVAR", "from-os")
	cfg := SubstituteConfig{Enabled: true}
	input := map[string]string{"KEY": "val-${VAULTPULL_TEST_OSVAR}"}
	out, err := SubstituteSecrets(cfg, input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["KEY"] != "val-from-os" {
		t.Errorf("expected OS var resolved, got %q", out["KEY"])
	}
}

func TestSubstituteSecrets_PrefixFilter_SkipsNonMatching(t *testing.T) {
	cfg := SubstituteConfig{Enabled: true, Prefix: "APP_"}
	input := map[string]string{
		"APP_HOST": "localhost",
		"URL":      "http://${APP_HOST}:${PORT}",
		"PORT":     "8080",
	}
	out, err := SubstituteSecrets(cfg, input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// APP_HOST matches prefix and resolves; PORT does not match prefix, kept as-is
	if out["URL"] != "http://localhost:${PORT}" {
		t.Errorf("expected partial substitution, got %q", out["URL"])
	}
}

func TestSubstituteSecrets_DoesNotMutateInput(t *testing.T) {
	cfg := SubstituteConfig{Enabled: true}
	input := map[string]string{"A": "hello", "B": "${A} world"}
	_, _ = SubstituteSecrets(cfg, input)
	if input["B"] != "${A} world" {
		t.Error("input map was mutated")
	}
}

func TestSubstituteConfigFromEnv_Defaults(t *testing.T) {
	os.Unsetenv("VAULTPULL_SUBSTITUTE_ENABLED")
	os.Unsetenv("VAULTPULL_SUBSTITUTE_ALLOW_EMPTY")
	os.Unsetenv("VAULTPULL_SUBSTITUTE_PREFIX")
	cfg := SubstituteConfigFromEnv()
	if cfg.Enabled {
		t.Error("expected Enabled=false by default")
	}
	if cfg.AllowEmpty {
		t.Error("expected AllowEmpty=false by default")
	}
	if cfg.Prefix != "" {
		t.Errorf("expected empty Prefix by default, got %q", cfg.Prefix)
	}
}

func TestSubstituteConfigFromEnv_Enabled(t *testing.T) {
	t.Setenv("VAULTPULL_SUBSTITUTE_ENABLED", "true")
	t.Setenv("VAULTPULL_SUBSTITUTE_ALLOW_EMPTY", "1")
	t.Setenv("VAULTPULL_SUBSTITUTE_PREFIX", "APP_")
	cfg := SubstituteConfigFromEnv()
	if !cfg.Enabled {
		t.Error("expected Enabled=true")
	}
	if !cfg.AllowEmpty {
		t.Error("expected AllowEmpty=true")
	}
	if cfg.Prefix != "APP_" {
		t.Errorf("expected Prefix=APP_, got %q", cfg.Prefix)
	}
}
