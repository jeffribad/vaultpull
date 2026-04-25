package sync

import (
	"os"
	"strings"
)

// PrefixAddConfigFromEnv loads prefix-add configuration from environment variables.
// VAULTPULL_PREFIX_ADD_ENABLED=true enables the feature.
// VAULTPULL_PREFIX_ADD_VALUE sets the prefix string to prepend to all keys.
// VAULTPULL_PREFIX_ADD_KEYS limits prefixing to specific comma-separated keys.
type PrefixAddConfig struct {
	Enabled bool
	Prefix  string
	Keys    []string
}

func PrefixAddConfigFromEnv() PrefixAddConfig {
	cfg := PrefixAddConfig{}
	cfg.Enabled = isTruthy(os.Getenv("VAULTPULL_PREFIX_ADD_ENABLED"))
	cfg.Prefix = os.Getenv("VAULTPULL_PREFIX_ADD_VALUE")
	raw := os.Getenv("VAULTPULL_PREFIX_ADD_KEYS")
	if raw != "" {
		cfg.Keys = splitTrimmed(raw, ",")
	}
	return cfg
}

// AddKeyPrefix prepends a prefix to all keys (or a subset) in the secrets map.
// Returns a new map; the original is not mutated.
func AddKeyPrefix(cfg PrefixAddConfig, secrets map[string]string) map[string]string {
	if !cfg.Enabled || cfg.Prefix == "" {
		return secrets
	}

	allowSet := make(map[string]bool, len(cfg.Keys))
	for _, k := range cfg.Keys {
		allowSet[strings.ToLower(k)] = true
	}

	result := make(map[string]string, len(secrets))
	for k, v := range secrets {
		if len(allowSet) == 0 || allowSet[strings.ToLower(k)] {
			result[cfg.Prefix+k] = v
		} else {
			result[k] = v
		}
	}
	return result
}
