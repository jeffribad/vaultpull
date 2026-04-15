package sync_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/vaultpull/internal/sync"
)

// TestRollback_Integration_BackupRestoredOnWriteFailure verifies that when a
// write fails after a backup has been created, the original file is restored.
func TestRollback_Integration_BackupRestoredOnWriteFailure(t *testing.T) {
	dir := t.TempDir()
	envFile := filepath.Join(dir, ".env")
	original := "ORIGINAL_KEY=original_value\n"

	if err := os.WriteFile(envFile, []byte(original), 0600); err != nil {
		t.Fatal(err)
	}

	bak, err := sync.CreateBackup(envFile)
	if err != nil {
		t.Fatalf("backup failed: %v", err)
	}
	if bak == nil {
		t.Fatal("expected non-nil backup")
	}

	// Simulate a failed write by corrupting the file then restoring.
	if err := os.WriteFile(envFile, []byte("CORRUPT=data\n"), 0600); err != nil {
		t.Fatal(err)
	}

	if err := bak.Restore(); err != nil {
		t.Fatalf("restore failed: %v", err)
	}

	got, err := os.ReadFile(envFile)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != original {
		t.Errorf("expected restored content %q, got %q", original, string(got))
	}

	// Backup file must be gone.
	if _, err := os.Stat(bak.BackupPath); !os.IsNotExist(err) {
		t.Error("backup file should not exist after successful restore")
	}
}

// TestRollback_Integration_DiscardLeavesOriginalIntact ensures Discard removes
// only the backup without touching the original file.
func TestRollback_Integration_DiscardLeavesOriginalIntact(t *testing.T) {
	dir := t.TempDir()
	envFile := filepath.Join(dir, ".env")
	content := []byte("KEY=keep_me\n")

	if err := os.WriteFile(envFile, content, 0600); err != nil {
		t.Fatal(err)
	}

	bak, err := sync.CreateBackup(envFile)
	if err != nil || bak == nil {
		t.Fatalf("backup setup failed: %v", err)
	}

	if err := bak.Discard(); err != nil {
		t.Fatalf("discard failed: %v", err)
	}

	got, err := os.ReadFile(envFile)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != string(content) {
		t.Errorf("original file modified unexpectedly: got %q", got)
	}
}
