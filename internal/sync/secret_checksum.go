package sync

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"sort"
	"strings"
)

// ChecksumConfig controls checksum generation for synced secrets.
type ChecksumConfig struct {
	Enabled bool
	OutputPath string
}

// ChecksumConfigFromEnv loads checksum config from environment variables.
func ChecksumConfigFromEnv() ChecksumConfig {
	enabled := strings.TrimSpace(os.Getenv("VAULTPULL_CHECKSUM_ENABLED"))
	path := strings.TrimSpace(os.Getenv("VAULTPULL_CHECKSUM_PATH"))
	if path == "" {
		path = ".env.sha256"
	}
	return ChecksumConfig{
		Enabled:    enabled == "1" || strings.EqualFold(enabled, "true"),
		OutputPath: path,
	}
}

// ComputeChecksum returns a deterministic SHA-256 hex digest of the secrets map.
func ComputeChecksum(secrets map[string]string) string {
	keys := make([]string, 0, len(secrets))
	for k := range secrets {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	h := sha256.New()
	for _, k := range keys {
		fmt.Fprintf(h, "%s=%s\n", k, secrets[k])
	}
	return hex.EncodeToString(h.Sum(nil))
}

// WriteChecksum writes the checksum of secrets to the configured output path.
// Returns the checksum string and any write error.
func WriteChecksum(cfg ChecksumConfig, secrets map[string]string) (string, error) {
	if !cfg.Enabled {
		return "", nil
	}
	checksum := ComputeChecksum(secrets)
	if err := os.WriteFile(cfg.OutputPath, []byte(checksum+"\n"), 0600); err != nil {
		return "", fmt.Errorf("write checksum: %w", err)
	}
	return checksum, nil
}
