package sync

import (
	"testing"
)

func TestEnforceLock_Disabled_ReturnsNil(t *testing.T) {
	cfg := LockConfig{Enabled: false, LockedKeys: map[string]string{"KEY": "val"}}
	secrets := map[string]string{"KEY": "other"}
	if err := EnforceLock(cfg, secrets); err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestEnforceLock_NoKeys_ReturnsNil(t *testing.T) {
	cfg := LockConfig{Enabled: true, LockedKeys: map[string]string{}}
	secrets := map[string]string{"KEY": "val"}
	if err := EnforceLock(cfg, secrets); err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestEnforceLock_MatchingValue_ReturnsNil(t *testing.T) {
	cfg := LockConfig{Enabled: true, LockedKeys: map[string]string{"ENV": "production"}}
	secrets := map[string]string{"ENV": "production"}
	if err := EnforceLock(cfg, secrets); err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestEnforceLock_MismatchedValue_ReturnsError(t *testing.T) {
	cfg := LockConfig{Enabled: true, LockedKeys: map[string]string{"ENV": "production"}}
	secrets := map[string]string{"ENV": "staging"}
	if err := EnforceLock(cfg, secrets); err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestEnforceLock_MissingKey_Skipped(t *testing.T) {
	cfg := LockConfig{Enabled: true, LockedKeys: map[string]string{"MISSING": "val"}}
	secrets := map[string]string{"OTHER": "val"}
	if err := EnforceLock(cfg, secrets); err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestLockConfigFromEnv_Defaults(t *testing.T) {
	t.Setenv("VAULTPULL_LOCK_ENABLED", "")
	t.Setenv("VAULTPULL_LOCK_KEYS", "")
	cfg := LockConfigFromEnv()
	if cfg.Enabled {
		t.Error("expected Enabled=false")
	}
	if len(cfg.LockedKeys) != 0 {
		t.Errorf("expected empty LockedKeys, got %v", cfg.LockedKeys)
	}
}

func TestLockConfigFromEnv_ParsesKeys(t *testing.T) {
	t.Setenv("VAULTPULL_LOCK_ENABLED", "true")
	t.Setenv("VAULTPULL_LOCK_KEYS", "APP_ENV=production, DB_DRIVER=postgres")
	cfg := LockConfigFromEnv()
	if !cfg.Enabled {
		t.Error("expected Enabled=true")
	}
	if cfg.LockedKeys["APP_ENV"] != "production" {
		t.Errorf("unexpected value: %q", cfg.LockedKeys["APP_ENV"])
	}
	if cfg.LockedKeys["DB_DRIVER"] != "postgres" {
		t.Errorf("unexpected value: %q", cfg.LockedKeys["DB_DRIVER"])
	}
}
