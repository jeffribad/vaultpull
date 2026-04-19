package sync

import (
	"os"
	"strings"
)

// EnvMapConfig controls key remapping from environment variable names.
type EnvMapConfig struct {
	Enabled bool
	Mappings map[string]string // vault key -> env key
}

// EnvMapConfigFromEnv loads EnvMapConfig from environment variables.
// VAULTPULL_ENVMAP_ENABLED=true
// VAULTPULL_ENVMAP_KEYS=VAULT_KEY:ENV_KEY,OTHER_KEY:OTHER_ENV
func EnvMapConfigFromEnv() EnvMapConfig {
	cfg := EnvMapConfig{
		Mappings: make(map[string]string),
	}

	val := os.Getenv("VAULTPULL_ENVMAP_ENABLED")
	cfg.Enabled = val == "true" || val == "1"

	raw := os.Getenv("VAULTPULL_ENVMAP_KEYS")
	if raw != "" {
		cfg.Mappings = parseEnvMapPairs(raw)
	}

	return cfg
}

func parseEnvMapPairs(raw string) map[string]string {
	result := make(map[string]string)
	for _, pair := range strings.Split(raw, ",") {
		parts := strings.SplitN(strings.TrimSpace(pair), ":", 2)
		if len(parts) != 2 {
			continue
		}
		k := strings.TrimSpace(parts[0])
		v := strings.TrimSpace(parts[1])
		if k != "" && v != "" {
			result[k] = v
		}
	}
	return result
}

// ApplyEnvMap renames keys in secrets according to the mapping.
// If disabled or no mappings, returns the original map unchanged.
func ApplyEnvMap(cfg EnvMapConfig, secrets map[string]string) map[string]string {
	if !cfg.Enabled || len(cfg.Mappings) == 0 {
		return secrets
	}

	result := make(map[string]string, len(secrets))
	for k, v := range secrets {
		if newKey, ok := cfg.Mappings[k]; ok {
			result[newKey] = v
		} else {
			result[k] = v
		}
	}
	return result
}
