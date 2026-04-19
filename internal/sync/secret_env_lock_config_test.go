package sync

import (
	"testing"
)

func TestLockConfigFromEnv_NumericEnabled(t *testing.T) {
	t.Setenv("VAULTPULL_LOCK_ENABLED", "1")
	t.Setenv("VAULTPULL_LOCK_KEYS", "")
	cfg := LockConfigFromEnv()
	if !cfg.Enabled {
		t.Error("expected Enabled=true for value '1'")
	}
}

func TestLockConfigFromEnv_FalseString(t *testing.T) {
	t.Setenv("VAULTPULL_LOCK_ENABLED", "false")
	t.Setenv("VAULTPULL_LOCK_KEYS", "")
	cfg := LockConfigFromEnv()
	if cfg.Enabled {
		t.Error("expected Enabled=false for value 'false'")
	}
}

func TestLockConfigFromEnv_SkipsMalformedPairs(t *testing.T) {
	t.Setenv("VAULTPULL_LOCK_ENABLED", "true")
	t.Setenv("VAULTPULL_LOCK_KEYS", "GOOD=val,BADENTRY,ANOTHER=ok")
	cfg := LockConfigFromEnv()
	if len(cfg.LockedKeys) != 2 {
		t.Errorf("expected 2 keys, got %d: %v", len(cfg.LockedKeys), cfg.LockedKeys)
	}
	if _, ok := cfg.LockedKeys["BADENTRY"]; ok {
		t.Error("malformed entry should be skipped")
	}
}
