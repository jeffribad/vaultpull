package sync_test

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/your-org/vaultpull/internal/audit"
	"github.com/your-org/vaultpull/internal/sync"
	"github.com/your-org/vaultpull/internal/vault"
)

func newTestLogger(buf *bytes.Buffer) *audit.Logger {
	return audit.NewLogger(buf)
}

func newTestClient(t *testing.T, secrets map[string]string) *vault.Client {
	t.Helper()
	client, err := vault.NewClient("http://127.0.0.1:8200", "test-token")
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	client.SetFakeSecrets(secrets)
	return client
}

func TestSyncer_DryRun_DoesNotWriteFile(t *testing.T) {
	var buf bytes.Buffer
	logger := newTestLogger(&buf)
	client := newTestClient(t, map[string]string{"DB_HOST": "localhost", "DB_PORT": "5432"})
	policies := vault.Policies{"backend": {"DB_HOST", "DB_PORT"}}

	s := sync.New(client, logger)
	outPath := filepath.Join(t.TempDir(), ".env")

	result, err := s.Run("secret/app", policies, sync.Options{
		Role: "backend", OutputPath: outPath, DryRun: true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Written != 0 {
		t.Errorf("expected 0 written, got %d", result.Written)
	}
	if _, err := os.Stat(outPath); !os.IsNotExist(err) {
		t.Error("expected no file to be created in dry-run mode")
	}
}

func TestSyncer_WritesFilteredSecrets(t *testing.T) {
	var buf bytes.Buffer
	logger := newTestLogger(&buf)
	client := newTestClient(t, map[string]string{"DB_HOST": "localhost", "SECRET_KEY": "abc"})
	policies := vault.Policies{"backend": {"DB_HOST"}}

	s := sync.New(client, logger)
	outPath := filepath.Join(t.TempDir(), ".env")

	result, err := s.Run("secret/app", policies, sync.Options{
		Role: "backend", OutputPath: outPath,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Written != 1 {
		t.Errorf("expected 1 written, got %d", result.Written)
	}

	var entry map[string]interface{}
	if err := json.NewDecoder(&buf).Decode(&entry); err != nil {
		t.Fatalf("audit log not valid JSON: %v", err)
	}
}

func TestSyncer_UnknownRole_ReturnsError(t *testing.T) {
	var buf bytes.Buffer
	client := newTestClient(t, map[string]string{"DB_HOST": "localhost"})
	policies := vault.Policies{"backend": {"DB_HOST"}}

	s := sync.New(client, audit.NewLogger(&buf))
	_, err := s.Run("secret/app", policies, sync.Options{Role: "unknown"})
	if err == nil {
		t.Fatal("expected error for unknown role")
	}
}
