package sync

import (
	"testing"
)

func TestNotifyConfigFromEnv_Defaults(t *testing.T) {
	t.Setenv("VAULTPULL_NOTIFY_ENABLED", "")
	t.Setenv("VAULTPULL_NOTIFY_CHANNEL", "")
	t.Setenv("VAULTPULL_NOTIFY_WEBHOOK", "")

	cfg := NotifyConfigFromEnv()
	if cfg.Enabled {
		t.Error("expected Enabled=false by default")
	}
	if cfg.Channel != NotifyStdout {
		t.Errorf("expected channel stdout, got %s", cfg.Channel)
	}
	if cfg.Webhook != "" {
		t.Errorf("expected empty webhook, got %s", cfg.Webhook)
	}
}

func TestNotifyConfigFromEnv_Enabled(t *testing.T) {
	t.Setenv("VAULTPULL_NOTIFY_ENABLED", "true")
	cfg := NotifyConfigFromEnv()
	if !cfg.Enabled {
		t.Error("expected Enabled=true")
	}
}

func TestNotifyConfigFromEnv_EnabledNumeric(t *testing.T) {
	t.Setenv("VAULTPULL_NOTIFY_ENABLED", "1")
	cfg := NotifyConfigFromEnv()
	if !cfg.Enabled {
		t.Error("expected Enabled=true for value '1'")
	}
}

func TestNotifyConfigFromEnv_SlackChannel(t *testing.T) {
	t.Setenv("VAULTPULL_NOTIFY_CHANNEL", "SLACK")
	t.Setenv("VAULTPULL_NOTIFY_WEBHOOK", "https://hooks.slack.com/test")
	cfg := NotifyConfigFromEnv()
	if cfg.Channel != NotifySlack {
		t.Errorf("expected slack channel, got %s", cfg.Channel)
	}
	if cfg.Webhook != "https://hooks.slack.com/test" {
		t.Errorf("unexpected webhook: %s", cfg.Webhook)
	}
}
