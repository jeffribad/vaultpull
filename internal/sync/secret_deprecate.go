package sync

import (
	"fmt"
	"os"
	"strings"
)

// DeprecateConfig controls deprecated key warnings.
type DeprecateConfig struct {
	Enabled     bool
	Deprecated  map[string]string // old key -> replacement key (or empty)
	FailOnUsage bool
}

// DeprecateConfigFromEnv loads deprecation config from environment variables.
// VAULTPULL_DEPRECATE_ENABLED=1
// VAULTPULL_DEPRECATED_KEYS=OLD_KEY:NEW_KEY,ANOTHER_OLD:
// VAULTPULL_DEPRECATE_FAIL=1
func DeprecateConfigFromEnv() DeprecateConfig {
	cfg := DeprecateConfig{
		Deprecated: make(map[string]string),
	}

	cfg.Enabled = isTruthy(os.Getenv("VAULTPULL_DEPRECATE_ENABLED"))
	cfg.FailOnUsage = isTruthy(os.Getenv("VAULTPULL_DEPRECATE_FAIL"))

	raw := strings.TrimSpace(os.Getenv("VAULTPULL_DEPRECATED_KEYS"))
	if raw == "" {
		return cfg
	}

	for _, pair := range strings.Split(raw, ",") {
		parts := strings.SplitN(strings.TrimSpace(pair), ":", 2)
		if len(parts) != 2 {
			continue
		}
		old := strings.TrimSpace(parts[0])
		new := strings.TrimSpace(parts[1])
		if old != "" {
			cfg.Deprecated[strings.ToUpper(old)] = strings.ToUpper(new)
		}
	}

	return cfg
}

// DeprecationViolation describes a deprecated key found in the secret set.
type DeprecationViolation struct {
	Key         string
	Replacement string
}

func (v DeprecationViolation) Error() string {
	if v.Replacement != "" {
		return fmt.Sprintf("deprecated key %q: use %q instead", v.Key, v.Replacement)
	}
	return fmt.Sprintf("deprecated key %q has no replacement and should be removed", v.Key)
}

// CheckDeprecated inspects secrets for deprecated keys and returns violations.
// If FailOnUsage is true, it returns the first violation as an error.
func CheckDeprecated(cfg DeprecateConfig, secrets map[string]string) ([]DeprecationViolation, error) {
	if !cfg.Enabled || len(cfg.Deprecated) == 0 {
		return nil, nil
	}

	var violations []DeprecationViolation
	for key := range secrets {
		upper := strings.ToUpper(key)
		if replacement, ok := cfg.Deprecated[upper]; ok {
			v := DeprecationViolation{Key: key, Replacement: replacement}
			violations = append(violations, v)
			if cfg.FailOnUsage {
				return violations, v
			}
		}
	}
	return violations, nil
}
