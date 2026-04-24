package sync

import (
	"os"
	"strings"
)

// BlacklistConfig controls which keys are explicitly blocked from sync.
type BlacklistConfig struct {
	Enabled bool
	Keys    []string
}

// BlacklistConfigFromEnv loads blacklist configuration from environment variables.
//
//	VAULTPULL_BLACKLIST_ENABLED=true
//	VAULTPULL_BLACKLIST_KEYS=SECRET_KEY,INTERNAL_TOKEN,DEBUG_PASSWORD
func BlacklistConfigFromEnv() BlacklistConfig {
	cfg := BlacklistConfig{}

	val := os.Getenv("VAULTPULL_BLACKLIST_ENABLED")
	cfg.Enabled = val == "true" || val == "1"

	raw := os.Getenv("VAULTPULL_BLACKLIST_KEYS")
	if raw != "" {
		cfg.Keys = splitTrimmed(raw, ",")
	}

	return cfg
}

// ApplyBlacklist removes any secrets whose keys appear in the blacklist.
// Keys are matched case-insensitively. Returns a new map without mutating input.
func ApplyBlacklist(cfg BlacklistConfig, secrets map[string]string) map[string]string {
	if !cfg.Enabled || len(cfg.Keys) == 0 {
		return secrets
	}

	blocked := make(map[string]struct{}, len(cfg.Keys))
	for _, k := range cfg.Keys {
		blocked[strings.ToLower(k)] = struct{}{}
	}

	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		if _, found := blocked[strings.ToLower(k)]; !found {
			out[k] = v
		}
	}
	return out
}
