package sync

import (
	"os"
	"strings"
)

// RedactConfig controls which keys should have their values redacted in output.
type RedactConfig struct {
	Enabled bool
	Keys    []string
}

// RedactConfigFromEnv loads redaction config from environment variables.
func RedactConfigFromEnv() RedactConfig {
	enabled := os.Getenv("VAULTPULL_REDACT_ENABLED")
	keys := os.Getenv("VAULTPULL_REDACT_KEYS")

	cfg := RedactConfig{
		Enabled: enabled == "true" || enabled == "1",
		Keys:    splitTrimmed(keys, ","),
	}
	return cfg
}

// RedactSecrets returns a copy of secrets with matching values replaced by "[REDACTED]".
func RedactSecrets(secrets map[string]string, cfg RedactConfig) map[string]string {
	result := make(map[string]string, len(secrets))
	for k, v := range secrets {
		result[k] = v
	}

	if !cfg.Enabled || len(cfg.Keys) == 0 {
		return result
	}

	for _, key := range cfg.Keys {
		normKey := strings.ToUpper(strings.TrimSpace(key))
		for k := range result {
			if strings.ToUpper(k) == normKey {
				result[k] = "[REDACTED]"
			}
		}
	}
	return result
}
