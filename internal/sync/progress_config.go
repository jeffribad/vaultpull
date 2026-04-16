package sync

import (
	"os"
	"strings"
)

// ProgressConfig holds configuration for the progress reporter.
type ProgressConfig struct {
	Enabled bool
	Quiet   bool
}

// ProgressConfigFromEnv reads progress config from environment variables.
//
//	VAULTPULL_PROGRESS=true   enables progress output (default: true)
//	VAULTPULL_QUIET=true      suppresses all progress output
func ProgressConfigFromEnv() ProgressConfig {
	cfg := ProgressConfig{
		Enabled: true,
		Quiet:   false,
	}

	if v := os.Getenv("VAULTPULL_PROGRESS"); v != "" {
		norm := strings.ToLower(strings.TrimSpace(v))
		cfg.Enabled = norm == "true" || norm == "1"
	}

	if v := os.Getenv("VAULTPULL_QUIET"); v != "" {
		norm := strings.ToLower(strings.TrimSpace(v))
		cfg.Quiet = norm == "true" || norm == "1"
	}

	return cfg
}
