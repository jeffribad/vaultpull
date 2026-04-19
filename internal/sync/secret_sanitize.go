package sync

import (
	"os"
	"strings"
)

// SanitizeConfig controls how secret values are sanitized before writing.
type SanitizeConfig struct {
	Enabled        bool
	StripControl   bool
	TrimWhitespace bool
	ReplaceNewlines bool
	NewlineReplacement string
}

// SanitizeConfigFromEnv loads sanitize config from environment variables.
func SanitizeConfigFromEnv() SanitizeConfig {
	enabled := isTruthy(os.Getenv("VAULTPULL_SANITIZE_ENABLED"))
	strip := os.Getenv("VAULTPULL_SANITIZE_STRIP_CONTROL")
	trim := os.Getenv("VAULTPULL_SANITIZE_TRIM_WHITESPACE")
	replaceNL := os.Getenv("VAULTPULL_SANITIZE_REPLACE_NEWLINES")
	nlReplacement := os.Getenv("VAULTPULL_SANITIZE_NEWLINE_REPLACEMENT")
	if nlReplacement == "" {
		nlReplacement = " "
	}
	return SanitizeConfig{
		Enabled:            enabled,
		StripControl:       enabled && (strip == "" || isTruthy(strip)),
		TrimWhitespace:     enabled && (trim == "" || isTruthy(trim)),
		ReplaceNewlines:    enabled && isTruthy(replaceNL),
		NewlineReplacement: nlReplacement,
	}
}

// SanitizeSecrets cleans secret values according to the provided config.
func SanitizeSecrets(secrets map[string]string, cfg SanitizeConfig) map[string]string {
	if !cfg.Enabled {
		return secrets
	}
	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		if cfg.TrimWhitespace {
			v = strings.TrimSpace(v)
		}
		if cfg.ReplaceNewlines {
			v = strings.ReplaceAll(v, "\r\n", cfg.NewlineReplacement)
			v = strings.ReplaceAll(v, "\n", cfg.NewlineReplacement)
			v = strings.ReplaceAll(v, "\r", cfg.NewlineReplacement)
		}
		if cfg.StripControl {
			v = stripControlChars(v)
		}
		out[k] = v
	}
	return out
}

func stripControlChars(s string) string {
	var b strings.Builder
	for _, r := range s {
		if r >= 32 || r == '\t' {
			b.WriteRune(r)
		}
	}
	return b.String()
}
