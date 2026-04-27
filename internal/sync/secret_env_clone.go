package sync

import (
	"os"
	"strings"
)

// CloneConfig controls the env-clone feature, which copies a full set of
// secrets under a new key prefix.
type CloneConfig struct {
	Enabled    bool
	FromPrefix string
	ToPrefix   string
	Overwrite  bool
}

// CloneConfigFromEnv reads CloneConfig from environment variables:
//
//	VAULTPULL_CLONE_ENABLED   – "1" / "true" to enable
//	VAULTPULL_CLONE_FROM      – source key prefix (e.g. "PROD_")
//	VAULTPULL_CLONE_TO        – destination key prefix (e.g. "STAGING_")
//	VAULTPULL_CLONE_OVERWRITE – "1" / "true" to overwrite existing destination keys
func CloneConfigFromEnv() CloneConfig {
	return CloneConfig{
		Enabled:    isTruthy(os.Getenv("VAULTPULL_CLONE_ENABLED")),
		FromPrefix: os.Getenv("VAULTPULL_CLONE_FROM"),
		ToPrefix:   os.Getenv("VAULTPULL_CLONE_TO"),
		Overwrite:  isTruthy(os.Getenv("VAULTPULL_CLONE_OVERWRITE")),
	}
}

// ApplyClone copies all secrets whose key starts with cfg.FromPrefix into new
// keys that start with cfg.ToPrefix (the original keys are preserved).
// If cfg.Overwrite is false, existing destination keys are not replaced.
func ApplyClone(secrets map[string]string, cfg CloneConfig) map[string]string {
	if !cfg.Enabled || cfg.FromPrefix == "" || cfg.ToPrefix == "" {
		return secrets
	}

	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		out[k] = v
	}

	for k, v := range secrets {
		upper := strings.ToUpper(k)
		if !strings.HasPrefix(upper, strings.ToUpper(cfg.FromPrefix)) {
			continue
		}
		suffix := k[len(cfg.FromPrefix):]
		dstKey := cfg.ToPrefix + suffix
		if _, exists := out[dstKey]; exists && !cfg.Overwrite {
			continue
		}
		out[dstKey] = v
	}

	return out
}
