package sync

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestTemplate_Integration_RenderAndWrite(t *testing.T) {
	dir := t.TempDir()
	tmplPath := filepath.Join(dir, "app.tmpl")
	outPath := filepath.Join(dir, "app.conf")

	tmplContent := strings.Join([]string{
		`database_url={{ index . "DATABASE_URL" }}`,
		`api_key={{ index . "API_KEY" }}`,
	}, "\n")
	if err := os.WriteFile(tmplPath, []byte(tmplContent), 0644); err != nil {
		t.Fatal(err)
	}

	secrets := map[string]string{
		"DATABASE_URL": "postgres://localhost/mydb",
		"API_KEY":      "supersecret",
	}

	cfg := TemplateConfig{TemplatePath: tmplPath, OutputPath: outPath}
	if err := WriteRenderedTemplate(cfg, secrets); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatal(err)
	}

	result := string(data)
	if !strings.Contains(result, "postgres://localhost/mydb") {
		t.Errorf("missing DATABASE_URL in output: %q", result)
	}
	if !strings.Contains(result, "supersecret") {
		t.Errorf("missing API_KEY in output: %q", result)
	}

	info, _ := os.Stat(outPath)
	if info.Mode().Perm() != 0600 {
		t.Errorf("expected file perm 0600, got %v", info.Mode().Perm())
	}
}
