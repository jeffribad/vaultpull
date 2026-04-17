package sync

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// ExpiryConfig controls warning behavior for secrets nearing expiration.
type ExpiryConfig struct {
	Enabled       bool
	WarnWithinDays int
}

// ExpiryConfigFromEnv loads expiry config from environment variables.
func ExpiryConfigFromEnv() ExpiryConfig {
	enabled := os.Getenv("VAULTPULL_EXPIRY_WARN_ENABLED")
	warnDays := os.Getenv("VAULTPULL_EXPIRY_WARN_DAYS")

	cfg := ExpiryConfig{
		Enabled:       enabled == "true" || enabled == "1",
		WarnWithinDays: 7,
	}

	if d, err := strconv.Atoi(warnDays); err == nil && d > 0 {
		cfg.WarnWithinDays = d
	}

	return cfg
}

// ExpiryWarning represents a secret that is nearing or past expiration.
type ExpiryWarning struct {
	Key       string
	ExpiresAt time.Time
	Expired   bool
}

func (w ExpiryWarning) String() string {
	if w.Expired {
		return fmt.Sprintf("secret %q has expired (expired at %s)", w.Key, w.ExpiresAt.Format(time.RFC3339))
	}
	return fmt.Sprintf("secret %q expires soon (at %s)", w.Key, w.ExpiresAt.Format(time.RFC3339))
}

// CheckSecretExpiry inspects secrets for expiry metadata and returns warnings.
func CheckSecretExpiry(cfg ExpiryConfig, secrets map[string]string) []ExpiryWarning {
	if !cfg.Enabled {
		return nil
	}

	var warnings []ExpiryWarning
	now := time.Now()
	threshold := now.Add(time.Duration(cfg.WarnWithinDays) * 24 * time.Hour)

	for key, val := range secrets {
		t, err := time.Parse(time.RFC3339, val)
		if err != nil {
			continue
		}
		if t.Before(threshold) {
			warnings = append(warnings, ExpiryWarning{
				Key:       key,
				ExpiresAt: t,
				Expired:   t.Before(now),
			})
		}
	}

	return warnings
}
