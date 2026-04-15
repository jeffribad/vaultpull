package audit_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/yourusername/vaultpull/internal/audit"
)

func TestLogger_LogSync_WritesJSON(t *testing.T) {
	var buf bytes.Buffer
	l := audit.NewLogger(&buf)

	err := l.LogSync("backend", "secret/app", []string{"DB_URL", "API_KEY"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var e audit.Event
	if err := json.Unmarshal(buf.Bytes(), &e); err != nil {
		t.Fatalf("output is not valid JSON: %v\noutput: %s", err, buf.String())
	}

	if e.Action != "sync" {
		t.Errorf("expected action=sync, got %q", e.Action)
	}
	if e.Role != "backend" {
		t.Errorf("expected role=backend, got %q", e.Role)
	}
	if e.Path != "secret/app" {
		t.Errorf("expected path=secret/app, got %q", e.Path)
	}
	if len(e.Keys) != 2 {
		t.Errorf("expected 2 keys, got %d", len(e.Keys))
	}
	if e.Timestamp.IsZero() {
		t.Error("expected non-zero timestamp")
	}
}

func TestLogger_LogError_WritesMessage(t *testing.T) {
	var buf bytes.Buffer
	l := audit.NewLogger(&buf)

	err := l.LogError("fetch", "vault unreachable")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var e audit.Event
	if err := json.Unmarshal(buf.Bytes(), &e); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}

	if e.Action != "fetch" {
		t.Errorf("expected action=fetch, got %q", e.Action)
	}
	if e.Message != "vault unreachable" {
		t.Errorf("expected message, got %q", e.Message)
	}
	if e.Role != "" || e.Path != "" {
		t.Error("expected empty role and path for error event")
	}
}

func TestLogger_NilWriter_DefaultsToStderr(t *testing.T) {
	// Just ensure no panic when nil writer is passed.
	l := audit.NewLogger(nil)
	if l == nil {
		t.Fatal("expected non-nil logger")
	}
}

func TestLogger_OutputEndsWithNewline(t *testing.T) {
	var buf bytes.Buffer
	l := audit.NewLogger(&buf)
	_ = l.LogSync("ops", "secret/ops", []string{"TOKEN"})

	if !strings.HasSuffix(buf.String(), "\n") {
		t.Error("expected output to end with newline")
	}
}
