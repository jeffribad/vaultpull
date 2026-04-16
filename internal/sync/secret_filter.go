package sync

import (
	"os"
	"strings"
)

// SecretFilterConfig holds configuration for key-based filtering.
type SecretFilterConfig struct {
	IncludeKeys []string
	ExcludeKeys []string
}

// SecretFilterConfigFromEnv reads filter config from environment variables.
// VAULTPULL_INCLUDE_KEYS and VAULTPULL_EXCLUDE_KEYS are comma-separated lists.
func SecretFilterConfigFromEnv() SecretFilterConfig {
	return SecretFilterConfig{
		IncludeKeys: splitTrimmed(os.Getenv("VAULTPULL_INCLUDE_KEYS")),
		ExcludeKeys: splitTrimmed(os.Getenv("VAULTPULL_EXCLUDE_KEYS")),
	}
}

func splitTrimmed(s string) []string {
	if s == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if t := strings.TrimSpace(p); t != "" {
			out = append(out, t)
		}
	}
	return out
}

// ApplySecretFilter returns a filtered map of secrets based on include/exclude rules.
// If IncludeKeys is non-empty, only those keys are kept.
// ExcludeKeys are always removed after include filtering.
func ApplySecretFilter(secrets map[string]string, cfg SecretFilterConfig) map[string]string {
	result := make(map[string]string, len(secrets))

	for k, v := range secrets {
		if len(cfg.IncludeKeys) > 0 && !containsKey(cfg.IncludeKeys, k) {
			continue
		}
		if containsKey(cfg.ExcludeKeys, k) {
			continue
		}
		result[k] = v
	}
	return result
}

func containsKey(list []string, key string) bool {
	upper := strings.ToUpper(key)
	for _, item := range list {
		if strings.ToUpper(item) == upper {
			return true
		}
	}
	return false
}
