package sync

import (
	"os"
	"strings"
)

// MergeConfig controls how secrets from multiple sources are merged.
type MergeConfig struct {
	Enabled  bool
	Strategy string // "vault-wins", "local-wins", "union"
}

// MergeConfigFromEnv loads merge configuration from environment variables.
func MergeConfigFromEnv() MergeConfig {
	strategy := strings.TrimSpace(os.Getenv("VAULTPULL_MERGE_STRATEGY"))
	if strategy == "" {
		strategy = "vault-wins"
	}
	enabled := isTruthy(os.Getenv("VAULTPULL_MERGE_ENABLED"))
	return MergeConfig{
		Enabled:  enabled,
		Strategy: strategy,
	}
}

// MergeSecrets merges vault secrets with a base map according to the configured strategy.
// "vault-wins"  — vault values overwrite base values (default)
// "local-wins"  — base values are preserved when keys conflict
// "union"       — all keys from both maps are included; vault fills missing keys only
func MergeSecrets(cfg MergeConfig, vault map[string]string, base map[string]string) map[string]string {
	result := make(map[string]string, len(vault)+len(base))

	if !cfg.Enabled {
		for k, v := range vault {
			result[k] = v
		}
		return result
	}

	switch cfg.Strategy {
	case "local-wins":
		for k, v := range vault {
			result[k] = v
		}
		for k, v := range base {
			result[k] = v // base overwrites vault
		}
	case "union":
		for k, v := range base {
			result[k] = v
		}
		for k, v := range vault {
			if _, exists := result[k]; !exists {
				result[k] = v
			}
		}
	default: // "vault-wins"
		for k, v := range base {
			result[k] = v
		}
		for k, v := range vault {
			result[k] = v // vault overwrites base
		}
	}

	return result
}
