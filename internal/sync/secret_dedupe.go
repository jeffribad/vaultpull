package sync

import "os"

// DedupeConfig controls deduplication of secret keys.
type DedupeConfig struct {
	Enabled       bool
	CaseSensitive bool
}

// DedupeConfigFromEnv loads deduplication config from environment variables.
func DedupeConfigFromEnv() DedupeConfig {
	enabled := os.Getenv("VAULTPULL_DEDUPE_ENABLED")
	caseSensitive := os.Getenv("VAULTPULL_DEDUPE_CASE_SENSITIVE")
	return DedupeConfig{
		Enabled:       enabled == "true" || enabled == "1",
		CaseSensitive: caseSensitive == "true" || caseSensitive == "1",
	}
}

// DedupeSecrets removes duplicate keys from secrets, keeping the last occurrence.
// If CaseSensitive is false, keys are compared case-insensitively and the
// canonical form of the last seen key is preserved.
func DedupeSecrets(secrets map[string]string, cfg DedupeConfig) map[string]string {
	if !cfg.Enabled {
		return secrets
	}
	if cfg.CaseSensitive {
		// map already has unique keys; return a copy
		out := make(map[string]string, len(secrets))
		for k, v := range secrets {
			out[k] = v
		}
		return out
	}
	// Case-insensitive: last write wins per lower-case key.
	// Preserve the original casing of the winning key.
	lower := make(map[string]string)   // lower -> canonical key
	out := make(map[string]string)
	for k, v := range secrets {
		lk := toLower(k)
		if prev, ok := lower[lk]; ok {
			delete(out, prev)
		}
		lower[lk] = k
		out[k] = v
	}
	return out
}

func toLower(s string) string {
	b := make([]byte, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= 'A' && c <= 'Z' {
			c += 'a' - 'A'
		}
		b[i] = c
	}
	return string(b)
}
