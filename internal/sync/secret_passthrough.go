package sync

import (
	"os"
	"strings"
)

// PassthroughConfig controls whether local environment variables can
// pass through and override secrets pulled from Vault.
type PassthroughConfig struct {
	Enabled bool
	Keys    []string // specific keys to pass through; empty means all
}

// PassthroughConfigFromEnv reads passthrough configuration from environment variables.
//
//	VAULTPULL_PASSTHROUGH_ENABLED=true
//	VAULTPULL_PASSTHROUGH_KEYS=DB_HOST,DB_PORT
func PassthroughConfigFromEnv() PassthroughConfig {
	enabled := os.Getenv("VAULTPULL_PASSTHROUGH_ENABLED")
	keys := os.Getenv("VAULTPULL_PASSTHROUGH_KEYS")

	cfg := PassthroughConfig{}

	switch strings.ToLower(strings.TrimSpace(enabled)) {
	case "true", "1", "yes":
		cfg.Enabled = true
	}

	if keys != "" {
		cfg.Keys = splitTrimmed(keys, ",")
	}

	return cfg
}

// ApplyPassthrough overrides secrets with local environment variable values
// when passthrough is enabled. If Keys is empty, all env vars that match
// existing secret keys are passed through.
func ApplyPassthrough(cfg PassthroughConfig, secrets map[string]string) map[string]string {
	if !cfg.Enabled {
		return secrets
	}

	result := make(map[string]string, len(secrets))
	for k, v := range secrets {
		result[k] = v
	}

	if len(cfg.Keys) == 0 {
		// pass through any env var that matches an existing secret key
		for k := range result {
			if val, ok := os.LookupEnv(k); ok {
				result[k] = val
			}
		}
		return result
	}

	for _, key := range cfg.Keys {
		if val, ok := os.LookupEnv(key); ok {
			result[key] = val
		}
	}

	return result
}
