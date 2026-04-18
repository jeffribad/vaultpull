package sync

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// SnapshotConfig controls whether secret snapshots are saved before sync.
type SnapshotConfig struct {
	Enabled bool
	Dir     string
}

// SnapshotConfigFromEnv loads snapshot config from environment variables.
func SnapshotConfigFromEnv() SnapshotConfig {
	enabled := os.Getenv("VAULTPULL_SNAPSHOT_ENABLED")
	dir := os.Getenv("VAULTPULL_SNAPSHOT_DIR")
	if dir == "" {
		dir = ".vaultpull/snapshots"
	}
	return SnapshotConfig{
		Enabled: enabled == "true" || enabled == "1",
		Dir:     dir,
	}
}

// Snapshot represents a point-in-time capture of secrets.
type Snapshot struct {
	Timestamp time.Time         `json:"timestamp"`
	Role      string            `json:"role"`
	Secrets   map[string]string `json:"secrets"`
}

// SaveSnapshot writes a JSON snapshot of secrets to the configured directory.
func SaveSnapshot(cfg SnapshotConfig, role string, secrets map[string]string) (string, error) {
	if !cfg.Enabled {
		return "", nil
	}
	if err := os.MkdirAll(cfg.Dir, 0700); err != nil {
		return "", fmt.Errorf("snapshot: create dir: %w", err)
	}
	snap := Snapshot{
		Timestamp: time.Now().UTC(),
		Role:      role,
		Secrets:   secrets,
	}
	data, err := json.MarshalIndent(snap, "", "  ")
	if err != nil {
		return "", fmt.Errorf("snapshot: marshal: %w", err)
	}
	filename := fmt.Sprintf("%s_%s.json", snap.Timestamp.Format("20060102T150405"), role)
	path := filepath.Join(cfg.Dir, filename)
	if err := os.WriteFile(path, data, 0600); err != nil {
		return "", fmt.Errorf("snapshot: write: %w", err)
	}
	return path, nil
}
