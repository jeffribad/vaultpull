package sync

import (
	"os"
	"testing"
)

func TestPinConfigFromEnv_Defaults(t *testing.T) {
	os.Unsetenv("VAULTPULL_PIN_ENABLED")
	os.Unsetenv("VAULTPULL_PIN_KEYS")
	os.Unsetenv("VAULTPULL_PIN_FAIL_FAST")
	cfg := PinConfigFromEnv()
	if cfg.Enabled {
		t.Error("expected Enabled=false")
	}
	if cfg.FailFast {
		t.Error("expected FailFast=false")
	}
	if len(cfg.Pins) != 0 {
		t.Errorf("expected empty pins, got %v", cfg.Pins)
	}
}

func TestPinConfigFromEnv_Enabled(t *testing.T) {
	t.Setenv("VAULTPULL_PIN_ENABLED", "true")
	cfg := PinConfigFromEnv()
	if !cfg.Enabled {
		t.Error("expected Enabled=true")
	}
}

func TestPinConfigFromEnv_ParsesPins(t *testing.T) {
	t.Setenv("VAULTPULL_PIN_ENABLED", "1")
	t.Setenv("VAULTPULL_PIN_KEYS", "API_KEY=sk-prod, DB_PASS=s3cr")
	cfg := PinConfigFromEnv()
	if len(cfg.Pins) != 2 {
		t.Fatalf("expected 2 pins, got %d", len(cfg.Pins))
	}
	if cfg.Pins["API_KEY"] != "sk-prod" {
		t.Errorf("unexpected pin for API_KEY: %q", cfg.Pins["API_KEY"])
	}
}

func TestPinConfigFromEnv_SkipsMalformed(t *testing.T) {
	t.Setenv("VAULTPULL_PIN_ENABLED", "1")
	t.Setenv("VAULTPULL_PIN_KEYS", "BADENTRY,API_KEY=sk-prod")
	cfg := PinConfigFromEnv()
	if len(cfg.Pins) != 1 {
		t.Errorf("expected 1 valid pin, got %d", len(cfg.Pins))
	}
}
