package sync

import (
	"os"
	"strings"
)

// InheritConfig controls whether secrets inherit values from the current OS environment.
type InheritConfig struct {
	Enabled  bool
	Keys     []string
	Prefix   string
	Override bool // if true, OS env wins over vault
}

// InheritConfigFromEnv loads InheritConfig from environment variables.
func InheritConfigFromEnv() InheritConfig {
	enabled := strings.TrimSpace(os.Getenv("VAULTPULL_INHERIT_ENABLED"))
	keys := splitTrimmed(os.Getenv("VAULTPULL_INHERIT_KEYS"), ",")
	prefix := strings.TrimSpace(os.Getenv("VAULTPULL_INHERIT_PREFIX"))
	override := strings.TrimSpace(os.Getenv("VAULTPULL_INHERIT_OVERRIDE"))

	return InheritConfig{
		Enabled:  enabled == "true" || enabled == "1",
		Keys:     keys,
		Prefix:   prefix,
		Override: override == "true" || override == "1",
	}
}

// ApplyInherit merges OS environment values into the secrets map.
// If cfg.Keys is set, only those keys are considered.
// If cfg.Prefix is set, only OS env vars with that prefix are considered (prefix is stripped).
// If cfg.Override is true, OS env values overwrite existing vault values.
func ApplyInherit(cfg InheritConfig, secrets map[string]string) map[string]string {
	if !cfg.Enabled {
		return secrets
	}

	result := make(map[string]string, len(secrets))
	for k, v := range secrets {
		result[k] = v
	}

	if len(cfg.Keys) > 0 {
		for _, key := range cfg.Keys {
			envVal, ok := os.LookupEnv(key)
			if !ok {
				continue
			}
			_, exists := result[key]
			if !exists || cfg.Override {
				result[key] = envVal
			}
		}
		return result
	}

	if cfg.Prefix != "" {
		for _, env := range os.Environ() {
			parts := strings.SplitN(env, "=", 2)
			if len(parts) != 2 {
				continue
			}
			envKey, envVal := parts[0], parts[1]
			if !strings.HasPrefix(envKey, cfg.Prefix) {
				continue
			}
			stripped := strings.TrimPrefix(envKey, cfg.Prefix)
			if stripped == "" {
				continue
			}
			_, exists := result[stripped]
			if !exists || cfg.Override {
				result[stripped] = envVal
			}
		}
		return result
	}

	return result
}
