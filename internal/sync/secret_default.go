package sync

import (
	"os"
	"strings"
)

// DefaultConfigFromEnv loads default value configuration from environment variables.
type DefaultConfig struct {
	Enabled  bool
	Defaults map[string]string
}

// DefaultConfigFromEnv reads VAULTPULL_DEFAULTS_ENABLED and VAULTPULL_DEFAULTS (KEY=VALUE,...)
func DefaultConfigFromEnv() DefaultConfig {
	enabled := strings.TrimSpace(os.Getenv("VAULTPULL_DEFAULTS_ENABLED"))
	defaults := strings.TrimSpace(os.Getenv("VAULTPULL_DEFAULTS"))

	cfg := DefaultConfig{
		Enabled:  enabled == "true" || enabled == "1",
		Defaults: make(map[string]string),
	}

	if defaults == "" {
		return cfg
	}

	for _, pair := range strings.Split(defaults, ",") {
		parts := strings.SplitN(strings.TrimSpace(pair), "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		val := strings.TrimSpace(parts[1])
		if key != "" {
			cfg.Defaults[key] = val
		}
	}

	return cfg
}

// ApplyDefaults fills in missing keys in secrets with configured default values.
// Only keys absent from the secrets map are populated.
func ApplyDefaults(cfg DefaultConfig, secrets map[string]string) map[string]string {
	if !cfg.Enabled || len(cfg.Defaults) == 0 {
		return secrets
	}

	result := make(map[string]string, len(secrets))
	for k, v := range secrets {
		result[k] = v
	}

	for k, v := range cfg.Defaults {
		if _, exists := result[k]; !exists {
			result[k] = v
		}
	}

	return result
}
