package sync

import (
	"bytes"
	"fmt"
	"os"
	"text/template"
)

// TemplateConfig holds configuration for template rendering.
type TemplateConfig struct {
	Enabled      bool
	TemplatePath string
	OutputPath   string
}

// RenderTemplate renders a Go template file using the provided secrets as data.
func RenderTemplate(cfg TemplateConfig, secrets map[string]string) ([]byte, error) {
	if cfg.TemplatePath == "" {
		return nil, fmt.Errorf("template path is required")
	}

	tmplBytes, err := os.ReadFile(cfg.TemplatePath)
	if err != nil {
		return nil, fmt.Errorf("reading template: %w", err)
	}

	tmpl, err := template.New("secrets").Option("missingkey=error").Parse(string(tmplBytes))
	if err != nil {
		return nil, fmt.Errorf("parsing template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, secrets); err != nil {
		return nil, fmt.Errorf("executing template: %w", err)
	}

	return buf.Bytes(), nil
}

// WriteRenderedTemplate renders and writes the output to cfg.OutputPath.
func WriteRenderedTemplate(cfg TemplateConfig, secrets map[string]string) error {
	if cfg.OutputPath == "" {
		return fmt.Errorf("output path is required")
	}

	data, err := RenderTemplate(cfg, secrets)
	if err != nil {
		return err
	}

	return os.WriteFile(cfg.OutputPath, data, 0600)
}
