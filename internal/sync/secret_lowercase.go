package sync

import (
	"os"
	"strings"
)

// LowercaseConfig controls whether secret keys are lowercased.
type LowercaseConfig struct {
	Enabled bool
}

// LowercaseConfigFromEnv reads lowercase config from environment variables.
func LowercaseConfigFromEnv() LowercaseConfig {
	val := strings.TrimSpace(os.Getenv("VAULTPULL_LOWERCASE_KEYS"))
	return LowercaseConfig{
		Enabled: val == "true" || val == "1",
	}
}

// LowercaseKeys returns a new map with all keys converted to lowercase.
// If two keys collide after lowercasing, the last one wins (map iteration order).
// If disabled, the original map is returned unchanged.
func LowercaseKeys(cfg LowercaseConfig, secrets map[string]string) map[string]string {
	if !cfg.Enabled {
		return secrets
	}

	result := make(map[string]string, len(secrets))
	for k, v := range secrets {
		result[strings.ToLower(k)] = v
	}
	return result
}
