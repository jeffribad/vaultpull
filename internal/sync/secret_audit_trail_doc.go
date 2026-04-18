// Package sync provides the secret_audit_trail module for vaultpull.
//
// The audit trail feature records every secret key accessed during a sync
// operation, including the action taken (read, skipped, renamed) and the
// role that triggered it. Entries are appended to a log file in either
// plain text (tab-separated) or JSON format.
//
// Configuration is driven by environment variables:
//
//	VAULTPULL_AUDIT_TRAIL_ENABLED  - set to "true" or "1" to enable (default: false)
//	VAULTPULL_AUDIT_TRAIL_FILE     - path to the audit log file (default: .vaultpull_audit.log)
//	VAULTPULL_AUDIT_TRAIL_FORMAT   - "text" or "json" (default: text)
//
// The audit log is append-only and uses file mode 0600 to restrict access.
// It is intended for compliance and debugging purposes and does NOT log
// secret values — only key names, actions, roles, and timestamps.
package sync
