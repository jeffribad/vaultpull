package sync

import (
	"os"
	"strings"
	"testing"
	"time"
)

func TestWriteAuditTrail_Disabled_ReturnsNil(t *testing.T) {
	cfg := AuditTrailConfig{Enabled: false}
	err := WriteAuditTrail(cfg, []AuditEntry{{Key: "FOO", Action: "read"}})
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestWriteAuditTrail_EmptyEntries_ReturnsNil(t *testing.T) {
	cfg := AuditTrailConfig{Enabled: true, FilePath: "/tmp/test_audit_empty.log", Format: "text"}
	err := WriteAuditTrail(cfg, nil)
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestWriteAuditTrail_TextFormat_WritesLines(t *testing.T) {
	f, _ := os.CreateTemp("", "audit_text_*.log")
	f.Close()
	defer os.Remove(f.Name())

	cfg := AuditTrailConfig{Enabled: true, FilePath: f.Name(), Format: "text"}
	entries := []AuditEntry{
		{Timestamp: time.Now(), Key: "DB_PASS", Action: "read", Role: "backend"},
	}
	if err := WriteAuditTrail(cfg, entries); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	data, _ := os.ReadFile(f.Name())
	if !strings.Contains(string(data), "DB_PASS") {
		t.Errorf("expected DB_PASS in output, got: %s", data)
	}
	if !strings.Contains(string(data), "read") {
		t.Errorf("expected action in output")
	}
}

func TestWriteAuditTrail_JSONFormat_WritesJSON(t *testing.T) {
	f, _ := os.CreateTemp("", "audit_json_*.log")
	f.Close()
	defer os.Remove(f.Name())

	cfg := AuditTrailConfig{Enabled: true, FilePath: f.Name(), Format: "json"}
	entries := []AuditEntry{
		{Timestamp: time.Now(), Key: "API_KEY", Action: "skipped", Role: "ops"},
	}
	if err := WriteAuditTrail(cfg, entries); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	data, _ := os.ReadFile(f.Name())
	if !strings.Contains(string(data), `"key"`) {
		t.Errorf("expected JSON key field, got: %s", data)
	}
}

func TestBuildAuditEntries_ReturnsOnePerSecret(t *testing.T) {
	secrets := map[string]string{"A": "1", "B": "2"}
	entries := BuildAuditEntries(secrets, "read", "backend")
	if len(entries) != 2 {
		t.Errorf("expected 2 entries, got %d", len(entries))
	}
	for _, e := range entries {
		if e.Action != "read" || e.Role != "backend" {
			t.Errorf("unexpected entry: %+v", e)
		}
	}
}
