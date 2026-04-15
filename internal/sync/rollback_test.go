package sync

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCreateBackup_NonExistentFile(t *testing.T) {
	bak, err := CreateBackup("/tmp/vaultpull_no_such_file_xyz.env")
	if err != nil {
		t.Fatalf("expected no error for missing file, got: %v", err)
	}
	if bak != nil {
		t.Errorf("expected nil backup for missing file, got: %+v", bak)
	}
}

func TestCreateBackup_CreatesBackupFile(t *testing.T) {
	dir := t.TempDir()
	original := filepath.Join(dir, ".env")
	content := []byte("KEY=value\nFOO=bar\n")
	if err := os.WriteFile(original, content, 0600); err != nil {
		t.Fatal(err)
	}

	bak, err := CreateBackup(original)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if bak == nil {
		t.Fatal("expected backup, got nil")
	}

	if !strings.HasSuffix(bak.BackupPath, ".bak") {
		t.Errorf("backup path should end with .bak, got: %s", bak.BackupPath)
	}

	got, err := os.ReadFile(bak.BackupPath)
	if err != nil {
		t.Fatalf("could not read backup: %v", err)
	}
	if string(got) != string(content) {
		t.Errorf("backup content mismatch: got %q, want %q", got, content)
	}
}

func TestBackup_Restore(t *testing.T) {
	dir := t.TempDir()
	original := filepath.Join(dir, ".env")
	origContent := []byte("ORIGINAL=yes\n")
	if err := os.WriteFile(original, origContent, 0600); err != nil {
		t.Fatal(err)
	}

	bak, err := CreateBackup(original)
	if err != nil || bak == nil {
		t.Fatalf("setup failed: %v", err)
	}

	// Overwrite original
	if err := os.WriteFile(original, []byte("MODIFIED=yes\n"), 0600); err != nil {
		t.Fatal(err)
	}

	if err := bak.Restore(); err != nil {
		t.Fatalf("restore failed: %v", err)
	}

	got, _ := os.ReadFile(original)
	if string(got) != string(origContent) {
		t.Errorf("restored content mismatch: got %q, want %q", got, origContent)
	}

	if _, err := os.Stat(bak.BackupPath); !os.IsNotExist(err) {
		t.Error("backup file should have been removed after restore")
	}
}

func TestBackup_Discard(t *testing.T) {
	dir := t.TempDir()
	original := filepath.Join(dir, ".env")
	if err := os.WriteFile(original, []byte("X=1\n"), 0600); err != nil {
		t.Fatal(err)
	}

	bak, err := CreateBackup(original)
	if err != nil || bak == nil {
		t.Fatalf("setup failed: %v", err)
	}

	if err := bak.Discard(); err != nil {
		t.Fatalf("discard failed: %v", err)
	}

	if _, err := os.Stat(bak.BackupPath); !os.IsNotExist(err) {
		t.Error("backup file should not exist after discard")
	}
}

func TestBackup_NilRestore(t *testing.T) {
	var b *Backup
	if err := b.Restore(); err != nil {
		t.Errorf("nil restore should be no-op, got: %v", err)
	}
}

func TestBackup_NilDiscard(t *testing.T) {
	var b *Backup
	if err := b.Discard(); err != nil {
		t.Errorf("nil discard should be no-op, got: %v", err)
	}
}
