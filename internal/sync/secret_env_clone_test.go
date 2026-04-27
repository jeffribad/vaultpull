package sync

import (
	"os"
	"testing"
)

func TestApplyClone_Disabled_ReturnsOriginal(t *testing.T) {
	secrets := map[string]string{"PROD_DB": "postgres"}
	cfg := CloneConfig{Enabled: false, FromPrefix: "PROD_", ToPrefix: "STAGING_"}
	out := ApplyClone(secrets, cfg)
	if _, ok := out["STAGING_DB"]; ok {
		t.Fatal("expected no cloned key when disabled")
	}
}

func TestApplyClone_EmptyFromPrefix_ReturnsOriginal(t *testing.T) {
	secrets := map[string]string{"PROD_DB": "postgres"}
	cfg := CloneConfig{Enabled: true, FromPrefix: "", ToPrefix: "STAGING_"}
	out := ApplyClone(secrets, cfg)
	if len(out) != 1 {
		t.Fatalf("expected 1 key, got %d", len(out))
	}
}

func TestApplyClone_CopiesMatchingKeys(t *testing.T) {
	secrets := map[string]string{"PROD_DB": "postgres", "PROD_API": "key123", "OTHER": "val"}
	cfg := CloneConfig{Enabled: true, FromPrefix: "PROD_", ToPrefix: "STAGING_"}
	out := ApplyClone(secrets, cfg)
	if out["STAGING_DB"] != "postgres" {
		t.Errorf("expected STAGING_DB=postgres, got %q", out["STAGING_DB"])
	}
	if out["STAGING_API"] != "key123" {
		t.Errorf("expected STAGING_API=key123, got %q", out["STAGING_API"])
	}
	if out["OTHER"] != "val" {
		t.Error("expected OTHER to be preserved")
	}
}

func TestApplyClone_NoOverwrite_KeepsExisting(t *testing.T) {
	secrets := map[string]string{"PROD_DB": "new", "STAGING_DB": "old"}
	cfg := CloneConfig{Enabled: true, FromPrefix: "PROD_", ToPrefix: "STAGING_", Overwrite: false}
	out := ApplyClone(secrets, cfg)
	if out["STAGING_DB"] != "old" {
		t.Errorf("expected existing STAGING_DB to be preserved, got %q", out["STAGING_DB"])
	}
}

func TestApplyClone_Overwrite_ReplacesExisting(t *testing.T) {
	secrets := map[string]string{"PROD_DB": "new", "STAGING_DB": "old"}
	cfg := CloneConfig{Enabled: true, FromPrefix: "PROD_", ToPrefix: "STAGING_", Overwrite: true}
	out := ApplyClone(secrets, cfg)
	if out["STAGING_DB"] != "new" {
		t.Errorf("expected STAGING_DB to be overwritten, got %q", out["STAGING_DB"])
	}
}

func TestApplyClone_DoesNotMutateInput(t *testing.T) {
	secrets := map[string]string{"PROD_KEY": "val"}
	cfg := CloneConfig{Enabled: true, FromPrefix: "PROD_", ToPrefix: "DEV_"}
	ApplyClone(secrets, cfg)
	if _, ok := secrets["DEV_KEY"]; ok {
		t.Fatal("input map was mutated")
	}
}

func TestCloneConfigFromEnv_Defaults(t *testing.T) {
	os.Unsetenv("VAULTPULL_CLONE_ENABLED")
	os.Unsetenv("VAULTPULL_CLONE_FROM")
	os.Unsetenv("VAULTPULL_CLONE_TO")
	os.Unsetenv("VAULTPULL_CLONE_OVERWRITE")
	cfg := CloneConfigFromEnv()
	if cfg.Enabled || cfg.FromPrefix != "" || cfg.ToPrefix != "" || cfg.Overwrite {
		t.Errorf("unexpected defaults: %+v", cfg)
	}
}

func TestCloneConfigFromEnv_Enabled(t *testing.T) {
	t.Setenv("VAULTPULL_CLONE_ENABLED", "true")
	t.Setenv("VAULTPULL_CLONE_FROM", "PROD_")
	t.Setenv("VAULTPULL_CLONE_TO", "STAGING_")
	t.Setenv("VAULTPULL_CLONE_OVERWRITE", "1")
	cfg := CloneConfigFromEnv()
	if !cfg.Enabled {
		t.Error("expected Enabled=true")
	}
	if cfg.FromPrefix != "PROD_" {
		t.Errorf("expected FromPrefix=PROD_, got %q", cfg.FromPrefix)
	}
	if cfg.ToPrefix != "STAGING_" {
		t.Errorf("expected ToPrefix=STAGING_, got %q", cfg.ToPrefix)
	}
	if !cfg.Overwrite {
		t.Error("expected Overwrite=true")
	}
}
