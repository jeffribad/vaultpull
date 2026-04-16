package sync

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTempTemplate(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "*.tmpl")
	if err != nil {
		t.Fatal(err)
	}
	_, _ = f.WriteString(content)
	_ = f.Close()
	return f.Name()
}

func TestRenderTemplate_BasicSubstitution(t *testing.T) {
	tmplPath := writeTempTemplate(t, `DB_HOST={{ index . "DB_HOST" }}`)
	secrets := map[string]string{"DB_HOST": "localhost"}
	out, err := RenderTemplate(TemplateConfig{TemplatePath: tmplPath}, secrets)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(out) != "DB_HOST=localhost" {
		t.Errorf("got %q", string(out))
	}
}

func TestRenderTemplate_MissingKey_ReturnsError(t *testing.T) {
	tmplPath := writeTempTemplate(t, `{{ index . "MISSING_KEY" }}`)
	_, err := RenderTemplate(TemplateConfig{TemplatePath: tmplPath}, map[string]string{})
	if err == nil {
		t.Fatal("expected error for missing key")
	}
}

func TestRenderTemplate_EmptyPath_ReturnsError(t *testing.T) {
	_, err := RenderTemplate(TemplateConfig{}, map[string]string{})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestWriteRenderedTemplate_WritesFile(t *testing.T) {
	tmplPath := writeTempTemplate(t, `HOST={{ index . "HOST" }}`)
	outPath := filepath.Join(t.TempDir(), "out.conf")
	cfg := TemplateConfig{TemplatePath: tmplPath, OutputPath: outPath}
	if err := WriteRenderedTemplate(cfg, map[string]string{"HOST": "example.com"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	data, _ := os.ReadFile(outPath)
	if string(data) != "HOST=example.com" {
		t.Errorf("got %q", string(data))
	}
}

func TestWriteRenderedTemplate_NoOutputPath_ReturnsError(t *testing.T) {
	tmplPath := writeTempTemplate(t, `hello`)
	cfg := TemplateConfig{TemplatePath: tmplPath}
	if err := WriteRenderedTemplate(cfg, map[string]string{}); err == nil {
		t.Fatal("expected error")
	}
}
