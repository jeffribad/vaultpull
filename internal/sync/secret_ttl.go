package sync

import (
	"os"
	"strconv"
	"time"
)

// TTLConfig holds configuration for secret TTL enforcement.
type TTLConfig struct {
	Enabled    bool
	MaxAgeDays int
}

// TTLConfigFromEnv loads TTL config from environment variables.
func TTLConfigFromEnv() TTLConfig {
	cfg := TTLConfig{
		MaxAgeDays: 30,
	}
	if v := os.Getenv("VAULTPULL_TTL_ENABLED"); v == "true" || v == "1" {
		cfg.Enabled = true
	}
	if v := os.Getenv("VAULTPULL_TTL_MAX_AGE_DAYS"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			cfg.MaxAgeDays = n
		}
	}
	return cfg
}

// TTLViolation describes a secret that has exceeded its allowed age.
type TTLViolation struct {
	Key     string
	AgeDays int
}

// CheckSecretTTL inspects secret metadata for created_time and returns
// violations for any secrets older than cfg.MaxAgeDays.
func CheckSecretTTL(cfg TTLConfig, secrets map[string]string, createdAt map[string]time.Time) []TTLViolation {
	if !cfg.Enabled {
		return nil
	}
	var violations []TTLViolation
	now := time.Now()
	maxAge := time.Duration(cfg.MaxAgeDays) * 24 * time.Hour
	for key := range secrets {
		t, ok := createdAt[key]
		if !ok {
			continue
		}
		age := now.Sub(t)
		if age > maxAge {
			violations = append(violations, TTLViolation{
				Key:     key,
				AgeDays: int(age.Hours() / 24),
			})
		}
	}
	return violations
}
