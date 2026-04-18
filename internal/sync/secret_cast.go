package sync

import (
	"os"
	"strconv"
	"strings"
)

// CastConfig controls type coercion of secret values.
type CastConfig struct {
	Enabled   bool
	BoolKeys  []string
	IntKeys   []string
}

// CastConfigFromEnv loads cast config from environment variables.
func CastConfigFromEnv() CastConfig {
	enabled := strings.TrimSpace(os.Getenv("VAULTPULL_CAST_ENABLED"))
	boolKeys := splitTrimmed(os.Getenv("VAULTPULL_CAST_BOOL_KEYS"), ",")
	intKeys := splitTrimmed(os.Getenv("VAULTPULL_CAST_INT_KEYS"), ",")
	return CastConfig{
		Enabled:  enabled == "true" || enabled == "1",
		BoolKeys: boolKeys,
		IntKeys:  intKeys,
	}
}

// CastSecrets normalises secret values according to the cast config.
// Bool keys: "true"/"1"/"yes" -> "true", anything else -> "false".
// Int keys: strips non-numeric characters, falls back to "0".
func CastSecrets(secrets map[string]string, cfg CastConfig) map[string]string {
	if !cfg.Enabled {
		return secrets
	}
	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		out[k] = v
	}
	for _, key := range cfg.BoolKeys {
		if v, ok := out[key]; ok {
			lower := strings.ToLower(strings.TrimSpace(v))
			if lower == "true" || lower == "1" || lower == "yes" {
				out[key] = "true"
			} else {
				out[key] = "false"
			}
		}
	}
	for _, key := range cfg.IntKeys {
		if v, ok := out[key]; ok {
			var digits strings.Builder
			for _, ch := range v {
				if ch >= '0' && ch <= '9' {
					digits.WriteRune(ch)
				}
			}
			s := digits.String()
			if s == "" {
				s = "0"
			}
			// normalise: strip leading zeros unless "0"
			if n, err := strconv.Atoi(s); err == nil {
				s = strconv.Itoa(n)
			}
			out[key] = s
		}
	}
	return out
}
