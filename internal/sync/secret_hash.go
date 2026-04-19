package sync

import (
	"crypto/sha256"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// HashConfig controls per-key hashing of secret values.
type HashConfig struct {
	Enabled bool
	Keys    []string // keys to hash; empty means none
}

// HashConfigFromEnv loads HashConfig from environment variables.
//
//	VAULTPULL_HASH_ENABLED=true
//	VAULTPULL_HASH_KEYS=API_SECRET,DB_PASSWORD
func HashConfigFromEnv() HashConfig {
	enabled := false
	raw := strings.TrimSpace(os.Getenv("VAULTPULL_HASH_ENABLED"))
	if v, err := strconv.ParseBool(raw); err == nil {
		enabled = v
	} else if raw == "1" {
		enabled = true
	}

	keys := splitTrimmed(os.Getenv("VAULTPULL_HASH_KEYS"), ",")
	return HashConfig{Enabled: enabled, Keys: keys}
}

// HashSecrets replaces the value of each configured key with its SHA-256 hex digest.
// If cfg.Keys is empty, no keys are hashed. Returns a new map; input is not mutated.
func HashSecrets(cfg HashConfig, secrets map[string]string) map[string]string {
	if !cfg.Enabled || len(cfg.Keys) == 0 {
		return secrets
	}

	keySet := make(map[string]struct{}, len(cfg.Keys))
	for _, k := range cfg.Keys {
		keySet[strings.ToUpper(k)] = struct{}{}
	}

	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		if _, ok := keySet[strings.ToUpper(k)]; ok {
			sum := sha256.Sum256([]byte(v))
			out[k] = fmt.Sprintf("%x", sum)
		} else {
			out[k] = v
		}
	}
	return out
}
