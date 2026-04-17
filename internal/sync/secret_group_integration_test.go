package sync

import (
	"os"
	"path/filepath"
	"testing"
)

func TestWriteGroups_Integration_CreatesFiles(t *testing.T) {
	dir := t.TempDir()

	groups := map[string]map[string]string{
		"backend":  {"DB_URL": "postgres://localhost"},
		"frontend": {"API_URL": "https://api.example.com"},
	}

	if err := WriteGroups(groups, dir, true); err != nil {
		t.Fatalf("WriteGroups failed: %v", err)
	}

	for group := range groups {
		path := filepath.Join(dir, group+".env")
		if _, err := os.Stat(path); os.IsNotExist(err) {
			t.Errorf("expected file %s to exist", path)
		}
	}
}

func TestWriteGroups_Integration_FullPipeline(t *testing.T) {
	dir := t.TempDir()

	secrets := map[string]string{
		"DB_HOST":   "localhost",
		"REDIS_URL": "redis://localhost",
		"STRIPE_SK": "sk_test_123",
	}
	labels := map[string]map[string]string{
		"DB_HOST":   {"group": "backend"},
		"REDIS_URL": {"group": "backend"},
		"STRIPE_SK": {"group": "payments"},
	}

	groups := GroupSecrets(secrets, labels, "group")
	if err := WriteGroups(groups, dir, true); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	backendPath := filepath.Join(dir, "backend.env")
	data, err := os.ReadFile(backendPath)
	if err != nil {
		t.Fatalf("could not read backend.env: %v", err)
	}
	content := string(data)
	if len(content) == 0 {
		t.Error("backend.env should not be empty")
	}
}
