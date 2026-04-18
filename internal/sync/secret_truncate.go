package sync

import (
	"fmt"
	"os"
	"strconv"
)

// TruncateConfig controls whether secret values are truncated to a max length.
type TruncateConfig struct {
	Enabled   bool
	MaxLength int
	Suffix    string
}

// TruncateConfigFromEnv loads truncation config from environment variables.
func TruncateConfigFromEnv() TruncateConfig {
	cfg := TruncateConfig{
		MaxLength: 256,
		Suffix:    "...",
	}

	if v := os.Getenv("VAULTPULL_TRUNCATE_ENABLED"); v == "true" || v == "1" {
		cfg.Enabled = true
	}

	if v := os.Getenv("VAULTPULL_TRUNCATE_MAX_LENGTH"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			cfg.MaxLength = n
		}
	}

	if v := os.Getenv("VAULTPULL_TRUNCATE_SUFFIX"); v != "" {
		cfg.Suffix = v
	}

	return cfg
}

// TruncateSecrets returns a new map with values truncated to MaxLength.
// If disabled, the original map is returned unchanged.
func TruncateSecrets(secrets map[string]string, cfg TruncateConfig) (map[string]string, []string) {
	if !cfg.Enabled {
		return secrets, nil
	}

	result := make(map[string]string, len(secrets))
	var truncated []string

	for k, v := range secrets {
		if len(v) > cfg.MaxLength {
			result[k] = v[:cfg.MaxLength] + cfg.Suffix
			truncated = append(truncated, fmt.Sprintf("%s (was %d chars)", k, len(v)))
		} else {
			result[k] = v
		}
	}

	return result, truncated
}
