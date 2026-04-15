package dotenv_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/your-org/vaultpull/internal/dotenv"
)

func TestWriter_WriteNewFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")

	w := dotenv.NewWriter(path, false)
	err := w.Write(map[string]string{
		"DB_HOST": "localhost",
		"DB_PORT": "5432",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	result, err := dotenv.Parse(path)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	if result["DB_HOST"] != "localhost" {
		t.Errorf("expected DB_HOST=localhost, got %q", result["DB_HOST"])
	}
	if result["DB_PORT"] != "5432" {
		t.Errorf("expected DB_PORT=5432, got %q", result["DB_PORT"])
	}
}

func TestWriter_MergesExistingKeys(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")

	_ = os.WriteFile(path, []byte("EXISTING_KEY=old_value\n"), 0600)

	w := dotenv.NewWriter(path, false)
	err := w.Write(map[string]string{"NEW_KEY": "new_value"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	result, err := dotenv.Parse(path)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	if result["EXISTING_KEY"] != "old_value" {
		t.Errorf("expected EXISTING_KEY to be preserved, got %q", result["EXISTING_KEY"])
	}
	if result["NEW_KEY"] != "new_value" {
		t.Errorf("expected NEW_KEY=new_value, got %q", result["NEW_KEY"])
	}
}

func TestWriter_OverwriteMode(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")

	_ = os.WriteFile(path, []byte("OLD_KEY=should_be_gone\n"), 0600)

	w := dotenv.NewWriter(path, true)
	err := w.Write(map[string]string{"ONLY_KEY": "value"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	result, err := dotenv.Parse(path)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	if _, ok := result["OLD_KEY"]; ok {
		t.Error("expected OLD_KEY to be removed in overwrite mode")
	}
}

func TestWriter_QuotesValuesWithSpaces(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")

	w := dotenv.NewWriter(path, true)
	err := w.Write(map[string]string{"MSG": "hello world"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	result, err := dotenv.Parse(path)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	if result["MSG"] != "hello world" {
		t.Errorf("expected MSG='hello world', got %q", result["MSG"])
	}
}
