package sync

import (
	"fmt"
	"os"
	"strings"
)

// ReadOnlyConfig controls enforcement of read-only keys that must not be written to .env.
type ReadOnlyConfig struct {
	Enabled bool
	Keys    []string
	FailFast bool
}

// ReadOnlyConfigFromEnv loads ReadOnlyConfig from environment variables.
func ReadOnlyConfigFromEnv() ReadOnlyConfig {
	enabled := strings.TrimSpace(os.Getenv("VAULTPULL_READONLY_ENABLED"))
	keys := strings.TrimSpace(os.Getenv("VAULTPULL_READONLY_KEYS"))
	failFast := strings.TrimSpace(os.Getenv("VAULTPULL_READONLY_FAIL_FAST"))

	cfg := ReadOnlyConfig{
		Enabled:  enabled == "true" || enabled == "1",
		FailFast: failFast == "true" || failFast == "1",
	}

	if keys != "" {
		for _, k := range strings.Split(keys, ",") {
			if t := strings.TrimSpace(k); t != "" {
				cfg.Keys = append(cfg.Keys, strings.ToUpper(t))
			}
		}
	}

	return cfg
}

// ReadOnlyViolation describes a key that is marked read-only but present in secrets.
type ReadOnlyViolation struct {
	Key string
}

func (v ReadOnlyViolation) Error() string {
	return fmt.Sprintf("key %q is marked read-only and must not be written", v.Key)
}

// EnforceReadOnly checks that no secret key matches the read-only list.
// If FailFast is true, it returns on the first violation.
func EnforceReadOnly(cfg ReadOnlyConfig, secrets map[string]string) []error {
	if !cfg.Enabled || len(cfg.Keys) == 0 {
		return nil
	}

	readOnly := make(map[string]struct{}, len(cfg.Keys))
	for _, k := range cfg.Keys {
		readOnly[k] = struct{}{}
	}

	var errs []error
	for k := range secrets {
		if _, found := readOnly[strings.ToUpper(k)]; found {
			errs = append(errs, ReadOnlyViolation{Key: k})
			if cfg.FailFast {
				return errs
			}
		}
	}
	return errs
}
