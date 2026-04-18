package sync

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestSnapshotConfigFromEnv_Defaults(t *testing.T) {
	os.Unsetenv("VAULTPULL_SNAPSHOT_ENABLED")
	os.Unsetenv("VAULTPULL_SNAPSHOT_DIR")
	cfg := SnapshotConfigFromEnv()
	if cfg.Enabled {
		t.Error("expected disabled by default")
	}
	if cfg.Dir != ".vaultpull/snapshots" {
		t.Errorf("unexpected default dir: %s", cfg.Dir)
	}
}

func TestSnapshotConfigFromEnv_Enabled(t *testing.T) {
	t.Setenv("VAULTPULL_SNAPSHOT_ENABLED", "true")
	t.Setenv("VAULTPULL_SNAPSHOT_DIR", "/tmp/snaps")
	cfg := SnapshotConfigFromEnv()
	if !cfg.Enabled {
		t.Error("expected enabled")
	}
	if cfg.Dir != "/tmp/snaps" {
		t.Errorf("unexpected dir: %s", cfg.Dir)
	}
}

func TestSaveSnapshot_Disabled_ReturnsEmpty(t *testing.T) {
	cfg := SnapshotConfig{Enabled: false}
	path, err := SaveSnapshot(cfg, "backend", map[string]string{"KEY": "val"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if path != "" {
		t.Errorf("expected empty path, got %s", path)
	}
}

func TestSaveSnapshot_WritesFile(t *testing.T) {
	dir := t.TempDir()
	cfg := SnapshotConfig{Enabled: true, Dir: dir}
	secrets := map[string]string{"DB_HOST": "localhost", "API_KEY": "secret"}
	path, err := SaveSnapshot(cfg, "backend", secrets)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if path == "" {
		t.Fatal("expected non-empty path")
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("cannot read snapshot: %v", err)
	}
	var snap Snapshot
	if err := json.Unmarshal(data, &snap); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if snap.Role != "backend" {
		t.Errorf("expected role backend, got %s", snap.Role)
	}
	if snap.Secrets["DB_HOST"] != "localhost" {
		t.Errorf("missing secret DB_HOST")
	}
}

func TestSaveSnapshot_CreatesDir(t *testing.T) {
	base := t.TempDir()
	dir := filepath.Join(base, "nested", "snaps")
	cfg := SnapshotConfig{Enabled: true, Dir: dir}
	_, err := SaveSnapshot(cfg, "ops", map[string]string{"X": "y"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		t.Error("expected directory to be created")
	}
}
