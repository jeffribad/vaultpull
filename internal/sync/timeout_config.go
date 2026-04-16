package sync

import (
	"os"
	"strconv"
	"time"
)

// TimeoutConfigFromEnv reads timeout settings from environment variables.
//
// VAULTPULL_VAULT_TIMEOUT_SEC  - seconds for vault read operations (default 10)
// VAULTPULL_WRITE_TIMEOUT_SEC  - seconds for file write operations (default 5)
// VAULTPULL_GLOBAL_TIMEOUT_SEC - seconds for the entire sync run (default 30)
func TimeoutConfigFromEnv() TimeoutConfig {
	cfg := DefaultTimeoutConfig()

	if v := os.Getenv("VAULTPULL_VAULT_TIMEOUT_SEC"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			cfg.VaultTimeout = time.Duration(n) * time.Second
		}
	}
	if v := os.Getenv("VAULTPULL_WRITE_TIMEOUT_SEC"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			cfg.WriteTimeout = time.Duration(n) * time.Second
		}
	}
	if v := os.Getenv("VAULTPULL_GLOBAL_TIMEOUT_SEC"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			cfg.GlobalTimeout = time.Duration(n) * time.Second
		}
	}
	return cfg
}
