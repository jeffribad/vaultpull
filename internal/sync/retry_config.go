package sync

import (
	"os"
	"strconv"
	"time"
)

// RetryConfigFromEnv builds a RetryConfig from environment variables,
// falling back to DefaultRetryConfig() for any missing or invalid values.
//
// Environment variables:
//   VAULTPULL_RETRY_ATTEMPTS  – maximum number of attempts (int, default 3)
//   VAULTPULL_RETRY_DELAY_MS  – initial delay in milliseconds (int, default 500)
//   VAULTPULL_RETRY_MULTIPLIER – backoff multiplier (float, default 2.0)
func RetryConfigFromEnv() RetryConfig {
	cfg := DefaultRetryConfig()

	if v := os.Getenv("VAULTPULL_RETRY_ATTEMPTS"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			cfg.MaxAttempts = n
		}
	}

	if v := os.Getenv("VAULTPULL_RETRY_DELAY_MS"); v != "" {
		if ms, err := strconv.Atoi(v); err == nil && ms >= 0 {
			cfg.Delay = time.Duration(ms) * time.Millisecond
		}
	}

	if v := os.Getenv("VAULTPULL_RETRY_MULTIPLIER"); v != "" {
		if m, err := strconv.ParseFloat(v, 64); err == nil && m >= 1.0 {
			cfg.Multiplier = m
		}
	}

	return cfg
}
