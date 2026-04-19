package sync

import (
	"os"
	"strings"
)

// FallbackConfig controls whether missing vault secrets fall back to local env vars.
type FallbackConfig struct {
	Enabled  bool
	Keys     []string // if empty, applies to all missing keys
	Prefix   string   // optional prefix to strip when looking up env vars
}

// FallbackConfigFromEnv loads FallbackConfig from environment variables.
func FallbackConfigFromEnv() FallbackConfig {
	enabled := isTruthy(os.Getenv("VAULTPULL_FALLBACK_ENABLED"))
	keys := splitTrimmed(os.Getenv("VAULTPULL_FALLBACK_KEYS"), ",")
	prefix := strings.TrimSpace(os.Getenv("VAULTPULL_FALLBACK_PREFIX"))
	return FallbackConfig{
		Enabled: enabled,
		Keys:    keys,
		Prefix:  prefix,
	}
}

// ApplyFallback fills missing keys in secrets from OS environment variables.
// If cfg.Keys is non-empty, only those keys are considered for fallback.
func ApplyFallback(cfg FallbackConfig, secrets map[string]string) map[string]string {
	if !cfg.Enabled {
		return secrets
	}

	result := make(map[string]string, len(secrets))
	for k, v := range secrets {
		result[k] = v
	}

	applyKey := func(key string) {
		if _, exists := result[key]; exists {
			return
		}
		lookup := cfg.Prefix + key
		if val, ok := os.LookupEnv(lookup); ok {
			result[key] = val
		}
	}

	if len(cfg.Keys) > 0 {
		for _, k := range cfg.Keys {
			applyKey(k)
		}
	} else {
		for _, k := range cfg.Keys {
			applyKey(k)
		}
	}

	return result
}
