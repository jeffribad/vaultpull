package sync

import (
	"os"
	"strings"
)

// UnsetConfig controls whether secrets present in the environment but
// absent from Vault should be removed from the output map.
type UnsetConfig struct {
	Enabled     bool
	Keys        []string // explicit keys to unset regardless of Vault
	SyncWithEnv bool     // remove keys that exist in env but not in vault secrets
}

// UnsetConfigFromEnv reads UnsetConfig from environment variables.
//
//	VAULTPULL_UNSET_ENABLED=true
//	VAULTPULL_UNSET_KEYS=OLD_KEY,DEPRECATED_KEY
//	VAULTPULL_UNSET_SYNC_WITH_ENV=true
func UnsetConfigFromEnv() UnsetConfig {
	cfg := UnsetConfig{}

	if v := os.Getenv("VAULTPULL_UNSET_ENABLED"); isTruthy(v) {
		cfg.Enabled = true
	}

	if v := os.Getenv("VAULTPULL_UNSET_KEYS"); v != "" {
		cfg.Keys = splitTrimmed(v, ",")
	}

	if v := os.Getenv("VAULTPULL_UNSET_SYNC_WITH_ENV"); isTruthy(v) {
		cfg.SyncWithEnv = true
	}

	return cfg
}

// ApplyUnset removes keys from secrets according to the UnsetConfig.
// Explicit keys are always removed when enabled. When SyncWithEnv is true,
// any key present in the current OS environment but absent from secrets is
// also removed from the output (so callers can detect stale env vars).
func ApplyUnset(cfg UnsetConfig, secrets map[string]string) map[string]string {
	if !cfg.Enabled {
		return secrets
	}

	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		out[k] = v
	}

	// Remove explicitly listed keys.
	for _, key := range cfg.Keys {
		delete(out, key)
	}

	// Remove keys present in the OS environment but not in the vault secrets.
	if cfg.SyncWithEnv {
		for _, pair := range os.Environ() {
			parts := strings.SplitN(pair, "=", 2)
			if len(parts) < 1 {
				continue
			}
			envKey := parts[0]
			if _, inVault := secrets[envKey]; !inVault {
				delete(out, envKey)
			}
		}
	}

	return out
}
