package sync

import (
	"testing"
	"time"
)

func TestTimeoutConfigFromEnv_Defaults(t *testing.T) {
	t.Setenv("VAULTPULL_VAULT_TIMEOUT_SEC", "")
	t.Setenv("VAULTPULL_WRITE_TIMEOUT_SEC", "")
	t.Setenv("VAULTPULL_GLOBAL_TIMEOUT_SEC", "")

	cfg := TimeoutConfigFromEnv()
	if cfg.VaultTimeout != 10*time.Second {
		t.Errorf("expected 10s, got %v", cfg.VaultTimeout)
	}
	if cfg.WriteTimeout != 5*time.Second {
		t.Errorf("expected 5s, got %v", cfg.WriteTimeout)
	}
	if cfg.GlobalTimeout != 30*time.Second {
		t.Errorf("expected 30s, got %v", cfg.GlobalTimeout)
	}
}

func TestTimeoutConfigFromEnv_CustomValues(t *testing.T) {
	t.Setenv("VAULTPULL_VAULT_TIMEOUT_SEC", "20")
	t.Setenv("VAULTPULL_WRITE_TIMEOUT_SEC", "8")
	t.Setenv("VAULTPULL_GLOBAL_TIMEOUT_SEC", "60")

	cfg := TimeoutConfigFromEnv()
	if cfg.VaultTimeout != 20*time.Second {
		t.Errorf("expected 20s, got %v", cfg.VaultTimeout)
	}
	if cfg.WriteTimeout != 8*time.Second {
		t.Errorf("expected 8s, got %v", cfg.WriteTimeout)
	}
	if cfg.GlobalTimeout != 60*time.Second {
		t.Errorf("expected 60s, got %v", cfg.GlobalTimeout)
	}
}

func TestTimeoutConfigFromEnv_InvalidValues_FallsBackToDefaults(t *testing.T) {
	t.Setenv("VAULTPULL_VAULT_TIMEOUT_SEC", "notanumber")
	t.Setenv("VAULTPULL_WRITE_TIMEOUT_SEC", "-1")
	t.Setenv("VAULTPULL_GLOBAL_TIMEOUT_SEC", "0")

	cfg := TimeoutConfigFromEnv()
	if cfg.VaultTimeout != 10*time.Second {
		t.Errorf("expected default 10s, got %v", cfg.VaultTimeout)
	}
	if cfg.WriteTimeout != 5*time.Second {
		t.Errorf("expected default 5s, got %v", cfg.WriteTimeout)
	}
	if cfg.GlobalTimeout != 30*time.Second {
		t.Errorf("expected default 30s, got %v", cfg.GlobalTimeout)
	}
}
