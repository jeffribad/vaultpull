package audit

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"
)

// Event represents a single audit log entry.
type Event struct {
	Timestamp time.Time `json:"timestamp"`
	Action    string    `json:"action"`
	Role      string    `json:"role,omitempty"`
	Path      string    `json:"path,omitempty"`
	Keys      []string  `json:"keys,omitempty"`
	Message   string    `json:"message,omitempty"`
}

// Logger writes structured audit events.
type Logger struct {
	out io.Writer
}

// NewLogger creates a Logger writing to the given writer.
// If w is nil, os.Stderr is used.
func NewLogger(w io.Writer) *Logger {
	if w == nil {
		w = os.Stderr
	}
	return &Logger{out: w}
}

// Log emits a JSON-encoded audit event.
func (l *Logger) Log(action, role, path string, keys []string, msg string) error {
	e := Event{
		Timestamp: time.Now().UTC(),
		Action:    action,
		Role:      role,
		Path:      path,
		Keys:      keys,
		Message:   msg,
	}
	b, err := json.Marshal(e)
	if err != nil {
		return fmt.Errorf("audit: marshal event: %w", err)
	}
	_, err = fmt.Fprintln(l.out, string(b))
	return err
}

// LogSync is a convenience wrapper for a secrets-sync event.
func (l *Logger) LogSync(role, path string, keys []string) error {
	return l.Log("sync", role, path, keys, "")
}

// LogError records an error event without key details.
func (l *Logger) LogError(action, msg string) error {
	return l.Log(action, "", "", nil, msg)
}
