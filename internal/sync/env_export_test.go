package sync

import (
	"bytes"
	"strings"
	"testing"
)

func TestExportSecrets_DotenvFormat(t *testing.T) {
	var buf bytes.Buffer
	secrets := map[string]string{"DB_HOST": "localhost", "DB_PORT": "5432"}
	if err := ExportSecrets(&buf, secrets, FormatDotenv); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "DB_HOST=localhost") {
		t.Errorf("expected DB_HOST=localhost in output, got: %s", out)
	}
	if !strings.Contains(out, "DB_PORT=5432") {
		t.Errorf("expected DB_PORT=5432 in output, got: %s", out)
	}
}

func TestExportSecrets_ExportFormat(t *testing.T) {
	var buf bytes.Buffer
	secrets := map[string]string{"API_KEY": "abc123"}
	if err := ExportSecrets(&buf, secrets, FormatExport); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.HasPrefix(out, "export API_KEY=") {
		t.Errorf("expected export prefix, got: %s", out)
	}
}

func TestExportSecrets_JSONFormat(t *testing.T) {
	var buf bytes.Buffer
	secrets := map[string]string{"FOO": "bar"}
	if err := ExportSecrets(&buf, secrets, FormatJSON); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, `"FOO"`) {
		t.Errorf("expected JSON key FOO, got: %s", out)
	}
	if !strings.Contains(out, `"bar"`) {
		t.Errorf("expected JSON value bar, got: %s", out)
	}
}

func TestExportSecrets_QuotesValuesWithSpaces(t *testing.T) {
	var buf bytes.Buffer
	secrets := map[string]string{"MSG": "hello world"}
	if err := ExportSecrets(&buf, secrets, FormatDotenv); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, `MSG="hello world"`) {
		t.Errorf("expected quoted value, got: %s", out)
	}
}

func TestExportSecrets_SortedOutput(t *testing.T) {
	var buf bytes.Buffer
	secrets := map[string]string{"Z_KEY": "z", "A_KEY": "a", "M_KEY": "m"}
	ExportSecrets(&buf, secrets, FormatDotenv)
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 3 || !strings.HasPrefix(lines[0], "A_KEY") {
		t.Errorf("expected sorted output, got: %v", lines)
	}
}
