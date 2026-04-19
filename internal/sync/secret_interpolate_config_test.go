package sync

import (
	"testing"
)

func TestInterpolateConfigFromEnv_Defaults(t *testing.T) {
	t.Setenv("VAULTPULL_INTERPOLATE_ENABLED", "")
	t.Setenv("VAULTPULL_INTERPOLATE_ALLOW_ENV", "")
	cfg := InterpolateConfigFromEnv()
	if cfg.Enabled {
		t.Error("expected Enabled=false by default")
	}
	if cfg.AllowEnv {
		t.Error("expected AllowEnv=false by default")
	}
}

func TestInterpolateConfigFromEnv_Enabled(t *testing.T) {
	t.Setenv("VAULTPULL_INTERPOLATE_ENABLED", "true")
	cfg := InterpolateConfigFromEnv()
	if !cfg.Enabled {
		t.Error("expected Enabled=true")
	}
}

func TestInterpolateConfigFromEnv_AllowEnv(t *testing.T) {
	t.Setenv("VAULTPULL_INTERPOLATE_ENABLED", "true")
	t.Setenv("VAULTPULL_INTERPOLATE_ALLOW_ENV", "true")
	cfg := InterpolateConfigFromEnv()
	if !cfg.AllowEnv {
		t.Error("expected AllowEnv=true")
	}
}

func TestInterpolateConfigFromEnv_NumericEnabled(t *testing.T) {
	t.Setenv("VAULTPULL_INTERPOLATE_ENABLED", "1")
	cfg := InterpolateConfigFromEnv()
	if !cfg.Enabled {
		t.Error("expected Enabled=true for value '1'")
	}
}
