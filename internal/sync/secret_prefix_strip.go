package sync

import (
	"os"
	"strings"
)

// PrefixStripConfig controls automatic stripping of a common prefix from secret keys.
type PrefixStripConfig struct {
	Enabled bool
	Prefix  string
}

// PrefixStripConfigFromEnv loads PrefixStripConfig from environment variables.
//
//	VAULTPULL_PREFIX_STRIP_ENABLED=true
//	VAULTPULL_PREFIX_STRIP_PREFIX=APP_
func PrefixStripConfigFromEnv() PrefixStripConfig {
	enabled := os.Getenv("VAULTPULL_PREFIX_STRIP_ENABLED")
	prefix := os.Getenv("VAULTPULL_PREFIX_STRIP_PREFIX")
	return PrefixStripConfig{
		Enabled: isTruthy(enabled),
		Prefix:  prefix,
	}
}

// StripKeyPrefix removes a configured prefix from all matching secret keys.
// Keys that do not have the prefix are kept as-is.
// If the config is disabled or the prefix is empty, the original map is returned unchanged.
func StripKeyPrefix(cfg PrefixStripConfig, secrets map[string]string) map[string]string {
	if !cfg.Enabled || cfg.Prefix == "" {
		return secrets
	}

	result := make(map[string]string, len(secrets))
	for k, v := range secrets {
		if strings.HasPrefix(k, cfg.Prefix) {
			newKey := strings.TrimPrefix(k, cfg.Prefix)
			if newKey == "" {
				// Avoid empty key — keep original
				result[k] = v
			} else {
				result[newKey] = v
			}
		} else {
			result[k] = v
		}
	}
	return result
}
