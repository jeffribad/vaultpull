package sync

import (
	"os"
	"testing"
)

func TestAuditTrailConfigFromEnv_Defaults(t *testing.T) {
	os.Unsetenv("VAULTPULL_AUDIT_TRAIL_ENABLED")
	os.Unsetenv("VAULTPULL_AUDIT_TRAIL_FILE")
	os.Unsetenv("VAULTPULL_AUDIT_TRAIL_FORMAT")

	cfg := AuditTrailConfigFromEnv()
	if cfg.Enabled {
		t.Error("expected disabled by default")
	}
	if cfg.FilePath != ".vaultpull_audit.log" {
		t.Errorf("unexpected default path: %s", cfg.FilePath)
	}
	if cfg.Format != "text" {
		t.Errorf("unexpected default format: %s", cfg.Format)
	}
}

func TestAuditTrailConfigFromEnv_Enabled(t *testing.T) {
	t.Setenv("VAULTPULL_AUDIT_TRAIL_ENABLED", "true")
	t.Setenv("VAULTPULL_AUDIT_TRAIL_FILE", "/tmp/audit.log")
	t.Setenv("VAULTPULL_AUDIT_TRAIL_FORMAT", "json")

	cfg := AuditTrailConfigFromEnv()
	if !cfg.Enabled {
		t.Error("expected enabled")
	}
	if cfg.FilePath != "/tmp/audit.log" {
		t.Errorf("unexpected path: %s", cfg.FilePath)
	}
	if cfg.Format != "json" {
		t.Errorf("unexpected format: %s", cfg.Format)
	}
}

func TestAuditTrailConfigFromEnv_NumericEnabled(t *testing.T) {
	t.Setenv("VAULTPULL_AUDIT_TRAIL_ENABLED", "1")
	cfg := AuditTrailConfigFromEnv()
	if !cfg.Enabled {
		t.Error("expected enabled via numeric 1")
	}
}
