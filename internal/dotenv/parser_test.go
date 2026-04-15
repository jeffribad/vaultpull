package dotenv

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTempEnvFile(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatalf("failed to write temp env file: %v", err)
	}
	return path
}

func TestParse_BasicKeyValues(t *testing.T) {
	path := writeTempEnvFile(t, "DB_HOST=localhost\nDB_PORT=5432\n")
	got, err := Parse(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["DB_HOST"] != "localhost" {
		t.Errorf("expected DB_HOST=localhost, got %q", got["DB_HOST"])
	}
	if got["DB_PORT"] != "5432" {
		t.Errorf("expected DB_PORT=5432, got %q", got["DB_PORT"])
	}
}

func TestParse_SkipsComments(t *testing.T) {
	path := writeTempEnvFile(t, "# this is a comment\nAPP_ENV=production\n")
	got, err := Parse(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := got["# this is a comment"]; ok {
		t.Error("comment line should not be parsed as a key")
	}
	if got["APP_ENV"] != "production" {
		t.Errorf("expected APP_ENV=production, got %q", got["APP_ENV"])
	}
}

func TestParse_StripQuotes(t *testing.T) {
	path := writeTempEnvFile(t, `SECRET_KEY="my secret value"` + "\n" + `TOKEN='abc123'` + "\n")
	got, err := Parse(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["SECRET_KEY"] != "my secret value" {
		t.Errorf("expected unquoted value, got %q", got["SECRET_KEY"])
	}
	if got["TOKEN"] != "abc123" {
		t.Errorf("expected unquoted value, got %q", got["TOKEN"])
	}
}

func TestParse_NonExistentFile(t *testing.T) {
	got, err := Parse("/tmp/does_not_exist_vaultpull.env")
	if err != nil {
		t.Fatalf("expected no error for missing file, got: %v", err)
	}
	if len(got) != 0 {
		t.Errorf("expected empty map for missing file, got %v", got)
	}
}

func TestParse_SkipsInvalidLines(t *testing.T) {
	path := writeTempEnvFile(t, "INVALID_LINE_NO_EQUALS\nVALID=yes\n")
	got, err := Parse(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := got["INVALID_LINE_NO_EQUALS"]; ok {
		t.Error("line without '=' should be skipped")
	}
	if got["VALID"] != "yes" {
		t.Errorf("expected VALID=yes, got %q", got["VALID"])
	}
}
