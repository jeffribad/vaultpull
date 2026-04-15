package sync

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Backup holds the path to a backup file and the original file it was created from.
type Backup struct {
	OriginalPath string
	BackupPath   string
	CreatedAt    time.Time
}

// CreateBackup copies the contents of the given file to a timestamped backup file.
// If the original file does not exist, it returns nil without error.
func CreateBackup(path string) (*Backup, error) {
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("rollback: read original file: %w", err)
	}

	dir := filepath.Dir(path)
	base := filepath.Base(path)
	timestamp := time.Now().UTC().Format("20060102T150405Z")
	backupPath := filepath.Join(dir, fmt.Sprintf(".%s.%s.bak", base, timestamp))

	if err := os.WriteFile(backupPath, data, 0600); err != nil {
		return nil, fmt.Errorf("rollback: write backup file: %w", err)
	}

	return &Backup{
		OriginalPath: path,
		BackupPath:   backupPath,
		CreatedAt:    time.Now().UTC(),
	}, nil
}

// Restore writes the backup file contents back to the original path and removes the backup.
func (b *Backup) Restore() error {
	if b == nil {
		return nil
	}
	data, err := os.ReadFile(b.BackupPath)
	if err != nil {
		return fmt.Errorf("rollback: read backup file: %w", err)
	}
	if err := os.WriteFile(b.OriginalPath, data, 0600); err != nil {
		return fmt.Errorf("rollback: restore original file: %w", err)
	}
	return os.Remove(b.BackupPath)
}

// Discard removes the backup file without restoring it.
func (b *Backup) Discard() error {
	if b == nil {
		return nil
	}
	return os.Remove(b.BackupPath)
}
