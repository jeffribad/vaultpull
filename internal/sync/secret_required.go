package sync

import (
	"fmt"
	"os"
	"strings"
)

// RequiredConfigFromEnv loads required-key enforcement config from environment.
type RequiredConfig struct {
	Enabled  bool
	Keys     []string
	FailFast bool
}

func RequiredConfigFromEnv() RequiredConfig {
	keys := os.Getenv("VAULTPULL_REQUIRED_KEYS")
	failFast := os.Getenv("VAULTPULL_REQUIRED_FAIL_FAST")

	cfg := RequiredConfig{
		Enabled:  len(strings.TrimSpace(keys)) > 0,
		FailFast: isTruthy(failFast),
	}

	if cfg.Enabled {
		cfg.Keys = splitTrimmed(keys, ",")
	}

	return cfg
}

// EnforceRequired checks that all required keys are present and non-empty.
// Returns a list of violation messages. If failFast is true, stops at first.
func EnforceRequired(cfg RequiredConfig, secrets map[string]string) []string {
	if !cfg.Enabled || len(cfg.Keys) == 0 {
		return nil
	}

	var violations []string

	for _, key := range cfg.Keys {
		val, ok := secrets[key]
		if !ok {
			violations = append(violations, fmt.Sprintf("required key %q is missing", key))
		} else if strings.TrimSpace(val) == "" {
			violations = append(violations, fmt.Sprintf("required key %q is empty", key))
		}

		if cfg.FailFast && len(violations) > 0 {
			return violations
		}
	}

	return violations
}
