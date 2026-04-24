package sync

import (
	"os"
	"strings"
)

// WhitelistConfig controls which secret keys are explicitly allowed.
// Keys not in the whitelist are removed from the output.
type WhitelistConfig struct {
	Enabled bool
	Keys    []string
}

// WhitelistConfigFromEnv loads whitelist configuration from environment variables.
//
//	VAULTPULL_WHITELIST_ENABLED=true
//	VAULTPULL_WHITELIST_KEYS=DB_HOST,DB_PORT,API_KEY
func WhitelistConfigFromEnv() WhitelistConfig {
	enabled := strings.TrimSpace(os.Getenv("VAULTPULL_WHITELIST_ENABLED"))
	keys := strings.TrimSpace(os.Getenv("VAULTPULL_WHITELIST_KEYS"))

	cfg := WhitelistConfig{}

	switch strings.ToLower(enabled) {
	case "true", "1", "yes":
		cfg.Enabled = true
	}

	if keys != "" {
		for _, k := range strings.Split(keys, ",") {
			if trimmed := strings.TrimSpace(k); trimmed != "" {
				cfg.Keys = append(cfg.Keys, trimmed)
			}
		}
	}

	return cfg
}

// ApplyWhitelist removes any secrets whose keys are not in the whitelist.
// If the whitelist is disabled or empty, the original map is returned unchanged.
func ApplyWhitelist(cfg WhitelistConfig, secrets map[string]string) map[string]string {
	if !cfg.Enabled || len(cfg.Keys) == 0 {
		return secrets
	}

	allowed := make(map[string]bool, len(cfg.Keys))
	for _, k := range cfg.Keys {
		allowed[strings.ToUpper(k)] = true
	}

	result := make(map[string]string, len(secrets))
	for k, v := range secrets {
		if allowed[strings.ToUpper(k)] {
			result[k] = v
		}
	}
	return result
}
