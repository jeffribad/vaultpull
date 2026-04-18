package sync

import (
	"os"
	"strings"
)

// EnvOverrideConfig controls whether local environment variables can override vault secrets.
type EnvOverrideConfig struct {
	Enabled  bool
	Prefix   string
	Priority string // "env" or "vault"
}

// EnvOverrideConfigFromEnv loads override config from environment.
func EnvOverrideConfigFromEnv() EnvOverrideConfig {
	enabled := strings.TrimSpace(os.Getenv("VAULTPULL_ENV_OVERRIDE_ENABLED"))
	prefix := strings.TrimSpace(os.Getenv("VAULTPULL_ENV_OVERRIDE_PREFIX"))
	priority := strings.TrimSpace(os.Getenv("VAULTPULL_ENV_OVERRIDE_PRIORITY"))

	if prefix == "" {
		prefix = "VAULTPULL_OVERRIDE_"
	}
	if priority == "" {
		priority = "env"
	}

	return EnvOverrideConfig{
		Enabled:  enabled == "true" || enabled == "1",
		Prefix:   prefix,
		Priority: priority,
	}
}

// ApplyEnvOverrides merges local environment overrides into the secrets map.
// If priority is "env", local env vars win. If "vault", vault values win.
func ApplyEnvOverrides(cfg EnvOverrideConfig, secrets map[string]string) map[string]string {
	if !cfg.Enabled {
		return secrets
	}

	result := make(map[string]string, len(secrets))
	for k, v := range secrets {
		result[k] = v
	}

	for _, env := range os.Environ() {
		parts := strings.SplitN(env, "=", 2)
		if len(parts) != 2 {
			continue
		}
		envKey, envVal := parts[0], parts[1]
		if !strings.HasPrefix(envKey, cfg.Prefix) {
			continue
		}
		secretKey := strings.TrimPrefix(envKey, cfg.Prefix)
		if secretKey == "" {
			continue
		}
		if cfg.Priority == "vault" {
			if _, exists := result[secretKey]; exists {
				continue
			}
		}
		result[secretKey] = envVal
	}

	return result
}
