package sync

import (
	"os"
	"strings"
)

// FormatConfig controls how secret values are formatted before writing.
type FormatConfig struct {
	Enabled    bool
	TrimSpace  bool
	NormalizeNewlines bool
	StripNulls bool
}

// FormatConfigFromEnv loads format config from environment variables.
func FormatConfigFromEnv() FormatConfig {
	enabled := os.Getenv("VAULTPULL_FORMAT_ENABLED")
	trim := os.Getenv("VAULTPULL_FORMAT_TRIM_SPACE")
	newlines := os.Getenv("VAULTPULL_FORMAT_NORMALIZE_NEWLINES")
	nulls := os.Getenv("VAULTPULL_FORMAT_STRIP_NULLS")

	return FormatConfig{
		Enabled:           isTruthy(enabled),
		TrimSpace:         isTruthy(trim),
		NormalizeNewlines: isTruthy(newlines),
		StripNulls:        isTruthy(nulls),
	}
}

func isTruthy(s string) bool {
	s = strings.ToLower(strings.TrimSpace(s))
	return s == "true" || s == "1" || s == "yes"
}

// ApplyFormat applies formatting rules to all secrets in the map.
func ApplyFormat(secrets map[string]string, cfg FormatConfig) map[string]string {
	if !cfg.Enabled {
		return secrets
	}

	result := make(map[string]string, len(secrets))
	for k, v := range secrets {
		if cfg.TrimSpace {
			v = strings.TrimSpace(v)
		}
		if cfg.NormalizeNewlines {
			v = strings.ReplaceAll(v, "\r\n", "\n")
			v = strings.ReplaceAll(v, "\r", "\n")
		}
		if cfg.StripNulls {
			v = strings.ReplaceAll(v, "\x00", "")
		}
		result[k] = v
	}
	return result
}
