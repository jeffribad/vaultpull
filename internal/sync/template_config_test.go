package sync

import (
	"os"
	"testing"
)

func TestTemplateConfigFromEnv_Defaults(t *testing.T) {
	os.Unsetenv("VAULTPULL_TEMPLATE_ENABLED")
	os.Unsetenv("VAULTPULL_TEMPLATE_PATH")
	os.Unsetenv("VAULTPULL_TEMPLATE_OUTPUT")

	cfg := TemplateConfigFromEnv()
	if cfg.Enabled {
		t.Error("expected Enabled=false by default")
	}
	if cfg.TemplatePath != "" {
		t.Errorf("expected empty TemplatePath, got %q", cfg.TemplatePath)
	}
	if cfg.OutputPath != "" {
		t.Errorf("expected empty OutputPath, got %q", cfg.OutputPath)
	}
}

func TestTemplateConfigFromEnv_Enabled(t *testing.T) {
	t.Setenv("VAULTPULL_TEMPLATE_ENABLED", "true")
	t.Setenv("VAULTPULL_TEMPLATE_PATH", "/etc/app/template.tmpl")
	t.Setenv("VAULTPULL_TEMPLATE_OUTPUT", "/etc/app/config.conf")

	cfg := TemplateConfigFromEnv()
	if !cfg.Enabled {
		t.Error("expected Enabled=true")
	}
	if cfg.TemplatePath != "/etc/app/template.tmpl" {
		t.Errorf("unexpected TemplatePath: %q", cfg.TemplatePath)
	}
	if cfg.OutputPath != "/etc/app/config.conf" {
		t.Errorf("unexpected OutputPath: %q", cfg.OutputPath)
	}
}

func TestTemplateConfigFromEnv_NumericEnabled(t *testing.T) {
	t.Setenv("VAULTPULL_TEMPLATE_ENABLED", "1")
	cfg := TemplateConfigFromEnv()
	if !cfg.Enabled {
		t.Error("expected Enabled=true for value '1'")
	}
}
