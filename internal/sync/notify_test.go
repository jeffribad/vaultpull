package sync

import (
	"bytes"
	"strings"
	"testing"
)

func TestNotifier_Stdout_Synced(t *testing.T) {
	var buf bytes.Buffer
	n := NewNotifier(NotifyStdout, "", &buf)
	err := n.Send(Notification{Role: "backend", Path: "secret/app", Written: 5})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "synced") {
		t.Errorf("expected 'synced' in output, got: %s", out)
	}
	if !strings.Contains(out, "keys=5") {
		t.Errorf("expected keys=5 in output, got: %s", out)
	}
}

func TestNotifier_Stdout_DryRun(t *testing.T) {
	var buf bytes.Buffer
	n := NewNotifier(NotifyStdout, "", &buf)
	err := n.Send(Notification{Role: "ops", Path: "secret/ops", Written: 3, DryRun: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "dry-run") {
		t.Errorf("expected 'dry-run' in output")
	}
}

func TestNotifier_Stdout_WithErrors(t *testing.T) {
	var buf bytes.Buffer
	n := NewNotifier(NotifyStdout, "", &buf)
	err := n.Send(Notification{Role: "frontend", Path: "secret/fe", Errors: []string{"timeout"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "errors=timeout") {
		t.Errorf("expected errors in output, got: %s", buf.String())
	}
}

func TestNotifier_Slack_NoWebhook(t *testing.T) {
	var buf bytes.Buffer
	n := NewNotifier(NotifySlack, "", &buf)
	err := n.Send(Notification{Role: "backend", Path: "secret/app", Written: 2})
	if err == nil {
		t.Fatal("expected error for missing webhook")
	}
}

func TestNotifier_Slack_WithWebhook(t *testing.T) {
	var buf bytes.Buffer
	n := NewNotifier(NotifySlack, "https://hooks.slack.com/stub", &buf)
	err := n.Send(Notification{Role: "backend", Path: "secret/app", Written: 2})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNewNotifier_NilWriter_DefaultsToStdout(t *testing.T) {
	n := NewNotifier(NotifyStdout, "", nil)
	if n.writer == nil {
		t.Error("expected non-nil writer")
	}
}
