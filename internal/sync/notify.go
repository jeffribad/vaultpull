package sync

import (
	"fmt"
	"io"
	"os"
	"strings"
)

// NotifyChannel represents a supported notification channel.
type NotifyChannel string

const (
	NotifyStdout NotifyChannel = "stdout"
	NotifySlack  NotifyChannel = "slack"
)

// Notification holds the data for a sync notification.
type Notification struct {
	Role    string
	Path    string
	Written int
	DryRun  bool
	Errors  []string
}

// Notifier sends notifications after a sync operation.
type Notifier struct {
	channel NotifyChannel
	writer  io.Writer
	webhook string
}

// NewNotifier creates a Notifier based on channel config.
func NewNotifier(channel NotifyChannel, webhook string, w io.Writer) *Notifier {
	if w == nil {
		w = os.Stdout
	}
	return &Notifier{channel: channel, writer: w, webhook: webhook}
}

// Send dispatches the notification to the configured channel.
func (n *Notifier) Send(notif Notification) error {
	switch n.channel {
	case NotifySlack:
		return n.sendSlack(notif)
	default:
		return n.sendStdout(notif)
	}
}

func (n *Notifier) sendStdout(notif Notification) error {
	status := "synced"
	if notif.DryRun {
		status = "dry-run"
	}
	msg := fmt.Sprintf("[vaultpull] %s | role=%s path=%s keys=%d",
		status, notif.Role, notif.Path, notif.Written)
	if len(notif.Errors) > 0 {
		msg += " errors=" + strings.Join(notif.Errors, ";")
	}
	_, err := fmt.Fprintln(n.writer, msg)
	return err
}

func (n *Notifier) sendSlack(notif Notification) error {
	if n.webhook == "" {
		return fmt.Errorf("slack webhook URL is not configured")
	}
	// Real implementation would POST to n.webhook.
	// Stubbed here to keep external dependencies minimal.
	return n.sendStdout(notif)
}
