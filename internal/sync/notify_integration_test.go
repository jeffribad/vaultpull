package sync

import (
	"bytes"
	"strings"
	"testing"
)

// TestNotify_Integration_ConfigDrivenNotifier verifies that a Notifier built
// from NotifyConfigFromEnv sends a correctly formatted message.
func TestNotify_Integration_ConfigDrivenNotifier(t *testing.T) {
	t.Setenv("VAULTPULL_NOTIFY_ENABLED", "true")
	t.Setenv("VAULTPULL_NOTIFY_CHANNEL", "stdout")
	t.Setenv("VAULTPULL_NOTIFY_WEBHOOK", "")

	cfg := NotifyConfigFromEnv()
	if !cfg.Enabled {
		t.Fatal("expected notifications enabled")
	}

	var buf bytes.Buffer
	n := NewNotifier(cfg.Channel, cfg.Webhook, &buf)

	err := n.Send(Notification{
		Role:    "backend",
		Path:    "secret/data/app",
		Written: 7,
		DryRun:  false,
	})
	if err != nil {
		t.Fatalf("unexpected send error: %v", err)
	}

	out := buf.String()
	for _, want := range []string{"synced", "backend", "keys=7"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected %q in output: %s", want, out)
		}
	}
}

// TestNotify_Integration_DisabledSkipsSend verifies callers can gate on Enabled.
func TestNotify_Integration_DisabledSkipsSend(t *testing.T) {
	t.Setenv("VAULTPULL_NOTIFY_ENABLED", "false")

	cfg := NotifyConfigFromEnv()
	if cfg.Enabled {
		t.Fatal("expected notifications disabled")
	}
	// No send call — just confirm the config gate works.
}
