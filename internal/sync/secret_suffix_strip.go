package sync

import (
	"os"
	"strings"
)

// SuffixStripConfig controls stripping of key suffixes.
type SuffixStripConfig struct {
	Enabled bool
	Suffix  string
}

// SuffixStripConfigFromEnv loads SuffixStripConfig from environment variables.
//
//	VAULTPULL_SUFFIX_STRIP_ENABLED=true
//	VAULTPULL_SUFFIX_STRIP_SUFFIX=_SECRET
func SuffixStripConfigFromEnv() SuffixStripConfig {
	enabled := isTruthy(os.Getenv("VAULTPULL_SUFFIX_STRIP_ENABLED"))
	suffix := os.Getenv("VAULTPULL_SUFFIX_STRIP_SUFFIX")
	return SuffixStripConfig{
		Enabled: enabled,
		Suffix:  suffix,
	}
}

// StripKeySuffix returns a new map with the configured suffix removed from
// matching keys. Keys that do not end with the suffix are left unchanged.
// If the config is disabled or the suffix is empty, the original map is
// returned as-is.
func StripKeySuffix(cfg SuffixStripConfig, secrets map[string]string) map[string]string {
	if !cfg.Enabled || cfg.Suffix == "" {
		return secrets
	}

	result := make(map[string]string, len(secrets))
	for k, v := range secrets {
		if strings.HasSuffix(k, cfg.Suffix) {
			newKey := strings.TrimSuffix(k, cfg.Suffix)
			if newKey == "" {
				// Avoid creating an empty key; keep the original.
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
