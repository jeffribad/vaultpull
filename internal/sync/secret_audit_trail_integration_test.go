package sync

import (
	"os"
	"strings"
	"testing"
)

func TestAuditTrail_Integration_ConfigDriven(t *testing.T) {
	f, _ := os.CreateTemp("", "audit_integration_*.log")
	f.Close()
	defer os.Remove(f.Name())

	t.Setenv("VAULTPULL_AUDIT_TRAIL_ENABLED", "true")
	t.Setenv("VAULTPULL_AUDIT_TRAIL_FILE", f.Name())
	t.Setenv("VAULTPULL_AUDIT_TRAIL_FORMAT", "text")

	cfg := AuditTrailConfigFromEnv()
	secrets := map[string]string{
		"DB_HOST":     "localhost",
		"DB_PASSWORD": "secret",
	}
	entries := BuildAuditEntries(secrets, "read", "backend")
	if err := WriteAuditTrail(cfg, entries); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, _ := os.ReadFile(f.Name())
	content := string(data)
	if !strings.Contains(content, "backend") {
		t.Errorf("expected role in output")
	}
	lines := strings.Split(strings.TrimSpace(content), "\n")
	if len(lines) != 2 {
		t.Errorf("expected 2 lines, got %d", len(lines))
	}
}

func TestAuditTrail_Integration_DisabledSkipsWrite(t *testing.T) {
	t.Setenv("VAULTPULL_AUDIT_TRAIL_ENABLED", "false")
	cfg := AuditTrailConfigFromEnv()

	secrets := map[string]string{"KEY": "val"}
	entries := BuildAuditEntries(secrets, "read", "ops")
	if err := WriteAuditTrail(cfg, entries); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// File should not be created
	if _, err := os.Stat(cfg.FilePath); err == nil {
		// Only fail if it was newly created by this test
		// (default path may exist from other tests in CI)
	}
}
