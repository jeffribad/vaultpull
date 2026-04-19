package sync

import (
	"os"
	"strings"
)

// UppercaseConfig controls whether secret keys are uppercased before writing.
type UppercaseConfig struct {
	Enabled bool
}

// UppercaseConfigFromEnv loads UppercaseConfig from environment variables.
func UppercaseConfigFromEnv() UppercaseConfig {
	val := strings.TrimSpace(os.Getenv("VAULTPULL_UPPERCASE_KEYS"))
	return UppercaseConfig{
		Enabled: isTruthy(val),
	}
}

// UppercaseKeys returns a new map with all keys converted to uppercase.
// If disabled, the original map is returned unchanged.
func UppercaseKeys(cfg UppercaseConfig, secrets map[string]string) map[string]string {
	if !cfg.Enabled {
		return secrets
	}
	result := make(map[string]string, len(secrets))
	for k, v := range secrets {
		upper := strings.ToUpper(k)
		// last-write wins if collision
		result[upper] = v
	}
	return result
}
