package sync

import (
	"fmt"
	"os"
	"strings"
)

// ImmutableConfig controls which secrets are treated as immutable (cannot be overwritten).
type ImmutableConfig struct {
	Enabled bool
	Keys    []string
}

// ImmutableConfigFromEnv loads immutable config from environment variables.
func ImmutableConfigFromEnv() ImmutableConfig {
	cfg := ImmutableConfig{}

	val := os.Getenv("VAULTPULL_IMMUTABLE_ENABLED")
	cfg.Enabled = val == "true" || val == "1"

	raw := os.Getenv("VAULTPULL_IMMUTABLE_KEYS")
	if raw != "" {
		for _, k := range strings.Split(raw, ",") {
			k = strings.TrimSpace(k)
			if k != "" {
				cfg.Keys = append(cfg.Keys, k)
			}
		}
	}

	return cfg
}

// EnforceImmutable returns an error if any secret in incoming would overwrite
// an existing value in current for a key marked as immutable.
func EnforceImmutable(cfg ImmutableConfig, current, incoming map[string]string) error {
	if !cfg.Enabled || len(cfg.Keys) == 0 {
		return nil
	}

	for _, key := range cfg.Keys {
		existing, hadExisting := current[strings.ToUpper(key)]
		if !hadExisting {
			existing, hadExisting = current[key]
		}
		if !hadExisting {
			continue
		}

		newVal, hasNew := incoming[strings.ToUpper(key)]
		if !hasNew {
			newVal, hasNew = incoming[key]
		}
		if hasNew && newVal != existing {
			return fmt.Errorf("immutable key %q cannot be overwritten", key)
		}
	}

	return nil
}
