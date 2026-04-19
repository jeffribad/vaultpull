package sync

import (
	"os"
	"testing"
)

func TestScopeConfigFromEnv_Defaults(t *testing.T) {
	os.Unsetenv("VAULTPULL_SCOPE")
	os.Unsetenv("VAULTPULL_SCOPE_STRIP")
	cfg := ScopeConfigFromEnv()
	if cfg.Enabled {
		t.Error("expected Enabled=false when no scope set")
	}
}

func TestScopeConfigFromEnv_Enabled(t *testing.T) {
	os.Setenv("VAULTPULL_SCOPE", "APP")
	defer os.Unsetenv("VAULTPULL_SCOPE")
	cfg := ScopeConfigFromEnv()
	if !cfg.Enabled {
		t.Error("expected Enabled=true when scope is set")
	}
	if cfg.Scope != "APP" {
		t.Errorf("expected scope APP, got %s", cfg.Scope)
	}
}

func TestApplyScope_Disabled_ReturnsOriginal(t *testing.T) {
	secrets := map[string]string{"APP_KEY": "val", "OTHER": "x"}
	cfg := ScopeConfig{Enabled: false, Scope: "APP"}
	out := ApplyScope(cfg, secrets)
	if len(out) != 2 {
		t.Errorf("expected 2 keys, got %d", len(out))
	}
}

func TestApplyScope_FiltersToScope(t *testing.T) {
	secrets := map[string]string{"APP_KEY": "val", "APP_SECRET": "s", "OTHER": "x"}
	cfg := ScopeConfig{Enabled: true, Scope: "APP"}
	out := ApplyScope(cfg, secrets)
	if len(out) != 2 {
		t.Errorf("expected 2 keys, got %d", len(out))
	}
	if _, ok := out["OTHER"]; ok {
		t.Error("OTHER should be excluded")
	}
}

func TestApplyScope_StripPrefix(t *testing.T) {
	secrets := map[string]string{"APP_KEY": "val", "APP_TOKEN": "tok"}
	cfg := ScopeConfig{Enabled: true, Scope: "APP", Strip: true}
	out := ApplyScope(cfg, secrets)
	if _, ok := out["KEY"]; !ok {
		t.Error("expected KEY after stripping APP_ prefix")
	}
	if _, ok := out["TOKEN"]; !ok {
		t.Error("expected TOKEN after stripping APP_ prefix")
	}
}

func TestApplyScope_CaseInsensitiveMatch(t *testing.T) {
	secrets := map[string]string{"app_key": "val", "OTHER": "x"}
	cfg := ScopeConfig{Enabled: true, Scope: "APP"}
	out := ApplyScope(cfg, secrets)
	if len(out) != 1 {
		t.Errorf("expected 1 key, got %d", len(out))
	}
}
