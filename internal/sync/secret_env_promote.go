package sync

import (
	"os"
	"strings"
)

// PromoteConfig controls promotion of secrets from one environment tier to another.
type PromoteConfig struct {
	Enabled    bool
	FromPrefix string
	ToPrefix   string
	Overwrite  bool
}

// PromoteConfigFromEnv loads promotion config from environment variables.
func PromoteConfigFromEnv() PromoteConfig {
	enabled := isTruthy(os.Getenv("VAULTPULL_PROMOTE_ENABLED"))
	from := strings.TrimSpace(os.Getenv("VAULTPULL_PROMOTE_FROM_PREFIX"))
	to := strings.TrimSpace(os.Getenv("VAULTPULL_PROMOTE_TO_PREFIX"))
	overwrite := isTruthy(os.Getenv("VAULTPULL_PROMOTE_OVERWRITE"))
	return PromoteConfig{
		Enabled:    enabled,
		FromPrefix: from,
		ToPrefix:   to,
		Overwrite:  overwrite,
	}
}

// ApplyPromotion copies secrets whose keys start with FromPrefix into new keys
// with ToPrefix substituted, optionally overwriting existing keys.
func ApplyPromotion(secrets map[string]string, cfg PromoteConfig) map[string]string {
	if !cfg.Enabled || cfg.FromPrefix == "" {
		return secrets
	}

	result := make(map[string]string, len(secrets))
	for k, v := range secrets {
		result[k] = v
	}

	for k, v := range secrets {
		if !strings.HasPrefix(strings.ToUpper(k), strings.ToUpper(cfg.FromPrefix)) {
			continue
		}
		suffix := k[len(cfg.FromPrefix):]
		newKey := cfg.ToPrefix + suffix
		if _, exists := result[newKey]; exists && !cfg.Overwrite {
			continue
		}
		result[newKey] = v
	}

	return result
}
