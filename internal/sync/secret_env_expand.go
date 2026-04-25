package sync

import (
	"os"
	"strings"
)

// ExpandConfig controls shell-style variable expansion within secret values.
type ExpandConfig struct {
	Enabled  bool
	AllowEnv bool // also expand from OS environment
}

// ExpandConfigFromEnv reads expansion settings from environment variables.
func ExpandConfigFromEnv() ExpandConfig {
	return ExpandConfig{
		Enabled:  isTruthy(os.Getenv("VAULTPULL_EXPAND_ENABLED")),
		AllowEnv: isTruthy(os.Getenv("VAULTPULL_EXPAND_ALLOW_ENV")),
	}
}

// ExpandSecrets performs shell-style ${VAR} and $VAR expansion within secret
// values, resolving references against other keys in the map and optionally
// the OS environment.
func ExpandSecrets(secrets map[string]string, cfg ExpandConfig) (map[string]string, error) {
	if !cfg.Enabled {
		return secrets, nil
	}

	result := make(map[string]string, len(secrets))

	mapping := func(key string) string {
		// Check secrets map first.
		if val, ok := secrets[key]; ok {
			return val
		}
		// Optionally fall back to OS environment.
		if cfg.AllowEnv {
			return os.Getenv(key)
		}
		return ""
	}

	for k, v := range secrets {
		result[k] = os.Expand(v, mapping)
	}

	return result, nil
}

// expandValue expands a single value using the provided lookup map.
func expandValue(value string, lookup map[string]string) string {
	return os.Expand(value, func(key string) string {
		if val, ok := lookup[key]; ok {
			return val
		}
		return strings.ToUpper(key)
	})
}
