package sync

import (
	"os"
	"testing"
)

func TestPromoteConfigFromEnv_Defaults(t *testing.T) {
	os.Unsetenv("VAULTPULL_PROMOTE_ENABLED")
	os.Unsetenv("VAULTPULL_PROMOTE_FROM_PREFIX")
	os.Unsetenv("VAULTPULL_PROMOTE_TO_PREFIX")
	os.Unsetenv("VAULTPULL_PROMOTE_OVERWRITE")

	cfg := PromoteConfigFromEnv()
	if cfg.Enabled {
		t.Error("expected Enabled=false by default")
	}
	if cfg.FromPrefix != "" {
		t.Errorf("expected empty FromPrefix, got %q", cfg.FromPrefix)
	}
	if cfg.ToPrefix != "" {
		t.Errorf("expected empty ToPrefix, got %q", cfg.ToPrefix)
	}
	if cfg.Overwrite {
		t.Error("expected Overwrite=false by default")
	}
}

func TestPromoteConfigFromEnv_Enabled(t *testing.T) {
	t.Setenv("VAULTPULL_PROMOTE_ENABLED", "true")
	t.Setenv("VAULTPULL_PROMOTE_FROM_PREFIX", "STAGING_")
	t.Setenv("VAULTPULL_PROMOTE_TO_PREFIX", "PROD_")
	t.Setenv("VAULTPULL_PROMOTE_OVERWRITE", "true")

	cfg := PromoteConfigFromEnv()
	if !cfg.Enabled {
		t.Error("expected Enabled=true")
	}
	if cfg.FromPrefix != "STAGING_" {
		t.Errorf("expected FromPrefix=STAGING_, got %q", cfg.FromPrefix)
	}
	if cfg.ToPrefix != "PROD_" {
		t.Errorf("expected ToPrefix=PROD_, got %q", cfg.ToPrefix)
	}
	if !cfg.Overwrite {
		t.Error("expected Overwrite=true")
	}
}

func TestPromoteConfigFromEnv_NumericEnabled(t *testing.T) {
	t.Setenv("VAULTPULL_PROMOTE_ENABLED", "1")
	cfg := PromoteConfigFromEnv()
	if !cfg.Enabled {
		t.Error("expected Enabled=true for value '1'")
	}
}
