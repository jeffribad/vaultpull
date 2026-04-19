package sync

import (
	"testing"
)

func TestEnvLock_Integration_ConfigDriven(t *testing.T) {
	t.Setenv("VAULTPULL_LOCK_ENABLED", "true")
	t.Setenv("VAULTPULL_LOCK_KEYS", "APP_ENV=production")

	cfg := LockConfigFromEnv()

	secrets := map[string]string{
		"APP_ENV": "production",
		"DB_URL":  "postgres://localhost/app",
	}

	if err := EnforceLock(cfg, secrets); err != nil {
		t.Fatalf("expected no error for matching locked value: %v", err)
	}
}

func TestEnvLock_Integration_DetectsViolation(t *testing.T) {
	t.Setenv("VAULTPULL_LOCK_ENABLED", "true")
	t.Setenv("VAULTPULL_LOCK_KEYS", "APP_ENV=production")

	cfg := LockConfigFromEnv()

	secrets := map[string]string{
		"APP_ENV": "staging",
	}

	if err := EnforceLock(cfg, secrets); err == nil {
		t.Fatal("expected error for lock violation")
	}
}

func TestEnvLock_Integration_DisabledPassthrough(t *testing.T) {
	t.Setenv("VAULTPULL_LOCK_ENABLED", "false")
	t.Setenv("VAULTPULL_LOCK_KEYS", "APP_ENV=production")

	cfg := LockConfigFromEnv()

	secrets := map[string]string{
		"APP_ENV": "staging",
	}

	if err := EnforceLock(cfg, secrets); err != nil {
		t.Fatalf("expected nil when disabled, got %v", err)
	}
}
