package sync

import (
	"testing"
)

func TestApplyPassthrough_Disabled_ReturnsOriginal(t *testing.T) {
	secrets := map[string]string{"DB_HOST": "vault-host"}
	t.Setenv("DB_HOST", "local-host")

	cfg := PassthroughConfig{Enabled: false}
	result := ApplyPassthrough(cfg, secrets)

	if result["DB_HOST"] != "vault-host" {
		t.Errorf("expected vault-host, got %s", result["DB_HOST"])
	}
}

func TestApplyPassthrough_AllKeys_OverridesFromEnv(t *testing.T) {
	secrets := map[string]string{"DB_HOST": "vault-host", "APP_PORT": "8080"}
	t.Setenv("DB_HOST", "local-host")

	cfg := PassthroughConfig{Enabled: true}
	result := ApplyPassthrough(cfg, secrets)

	if result["DB_HOST"] != "local-host" {
		t.Errorf("expected local-host, got %s", result["DB_HOST"])
	}
	if result["APP_PORT"] != "8080" {
		t.Errorf("expected 8080, got %s", result["APP_PORT"])
	}
}

func TestApplyPassthrough_SpecificKeys_OnlyOverridesListed(t *testing.T) {
	secrets := map[string]string{"DB_HOST": "vault-host", "DB_PASS": "vault-pass"}
	t.Setenv("DB_HOST", "local-host")
	t.Setenv("DB_PASS", "local-pass")

	cfg := PassthroughConfig{Enabled: true, Keys: []string{"DB_HOST"}}
	result := ApplyPassthrough(cfg, secrets)

	if result["DB_HOST"] != "local-host" {
		t.Errorf("expected local-host, got %s", result["DB_HOST"])
	}
	if result["DB_PASS"] != "vault-pass" {
		t.Errorf("expected vault-pass, got %s", result["DB_PASS"])
	}
}

func TestApplyPassthrough_DoesNotMutateInput(t *testing.T) {
	secrets := map[string]string{"KEY": "original"}
	t.Setenv("KEY", "override")

	cfg := PassthroughConfig{Enabled: true}
	ApplyPassthrough(cfg, secrets)

	if secrets["KEY"] != "original" {
		t.Error("input map was mutated")
	}
}

func TestPassthroughConfigFromEnv_Defaults(t *testing.T) {
	t.Setenv("VAULTPULL_PASSTHROUGH_ENABLED", "")
	t.Setenv("VAULTPULL_PASSTHROUGH_KEYS", "")

	cfg := PassthroughConfigFromEnv()

	if cfg.Enabled {
		t.Error("expected Enabled=false by default")
	}
	if len(cfg.Keys) != 0 {
		t.Errorf("expected no keys, got %v", cfg.Keys)
	}
}

func TestPassthroughConfigFromEnv_Enabled(t *testing.T) {
	t.Setenv("VAULTPULL_PASSTHROUGH_ENABLED", "true")
	t.Setenv("VAULTPULL_PASSTHROUGH_KEYS", "FOO, BAR")

	cfg := PassthroughConfigFromEnv()

	if !cfg.Enabled {
		t.Error("expected Enabled=true")
	}
	if len(cfg.Keys) != 2 || cfg.Keys[0] != "FOO" || cfg.Keys[1] != "BAR" {
		t.Errorf("unexpected keys: %v", cfg.Keys)
	}
}
