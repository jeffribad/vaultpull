package sync

import (
	"testing"
)

func TestApplyInherit_Disabled_ReturnsOriginal(t *testing.T) {
	cfg := InheritConfig{Enabled: false}
	secrets := map[string]string{"KEY": "vault"}
	out := ApplyInherit(cfg, secrets)
	if out["KEY"] != "vault" {
		t.Errorf("expected vault, got %s", out["KEY"])
	}
}

func TestApplyInherit_SpecificKeys_FillsMissing(t *testing.T) {
	t.Setenv("MY_KEY", "from-env")
	cfg := InheritConfig{Enabled: true, Keys: []string{"MY_KEY"}}
	secrets := map[string]string{}
	out := ApplyInherit(cfg, secrets)
	if out["MY_KEY"] != "from-env" {
		t.Errorf("expected from-env, got %s", out["MY_KEY"])
	}
}

func TestApplyInherit_SpecificKeys_NoOverride_KeepsVault(t *testing.T) {
	t.Setenv("MY_KEY", "from-env")
	cfg := InheritConfig{Enabled: true, Keys: []string{"MY_KEY"}, Override: false}
	secrets := map[string]string{"MY_KEY": "vault-value"}
	out := ApplyInherit(cfg, secrets)
	if out["MY_KEY"] != "vault-value" {
		t.Errorf("expected vault-value, got %s", out["MY_KEY"])
	}
}

func TestApplyInherit_SpecificKeys_Override_UsesEnv(t *testing.T) {
	t.Setenv("MY_KEY", "from-env")
	cfg := InheritConfig{Enabled: true, Keys: []string{"MY_KEY"}, Override: true}
	secrets := map[string]string{"MY_KEY": "vault-value"}
	out := ApplyInherit(cfg, secrets)
	if out["MY_KEY"] != "from-env" {
		t.Errorf("expected from-env, got %s", out["MY_KEY"])
	}
}

func TestApplyInherit_Prefix_StripsAndMerges(t *testing.T) {
	t.Setenv("APP_DB_HOST", "localhost")
	t.Setenv("APP_DB_PORT", "5432")
	cfg := InheritConfig{Enabled: true, Prefix: "APP_"}
	secrets := map[string]string{}
	out := ApplyInherit(cfg, secrets)
	if out["DB_HOST"] != "localhost" {
		t.Errorf("expected localhost, got %s", out["DB_HOST"])
	}
	if out["DB_PORT"] != "5432" {
		t.Errorf("expected 5432, got %s", out["DB_PORT"])
	}
}

func TestApplyInherit_DoesNotMutateInput(t *testing.T) {
	t.Setenv("SOME_KEY", "env-val")
	cfg := InheritConfig{Enabled: true, Keys: []string{"SOME_KEY"}, Override: true}
	secrets := map[string]string{"SOME_KEY": "original"}
	ApplyInherit(cfg, secrets)
	if secrets["SOME_KEY"] != "original" {
		t.Error("input map was mutated")
	}
}
