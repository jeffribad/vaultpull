package sync

import (
	"os"
	"strings"
)

// PrefixFilterConfig controls filtering secrets by key prefix.
type PrefixFilterConfig struct {
	Enabled        bool
	AllowPrefixes  []string
	DenyPrefixes   []string
}

// PrefixFilterConfigFromEnv loads prefix filter config from environment variables.
func PrefixFilterConfigFromEnv() PrefixFilterConfig {
	return PrefixFilterConfig{
		Enabled:       os.Getenv("VAULTPULL_PREFIX_FILTER_ENABLED") == "1" || os.Getenv("VAULTPULL_PREFIX_FILTER_ENABLED") == "true",
		AllowPrefixes: splitTrimmed(os.Getenv("VAULTPULL_PREFIX_ALLOW"), ","),
		DenyPrefixes:  splitTrimmed(os.Getenv("VAULTPULL_PREFIX_DENY"), ","),
	}
}

// ApplyPrefixFilter returns a filtered copy of secrets based on allow/deny prefix rules.
// If allow prefixes are set, only matching keys are kept.
// Deny prefixes are applied after allow filtering.
func ApplyPrefixFilter(secrets map[string]string, cfg PrefixFilterConfig) map[string]string {
	if !cfg.Enabled {
		return secrets
	}

	result := make(map[string]string, len(secrets))

	for k, v := range secrets {
		if len(cfg.AllowPrefixes) > 0 && !hasAnyPrefix(k, cfg.AllowPrefixes) {
			continue
		}
		if hasAnyPrefix(k, cfg.DenyPrefixes) {
			continue
		}
		result[k] = v
	}

	return result
}

func hasAnyPrefix(key string, prefixes []string) bool {
	upper := strings.ToUpper(key)
	for _, p := range prefixes {
		if strings.HasPrefix(upper, strings.ToUpper(p)) {
			return true
		}
	}
	return false
}
