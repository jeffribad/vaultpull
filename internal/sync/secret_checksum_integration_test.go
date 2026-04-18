package sync

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestChecksum_Integration_ConfigDriven(t *testing.T) {
	dir := t.TempDir()
	outPath := filepath.Join(dir, ".env.sha256")

	t.Setenv("VAULTPULL_CHECKSUM_ENABLED", "1")
	t.Setenv("VAULTPULL_CHECKSUM_PATH", outPath)

	cfg := ChecksumConfigFromEnv()
	secrets := map[string]string{
		"API_KEY": "abc123",
		"DB_URL":  "postgres://localhost/db",
	}

	sum, err := WriteChecksum(cfg, secrets)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(sum) != 64 {
		t.Errorf("expected 64-char sha256 hex, got %d chars", len(sum))
	}

	data, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatalf("could not read checksum file: %v", err)
	}
	if strings.TrimSpace(string(data)) != sum {
		t.Error("file content does not match returned checksum")
	}
}

func TestChecksum_Integration_DisabledSkipsWrite(t *testing.T) {
	dir := t.TempDir()
	outPath := filepath.Join(dir, ".env.sha256")

	t.Setenv("VAULTPULL_CHECKSUM_ENABLED", "false")
	t.Setenv("VAULTPULL_CHECKSUM_PATH", outPath)

	cfg := ChecksumConfigFromEnv()
	_, err := WriteChecksum(cfg, map[string]string{"X": "y"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, statErr := os.Stat(outPath); !os.IsNotExist(statErr) {
		t.Error("expected no file written when disabled")
	}
}
