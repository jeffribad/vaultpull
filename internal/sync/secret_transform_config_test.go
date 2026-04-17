package sync

import (
	"os"
	"testing"
)

func TestTransformConfigFromEnv_Defaults(t *testing.T) {
	os.Unsetenv("VAULTPULL_KEY_PREFIX")
	os.Unsetenv("VAULTPULL_KEY_SUFFIX")
	os.Unsetenv("VAULTPULL_KEY_UPPERCASE")
	os.Unsetenv("VAULTPULL_KEY_LOWERCASE")

	cfg := TransformConfigFromEnv()
	if cfg.Prefix != "" || cfg.Suffix != "" || cfg.UpperCase || cfg.LowerCase {
		t.Errorf("expected zero-value defaults, got %+v", cfg)
	}
}

func TestTransformConfigFromEnv_Prefix(t *testing.T) {
	t.Setenv("VAULTPULL_KEY_PREFIX", "APP_")
	cfg := TransformConfigFromEnv()
	if cfg.Prefix != "APP_" {
		t.Errorf("expected prefix APP_, got %q", cfg.Prefix)
	}
}

func TestTransformConfigFromEnv_UpperCase(t *testing.T) {
	t.Setenv("VAULTPULL_KEY_UPPERCASE", "1")
	cfg := TransformConfigFromEnv()
	if !cfg.UpperCase {
		t.Error("expected UpperCase to be true")
	}
}
