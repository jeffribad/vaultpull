package sync

import (
	"fmt"
	"os"
	"strings"
)

// PinConfig controls secret version pinning — ensuring specific keys
// always resolve to an expected value fingerprint (first 8 chars of value).
type PinConfig struct {
	Enabled bool
	Pins    map[string]string // key -> expected value prefix
	FailFast bool
}

func PinConfigFromEnv() PinConfig {
	cfg := PinConfig{
		Pins: make(map[string]string),
	}
	cfg.Enabled = isTruthy(os.Getenv("VAULTPULL_PIN_ENABLED"))
	cfg.FailFast = isTruthy(os.Getenv("VAULTPULL_PIN_FAIL_FAST"))

	raw := os.Getenv("VAULTPULL_PIN_KEYS")
	for _, pair := range splitTrimmed(raw, ",") {
		parts := strings.SplitN(pair, "=", 2)
		if len(parts) == 2 && parts[0] != "" && parts[1] != "" {
			cfg.Pins[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
		}
	}
	return cfg
}

// PinViolation describes a key whose value does not match its pinned prefix.
type PinViolation struct {
	Key      string
	Expected string
	Actual   string
}

func (v PinViolation) Error() string {
	return fmt.Sprintf("pin violation: key %q expected prefix %q, got %q", v.Key, v.Expected, v.Actual)
}

// CheckPins validates that each pinned key's value starts with the expected prefix.
func CheckPins(cfg PinConfig, secrets map[string]string) []error {
	if !cfg.Enabled || len(cfg.Pins) == 0 {
		return nil
	}
	var errs []error
	for key, expected := range cfg.Pins {
		val, ok := secrets[key]
		if !ok {
			continue
		}
		actual := val
		if len(actual) > 8 {
			actual = actual[:8]
		}
		if !strings.HasPrefix(val, expected) {
			v := PinViolation{Key: key, Expected: expected, Actual: actual}
			if cfg.FailFast {
				return []error{v}
			}
			errs = append(errs, v)
		}
	}
	return errs
}
