package sync

import (
	"testing"
)

func TestInheritConfigFromEnv_Defaults(t *testing.T) {
	t.Setenv("VAULTPULL_INHERIT_ENABLED", "")
	t.Setenv("VAULTPULL_INHERIT_KEYS", "")
	t.Setenv("VAULTPULL_INHERIT_PREFIX", "")
	t.Setenv("VAULTPULL_INHERIT_OVERRIDE", "")

	cfg := InheritConfigFromEnv()
	if cfg.Enabled {
		t.Error("expected Enabled=false by default")
	}
	if len(cfg.Keys) != 0 {
		t.Errorf("expected empty keys, got %v", cfg.Keys)
	}
	if cfg.Prefix != "" {
		t.Errorf("expected empty prefix, got %s", cfg.Prefix)
	}
	if cfg.Override {
		t.Error("expected Override=false by default")
	}
}

func TestInheritConfigFromEnv_Enabled(t *testing.T) {
	t.Setenv("VAULTPULL_INHERIT_ENABLED", "true")
	cfg := InheritConfigFromEnv()
	if !cfg.Enabled {
		t.Error("expected Enabled=true")
	}
}

func TestInheritConfigFromEnv_NumericEnabled(t *testing.T) {
	t.Setenv("VAULTPULL_INHERIT_ENABLED", "1")
	cfg := InheritConfigFromEnv()
	if !cfg.Enabled {
		t.Error("expected Enabled=true for value '1'")
	}
}

func TestInheritConfigFromEnv_ParsesKeys(t *testing.T) {
	t.Setenv("VAULTPULL_INHERIT_ENABLED", "true")
	t.Setenv("VAULTPULL_INHERIT_KEYS", "KEY_A, KEY_B , KEY_C")
	cfg := InheritConfigFromEnv()
	if len(cfg.Keys) != 3 {
		t.Fatalf("expected 3 keys, got %d", len(cfg.Keys))
	}
	if cfg.Keys[1] != "KEY_B" {
		t.Errorf("expected KEY_B, got %s", cfg.Keys[1])
	}
}

func TestInheritConfigFromEnv_ParsesPrefix(t *testing.T) {
	t.Setenv("VAULTPULL_INHERIT_ENABLED", "true")
	t.Setenv("VAULTPULL_INHERIT_PREFIX", "MYAPP_")
	cfg := InheritConfigFromEnv()
	if cfg.Prefix != "MYAPP_" {
		t.Errorf("expected MYAPP_, got %s", cfg.Prefix)
	}
}
