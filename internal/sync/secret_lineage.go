package sync

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// LineageConfig controls secret lineage tracking (origin metadata injection).
type LineageConfig struct {
	Enabled    bool
	VaultAddr  string
	SecretPath string
	Annotate   bool
}

// LineageConfigFromEnv reads lineage config from environment variables.
func LineageConfigFromEnv() LineageConfig {
	enabled := false
	if v := os.Getenv("VAULTPULL_LINEAGE_ENABLED"); v != "" {
		if b, err := strconv.ParseBool(v); err == nil {
			enabled = b
		} else if v == "1" {
			enabled = true
		}
	}
	return LineageConfig{
		Enabled:    enabled,
		VaultAddr:  os.Getenv("VAULT_ADDR"),
		SecretPath: os.Getenv("VAULTPULL_SECRET_PATH"),
		Annotate:   os.Getenv("VAULTPULL_LINEAGE_ANNOTATE") == "true" || os.Getenv("VAULTPULL_LINEAGE_ANNOTATE") == "1",
	}
}

// LineageRecord holds origin metadata for a synced secret set.
type LineageRecord struct {
	SyncedAt   time.Time
	VaultAddr  string
	SecretPath string
	KeyCount   int
}

// String returns a human-readable lineage summary.
func (r LineageRecord) String() string {
	return fmt.Sprintf("synced %d keys from %s at %s on %s",
		r.KeyCount, r.SecretPath, r.VaultAddr, r.SyncedAt.Format(time.RFC3339))
}

// InjectLineage adds lineage comment lines to the top of a .env output slice.
// Returns nil if disabled or no secrets provided.
func InjectLineage(cfg LineageConfig, secrets map[string]string, lines []string) []string {
	if !cfg.Enabled || !cfg.Annotate || len(secrets) == 0 {
		return lines
	}
	header := []string{
		fmt.Sprintf("# vaultpull lineage: synced %d keys", len(secrets)),
		fmt.Sprintf("# source: %s", strings.TrimRight(cfg.VaultAddr, "/")+"/"+strings.TrimLeft(cfg.SecretPath, "/")),
		fmt.Sprintf("# synced_at: %s", time.Now().UTC().Format(time.RFC3339)),
		"",
	}
	return append(header, lines...)
}
