package sync

import (
	"os"
	"testing"
)

func TestEnvMapConfigFromEnv_Defaults(t *testing.T) {
	os.Unsetenv("VAULTPULL_ENVMAP_ENABLED")
	os.Unsetenv("VAULTPULL_ENVMAP_KEYS")

	cfg := EnvMapConfigFromEnv()
	if cfg.Enabled {
		t.Error("expected Enabled=false by default")
	}
	if len(cfg.Mappings) != 0 {
		t.Errorf("expected empty mappings, got %v", cfg.Mappings)
	}
}

func TestEnvMapConfigFromEnv_Enabled(t *testing.T) {
	t.Setenv("VAULTPULL_ENVMAP_ENABLED", "true")
	os.Unsetenv("VAULTPULL_ENVMAP_KEYS")

	cfg := EnvMapConfigFromEnv()
	if !cfg.Enabled {
		t.Error("expected Enabled=true")
	}
}

func TestEnvMapConfigFromEnv_NumericEnabled(t *testing.T) {
	t.Setenv("VAULTPULL_ENVMAP_ENABLED", "1")
	cfg := EnvMapConfigFromEnv()
	if !cfg.Enabled {
		t.Error("expected Enabled=true for value '1'")
	}
}

func TestEnvMapConfigFromEnv_ParsesMappings(t *testing.T) {
	t.Setenv("VAULTPULL_ENVMAP_ENABLED", "true")
	t.Setenv("VAULTPULL_ENVMAP_KEYS", "DB_PASS:DATABASE_PASSWORD,API_KEY:SERVICE_API_KEY")

	cfg := EnvMapConfigFromEnv()
	if cfg.Mappings["DB_PASS"] != "DATABASE_PASSWORD" {
		t.Errorf("unexpected mapping: %v", cfg.Mappings)
	}
	if cfg.Mappings["API_KEY"] != "SERVICE_API_KEY" {
		t.Errorf("unexpected mapping: %v", cfg.Mappings)
	}
}

func TestEnvMapConfigFromEnv_SkipsMalformed(t *testing.T) {
	t.Setenv("VAULTPULL_ENVMAP_ENABLED", "true")
	t.Setenv("VAULTPULL_ENVMAP_KEYS", "GOOD_KEY:GOOD_ENV,BADENTRY,ANOTHER:VALID")

	cfg := EnvMapConfigFromEnv()
	if len(cfg.Mappings) != 2 {
		t.Errorf("expected 2 valid mappings, got %d", len(cfg.Mappings))
	}
}
