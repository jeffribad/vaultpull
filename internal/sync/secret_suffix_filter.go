package sync

import (
	"os"
	"strings"
)

// SuffixFilterConfig controls filtering secrets by key suffix.
type SuffixFilterConfig struct {
	Enabled     bool
	AllowSuffix []string
	DenySuffix  []string
}

// SuffixFilterConfigFromEnv loads SuffixFilterConfig from environment variables.
//
//	VAULTPULL_SUFFIX_FILTER_ENABLED=true
//	VAULTPULL_SUFFIX_FILTER_ALLOW=_URL,_HOST
//	VAULTPULL_SUFFIX_FILTER_DENY=_SECRET,_KEY
func SuffixFilterConfigFromEnv() SuffixFilterConfig {
	return SuffixFilterConfig{
		Enabled:     os.Getenv("VAULTPULL_SUFFIX_FILTER_ENABLED") == "true" || os.Getenv("VAULTPULL_SUFFIX_FILTER_ENABLED") == "1",
		AllowSuffix: splitTrimmed(os.Getenv("VAULTPULL_SUFFIX_FILTER_ALLOW"), ","),
		DenySuffix:  splitTrimmed(os.Getenv("VAULTPULL_SUFFIX_FILTER_DENY"), ","),
	}
}

// ApplySuffixFilter filters secrets by key suffix.
// If AllowSuffix is set, only keys ending with one of those suffixes are kept.
// If DenySuffix is set, keys ending with those suffixes are removed.
// Allow takes precedence over deny when both are specified.
func ApplySuffixFilter(secrets map[string]string, cfg SuffixFilterConfig) map[string]string {
	if !cfg.Enabled {
		return secrets
	}

	result := make(map[string]string, len(secrets))
	for k, v := range secrets {
		upper := strings.ToUpper(k)

		if len(cfg.AllowSuffix) > 0 {
			if hasAnySuffix(upper, cfg.AllowSuffix) {
				result[k] = v
			}
			continue
		}

		if len(cfg.DenySuffix) > 0 && hasAnySuffix(upper, cfg.DenySuffix) {
			continue
		}

		result[k] = v
	}
	return result
}

func hasAnySuffix(key string, suffixes []string) bool {
	for _, s := range suffixes {
		if strings.HasSuffix(key, strings.ToUpper(s)) {
			return true
		}
	}
	return false
}
