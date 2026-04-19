package sync

import (
	"fmt"
	"os"
	"strings"
)

// LockConfig controls which secrets are locked to specific values.
type LockConfig struct {
	Enabled bool
	LockedKeys map[string]string // key -> expected value
}

// LockConfigFromEnv reads lock configuration from environment variables.
// VAULTPULL_LOCK_ENABLED=true
// VAULTPULL_LOCK_KEYS=KEY1=val1,KEY2=val2
func LockConfigFromEnv() LockConfig {
	enabled := strings.TrimSpace(os.Getenv("VAULTPULL_LOCK_ENABLED"))
	raw := strings.TrimSpace(os.Getenv("VAULTPULL_LOCK_KEYS"))

	cfg := LockConfig{
		Enabled:    enabled == "true" || enabled == "1",
		LockedKeys: make(map[string]string),
	}

	if raw == "" {
		return cfg
	}

	for _, pair := range strings.Split(raw, ",") {
		parts := strings.SplitN(strings.TrimSpace(pair), "=", 2)
		if len(parts) != 2 {
			continue
		}
		cfg.LockedKeys[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
	}

	return cfg
}

// EnforceLock checks that locked keys in secrets match their expected values.
// Returns an error if any locked key has an unexpected value.
func EnforceLock(cfg LockConfig, secrets map[string]string) error {
	if !cfg.Enabled || len(cfg.LockedKeys) == 0 {
		return nil
	}

	for key, expected := range cfg.LockedKeys {
		val, ok := secrets[key]
		if !ok {
			continue
		}
		if val != expected {
			return fmt.Errorf("lock violation: key %q expected %q but got %q", key, expected, val)
		}
	}

	return nil
}
