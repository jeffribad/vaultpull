package sync

import (
	"os"
	"strings"
)

// CoerceConfig controls type coercion of secret values to normalized forms.
type CoerceConfig struct {
	Enabled    bool
	BoolKeys   []string
	NumberKeys []string
	JSONKeys   []string
}

// CoerceConfigFromEnv loads coercion config from environment variables.
func CoerceConfigFromEnv() CoerceConfig {
	return CoerceConfig{
		Enabled:    isTruthy(os.Getenv("VAULTPULL_COERCE_ENABLED")),
		BoolKeys:   splitTrimmed(os.Getenv("VAULTPULL_COERCE_BOOL_KEYS"), ","),
		NumberKeys: splitTrimmed(os.Getenv("VAULTPULL_COERCE_NUMBER_KEYS"), ","),
		JSONKeys:   splitTrimmed(os.Getenv("VAULTPULL_COERCE_JSON_KEYS"), ","),
	}
}

// CoerceSecrets normalizes secret values into canonical string representations
// based on their declared type (bool, number, json).
func CoerceSecrets(secrets map[string]string, cfg CoerceConfig) map[string]string {
	if !cfg.Enabled {
		return secrets
	}

	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		out[k] = v
	}

	for _, key := range cfg.BoolKeys {
		val, ok := findKey(out, key)
		if !ok {
			continue
		}
		switch strings.ToLower(strings.TrimSpace(val)) {
		case "1", "true", "yes", "on":
			out[key] = "true"
		default:
			out[key] = "false"
		}
	}

	for _, key := range cfg.NumberKeys {
		val, ok := findKey(out, key)
		if !ok {
			continue
		}
		out[key] = strings.TrimSpace(val)
	}

	for _, key := range cfg.JSONKeys {
		val, ok := findKey(out, key)
		if !ok {
			continue
		}
		compact := strings.Join(strings.Fields(val), " ")
		out[key] = compact
	}

	return out
}

// findKey looks up a key in m using exact match first, then case-insensitive
// fallback. Returns the value and whether any match was found.
func findKey(m map[string]string, key string) (string, bool) {
	if v, ok := m[key]; ok {
		return v, true
	}
	lower := strings.ToLower(key)
	for k, v := range m {
		if strings.ToLower(k) == lower {
			return v, true
		}
	}
	return "", false
}

// HasCoercionKeys reports whether cfg declares any keys to coerce.
// This is useful for skipping coercion setup when no keys are configured.
func (cfg CoerceConfig) HasCoercionKeys() bool {
	return len(cfg.BoolKeys) > 0 || len(cfg.NumberKeys) > 0 || len(cfg.JSONKeys) > 0
}
