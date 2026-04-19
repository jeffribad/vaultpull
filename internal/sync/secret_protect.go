package sync

import (
	"fmt"
	"os"
	"strings"
)

// ProtectConfig controls protection of specific secret keys from being overwritten.
type ProtectConfig struct {
	Enabled bool
	Keys    []string
}

// ProtectConfigFromEnv loads protection config from environment variables.
func ProtectConfigFromEnv() ProtectConfig {
	enabled := isTruthy(os.Getenv("VAULTPULL_PROTECT_ENABLED"))
	raw := os.Getenv("VAULTPULL_PROTECT_KEYS")
	keys := splitTrimmed(raw, ",")
	return ProtectConfig{
		Enabled: enabled,
		Keys:    keys,
	}
}

// EnforceProtect checks that protected keys in existing are not overwritten by incoming.
// Returns a list of violation error messages, or nil if none.
func EnforceProtect(cfg ProtectConfig, existing, incoming map[string]string) []error {
	if !cfg.Enabled || len(cfg.Keys) == 0 {
		return nil
	}

	var errs []error
	for _, key := range cfg.Keys {
		norm := strings.ToUpper(strings.TrimSpace(key))
		existingVal, hasExisting := existing[norm]
		if !hasExisting {
			continue
		}
		incomingVal, hasIncoming := incoming[norm]
		if hasIncoming && incomingVal != existingVal {
			errs = append(errs, fmt.Errorf("protected key %q cannot be overwritten", norm))
		}
	}
	if len(errs) == 0 {
		return nil
	}
	return errs
}
