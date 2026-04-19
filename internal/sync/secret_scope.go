package sync

import (
	"os"
	"strings"
)

// ScopeConfig controls filtering secrets by a named scope prefix.
type ScopeConfig struct {
	Enabled bool
	Scope   string
	Strip   bool // strip the scope prefix from the key
}

// ScopeConfigFromEnv loads scope filter config from environment variables.
func ScopeConfigFromEnv() ScopeConfig {
	scope := strings.TrimSpace(os.Getenv("VAULTPULL_SCOPE"))
	strip := isTruthy(os.Getenv("VAULTPULL_SCOPE_STRIP"))
	enabled := scope != ""
	return ScopeConfig{
		Enabled: enabled,
		Scope:   scope,
		Strip:   strip,
	}
}

// ApplyScope filters secrets to only those whose keys start with the given scope.
// If Strip is true, the scope prefix is removed from the resulting keys.
func ApplyScope(cfg ScopeConfig, secrets map[string]string) map[string]string {
	if !cfg.Enabled || cfg.Scope == "" {
		return secrets
	}

	prefix := cfg.Scope
	if !strings.HasSuffix(prefix, "_") {
		prefix += "_"
	}
	upper := strings.ToUpper(prefix)

	out := make(map[string]string)
	for k, v := range secrets {
		if strings.HasPrefix(strings.ToUpper(k), upper) {
			key := k
			if cfg.Strip {
				key = k[len(prefix):]
				if len(key) == 0 {
					continue
				}
			}
			out[key] = v
		}
	}
	return out
}
