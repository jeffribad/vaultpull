package sync

import (
	"os"
	"strings"
)

// NotifyConfig holds notification settings loaded from the environment.
type NotifyConfig struct {
	Enabled  bool
	Channel  NotifyChannel
	Webhook  string
}

// NotifyConfigFromEnv reads notification config from environment variables.
//
//	VAULTPULL_NOTIFY_ENABLED=true
//	VAULTPULL_NOTIFY_CHANNEL=slack
//	VAULTPULL_NOTIFY_WEBHOOK=https://hooks.slack.com/...
func NotifyConfigFromEnv() NotifyConfig {
	cfg := NotifyConfig{
		Channel: NotifyStdout,
	}

	if v := os.Getenv("VAULTPULL_NOTIFY_ENABLED"); strings.EqualFold(v, "true") || v == "1" {
		cfg.Enabled = true
	}

	if ch := os.Getenv("VAULTPULL_NOTIFY_CHANNEL"); ch != "" {
		cfg.Channel = NotifyChannel(strings.ToLower(ch))
	}

	if wh := os.Getenv("VAULTPULL_NOTIFY_WEBHOOK"); wh != "" {
		cfg.Webhook = wh
	}

	return cfg
}
