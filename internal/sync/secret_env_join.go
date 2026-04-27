package sync

import (
	"fmt"
	"os"
	"strings"
)

// JoinConfig controls how multiple secret values are joined into a single key.
type JoinConfig struct {
	Enabled   bool
	SourceKeys []string
	DestKey   string
	Separator string
}

// JoinConfigFromEnv loads JoinConfig from environment variables.
//
//	VAULTPULL_JOIN_ENABLED=true
//	VAULTPULL_JOIN_SOURCE_KEYS=DB_HOST,DB_PORT
//	VAULTPULL_JOIN_DEST_KEY=DB_ADDR
//	VAULTPULL_JOIN_SEPARATOR=:
func JoinConfigFromEnv() JoinConfig {
	enabled := os.Getenv("VAULTPULL_JOIN_ENABLED")
	sourceRaw := os.Getenv("VAULTPULL_JOIN_SOURCE_KEYS")
	destKey := os.Getenv("VAULTPULL_JOIN_DEST_KEY")
	separator := os.Getenv("VAULTPULL_JOIN_SEPARATOR")

	if separator == "" {
		separator = ","
	}

	return JoinConfig{
		Enabled:    isTruthy(enabled),
		SourceKeys: splitTrimmed(sourceRaw, ","),
		DestKey:    destKey,
		Separator:  separator,
	}
}

// ApplyJoin joins multiple secret values into a single destination key.
// Source keys are resolved case-insensitively. The destination key is added
// to the output map; existing keys are preserved.
func ApplyJoin(secrets map[string]string, cfg JoinConfig) (map[string]string, error) {
	if !cfg.Enabled {
		return secrets, nil
	}
	if len(cfg.SourceKeys) == 0 || cfg.DestKey == "" {
		return secrets, nil
	}

	// Build a lowercase index for case-insensitive lookup.
	lower := make(map[string]string, len(secrets))
	for k, v := range secrets {
		lower[strings.ToLower(k)] = v
	}

	parts := make([]string, 0, len(cfg.SourceKeys))
	for _, src := range cfg.SourceKeys {
		v, ok := lower[strings.ToLower(src)]
		if !ok {
			return nil, fmt.Errorf("join: source key %q not found in secrets", src)
		}
		parts = append(parts, v)
	}

	out := make(map[string]string, len(secrets)+1)
	for k, v := range secrets {
		out[k] = v
	}
	out[cfg.DestKey] = strings.Join(parts, cfg.Separator)
	return out, nil
}
