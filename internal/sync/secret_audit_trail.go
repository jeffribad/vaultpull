package sync

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// AuditTrailConfig controls whether an audit trail file is written after sync.
type AuditTrailConfig struct {
	Enabled  bool
	FilePath string
	Format   string // "text" or "json"
}

// AuditTrailConfigFromEnv loads audit trail config from environment variables.
func AuditTrailConfigFromEnv() AuditTrailConfig {
	enabled := false
	if v := os.Getenv("VAULTPULL_AUDIT_TRAIL_ENABLED"); v != "" {
		if b, err := strconv.ParseBool(v); err == nil {
			enabled = b
		} else if v == "1" {
			enabled = true
		}
	}
	path := os.Getenv("VAULTPULL_AUDIT_TRAIL_FILE")
	if path == "" {
		path = ".vaultpull_audit.log"
	}
	format := strings.ToLower(os.Getenv("VAULTPULL_AUDIT_TRAIL_FORMAT"))
	if format != "json" {
		format = "text"
	}
	return AuditTrailConfig{Enabled: enabled, FilePath: path, Format: format}
}

// AuditEntry represents a single secret access event.
type AuditEntry struct {
	Timestamp time.Time
	Key       string
	Action    string // "read", "skipped", "renamed"
	Role      string
}

// WriteAuditTrail appends audit entries to the configured file.
func WriteAuditTrail(cfg AuditTrailConfig, entries []AuditEntry) error {
	if !cfg.Enabled || len(entries) == 0 {
		return nil
	}
	f, err := os.OpenFile(cfg.FilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return fmt.Errorf("audit trail: open file: %w", err)
	}
	defer f.Close()

	for _, e := range entries {
		var line string
		if cfg.Format == "json" {
			line = fmt.Sprintf(`{"ts":%q,"key":%q,"action":%q,"role":%q}`,
				e.Timestamp.UTC().Format(time.RFC3339), e.Key, e.Action, e.Role)
		} else {
			line = fmt.Sprintf("%s\t%s\t%s\t%s",
				e.Timestamp.UTC().Format(time.RFC3339), e.Action, e.Key, e.Role)
		}
		if _, err := fmt.Fprintln(f, line); err != nil {
			return fmt.Errorf("audit trail: write: %w", err)
		}
	}
	return nil
}

// BuildAuditEntries creates AuditEntry records from a secrets map.
func BuildAuditEntries(secrets map[string]string, action, role string) []AuditEntry {
	entries := make([]AuditEntry, 0, len(secrets))
	now := time.Now()
	for k := range secrets {
		entries = append(entries, AuditEntry{Timestamp: now, Key: k, Action: action, Role: role})
	}
	return entries
}
