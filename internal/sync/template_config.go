package sync

import (
	"os"
	"strconv"
)

// TemplateConfigFromEnv loads TemplateConfig from environment variables.
//
//	VAULTPULL_TEMPLATE_ENABLED=true
//	VAULTPULL_TEMPLATE_PATH=/path/to/template.tmpl
//	VAULTPULL_TEMPLATE_OUTPUT=/path/to/output.conf
func TemplateConfigFromEnv() TemplateConfig {
	enabled, _ := strconv.ParseBool(os.Getenv("VAULTPULL_TEMPLATE_ENABLED"))
	return TemplateConfig{
		Enabled:      enabled,
		TemplatePath: os.Getenv("VAULTPULL_TEMPLATE_PATH"),
		OutputPath:   os.Getenv("VAULTPULL_TEMPLATE_OUTPUT"),
	}
}
