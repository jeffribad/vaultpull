package audit_test

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/yourusername/vaultpull/internal/audit"
)

// TestLogger_MultipleEvents verifies sequential events are each on their own line.
func TestLogger_MultipleEvents(t *testing.T) {
	var buf bytes.Buffer
	l := audit.NewLogger(&buf)

	events := []struct {
		role string
		path string
		keys []string
	}{
		{"backend", "secret/backend", []string{"DB_HOST", "DB_PASS"}},
		{"frontend", "secret/frontend", []string{"REACT_APP_URL"}},
		{"ops", "secret/ops", []string{"SSH_KEY", "DEPLOY_TOKEN", "SLACK_WEBHOOK"}},
	}

	for _, ev := range events {
		if err := l.LogSync(ev.role, ev.path, ev.keys); err != nil {
			t.Fatalf("LogSync(%q): %v", ev.role, err)
		}
	}

	lines := bytes.Split(bytes.TrimRight(buf.Bytes(), "\n"), []byte("\n"))
	if len(lines) != len(events) {
		t.Fatalf("expected %d lines, got %d", len(events), len(lines))
	}

	for i, line := range lines {
		var e audit.Event
		if err := json.Unmarshal(line, &e); err != nil {
			t.Errorf("line %d is not valid JSON: %v", i, err)
			continue
		}
		if e.Role != events[i].role {
			t.Errorf("line %d: expected role=%q, got %q", i, events[i].role, e.Role)
		}
		if len(e.Keys) != len(events[i].keys) {
			t.Errorf("line %d: expected %d keys, got %d", i, len(events[i].keys), len(e.Keys))
		}
	}
}
