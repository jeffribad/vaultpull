package sync

import (
	"os"
	"strings"
)

// NormalizeConfig controls key normalization behavior.
type NormalizeConfig struct {
	Enabled     bool
	UpperKeys   bool
	SnakeCase   bool
	StripDashes bool
}

// NormalizeConfigFromEnv loads normalization config from environment variables.
func NormalizeConfigFromEnv() NormalizeConfig {
	return NormalizeConfig{
		Enabled:     isTruthy(os.Getenv("VAULTPULL_NORMALIZE_ENABLED")),
		UpperKeys:   isTruthy(os.Getenv("VAULTPULL_NORMALIZE_UPPER_KEYS")),
		SnakeCase:   isTruthy(os.Getenv("VAULTPULL_NORMALIZE_SNAKE_CASE")),
		StripDashes: isTruthy(os.Getenv("VAULTPULL_NORMALIZE_STRIP_DASHES")),
	}
}

// NormalizeSecrets applies key normalization rules to the secrets map.
// Returns a new map with normalized keys; values are unchanged.
func NormalizeSecrets(cfg NormalizeConfig, secrets map[string]string) (map[string]string, error) {
	if !cfg.Enabled {
		return secrets, nil
	}

	result := make(map[string]string, len(secrets))
	for k, v := range secrets {
		nk := normalizeKey(cfg, k)
		result[nk] = v
	}
	return result, nil
}

func normalizeKey(cfg NormalizeConfig, key string) string {
	if cfg.SnakeCase {
		key = toSnakeCase(key)
	}
	if cfg.StripDashes {
		key = strings.ReplaceAll(key, "-", "_")
	}
	if cfg.UpperKeys {
		key = strings.ToUpper(key)
	}
	return key
}

// toSnakeCase converts camelCase or PascalCase to snake_case.
func toSnakeCase(s string) string {
	var b strings.Builder
	for i, r := range s {
		if r >= 'A' && r <= 'Z' {
			if i > 0 {
				b.WriteRune('_')
			}
			b.WriteRune(r + 32)
		} else {
			b.WriteRune(r)
		}
	}
	return b.String()
}
