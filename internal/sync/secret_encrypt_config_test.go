package sync

import (
	"os"
	"testing"
)

func TestEncryptConfigFromEnv_Defaults(t *testing.T) {
	os.Unsetenv("VAULTPULL_ENCRYPT_OUTPUT")
	os.Unsetenv("VAULTPULL_ENCRYPT_KEY")
	cfg := EncryptConfigFromEnv()
	if cfg.Enabled {
		t.Error("expected Enabled=false by default")
	}
	if cfg.Key != "" {
		t.Error("expected empty key by default")
	}
}

func TestEncryptConfigFromEnv_Enabled(t *testing.T) {
	os.Setenv("VAULTPULL_ENCRYPT_OUTPUT", "true")
	os.Setenv("VAULTPULL_ENCRYPT_KEY", "mykey")
	defer os.Unsetenv("VAULTPULL_ENCRYPT_OUTPUT")
	defer os.Unsetenv("VAULTPULL_ENCRYPT_KEY")
	cfg := EncryptConfigFromEnv()
	if !cfg.Enabled {
		t.Error("expected Enabled=true")
	}
	if cfg.Key != "mykey" {
		t.Errorf("expected key=mykey, got %q", cfg.Key)
	}
}

func TestEncryptConfigFromEnv_NumericEnabled(t *testing.T) {
	os.Setenv("VAULTPULL_ENCRYPT_OUTPUT", "1")
	defer os.Unsetenv("VAULTPULL_ENCRYPT_OUTPUT")
	cfg := EncryptConfigFromEnv()
	if !cfg.Enabled {
		t.Error("expected Enabled=true for value '1'")
	}
}
