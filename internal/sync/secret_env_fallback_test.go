package sync

import (
	"testing"
)

func TestApplyFallback_Disabled_ReturnsOriginal(t *testing.T) {
	cfg := FallbackConfig{Enabled: false}
	secrets := map[string]string{"KEY": "vault"}
	out := ApplyFallback(cfg, secrets)
	if out["KEY"] != "vault" {
		t.Errorf("expected vault, got %s", out["KEY"])
	}
}

func TestApplyFallback_DoesNotOverrideExisting(t *testing.T) {
	t.Setenv("DB_URL", "env-value")
	cfg := FallbackConfig{Enabled: true, Keys: []string{"DB_URL"}}
	secrets := map[string]string{"DB_URL": "vault-value"}
	out := ApplyFallback(cfg, secrets)
	if out["DB_URL"] != "vault-value" {
		t.Errorf("expected vault-value, got %s", out["DB_URL"])
	}
}

func TestApplyFallback_FillsMissingKey(t *testing.T) {
	t.Setenv("MISSING_KEY", "from-env")
	cfg := FallbackConfig{Enabled: true, Keys: []string{"MISSING_KEY"}}
	secrets := map[string]string{}
	out := ApplyFallback(cfg, secrets)
	if out["MISSING_KEY"] != "from-env" {
		t.Errorf("expected from-env, got %s", out["MISSING_KEY"])
	}
}

func TestApplyFallback_WithPrefix(t *testing.T) {
	t.Setenv("FALLBACK_TOKEN", "prefixed-val")
	cfg := FallbackConfig{Enabled: true, Keys: []string{"TOKEN"}, Prefix: "FALLBACK_"}
	secrets := map[string]string{}
	out := ApplyFallback(cfg, secrets)
	if out["TOKEN"] != "prefixed-val" {
		t.Errorf("expected prefixed-val, got %s", out["TOKEN"])
	}
}

func TestApplyFallback_DoesNotMutateInput(t *testing.T) {
	t.Setenv("NEW_KEY", "injected")
	cfg := FallbackConfig{Enabled: true, Keys: []string{"NEW_KEY"}}
	original := map[string]string{"EXISTING": "val"}
	ApplyFallback(cfg, original)
	if _, ok := original["NEW_KEY"]; ok {
		t.Error("original map was mutated")
	}
}
