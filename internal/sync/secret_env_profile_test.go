package sync

import (
	"testing"
)

func TestApplyEnvProfile_Disabled_ReturnsOriginal(t *testing.T) {
	secrets := map[string]string{"KEY__prod": "val", "OTHER": "x"}
	cfg := EnvProfileConfig{Enabled: false, Profile: "prod", Profiles: []string{"prod"}}
	out := ApplyEnvProfile(cfg, secrets)
	if len(out) != 2 {
		t.Errorf("expected 2 keys, got %d", len(out))
	}
}

func TestApplyEnvProfile_NoProfile_ReturnsOriginal(t *testing.T) {
	secrets := map[string]string{"KEY__prod": "val"}
	cfg := EnvProfileConfig{Enabled: true, Profile: "", Profiles: []string{"prod"}}
	out := ApplyEnvProfile(cfg, secrets)
	if len(out) != 1 {
		t.Errorf("expected 1 key, got %d", len(out))
	}
}

func TestApplyEnvProfile_RenamesMatchingSuffix(t *testing.T) {
	secrets := map[string]string{"DB_URL__prod": "postgres://prod", "APP_NAME": "myapp"}
	cfg := EnvProfileConfig{Enabled: true, Profile: "prod", Profiles: []string{"dev", "prod"}}
	out := ApplyEnvProfile(cfg, secrets)
	if v, ok := out["DB_URL"]; !ok || v != "postgres://prod" {
		t.Errorf("expected DB_URL=postgres://prod, got %v", out)
	}
	if _, ok := out["DB_URL__prod"]; ok {
		t.Error("suffixed key should be removed")
	}
}

func TestApplyEnvProfile_ExcludesOtherProfileKeys(t *testing.T) {
	secrets := map[string]string{
		"KEY__dev":  "dev-val",
		"KEY__prod": "prod-val",
		"SHARED":    "shared",
	}
	cfg := EnvProfileConfig{Enabled: true, Profile: "prod", Profiles: []string{"dev", "prod"}}
	out := ApplyEnvProfile(cfg, secrets)
	if _, ok := out["KEY__dev"]; ok {
		t.Error("dev-profile key should be excluded")
	}
	if v, ok := out["KEY"]; !ok || v != "prod-val" {
		t.Errorf("expected KEY=prod-val, got %v", out)
	}
	if _, ok := out["SHARED"]; !ok {
		t.Error("shared key should be included")
	}
}

func TestApplyEnvProfile_DoesNotMutateInput(t *testing.T) {
	secrets := map[string]string{"KEY__prod": "v"}
	orig := map[string]string{"KEY__prod": "v"}
	cfg := EnvProfileConfig{Enabled: true, Profile: "prod", Profiles: []string{"prod"}}
	ApplyEnvProfile(cfg, secrets)
	for k, v := range orig {
		if secrets[k] != v {
			t.Error("input map was mutated")
		}
	}
}
