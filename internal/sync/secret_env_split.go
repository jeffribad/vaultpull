package sync

import (
	"os"
	"strings"
)

// SplitConfig controls splitting a single secret value into multiple keys
// using a delimiter. For example, splitting "a=1,b=2" into separate keys.
type SplitConfig struct {
	Enabled   bool
	SourceKey string
	Delimiter string
	Separator string // separates key=value within each segment
}

// SplitConfigFromEnv reads split configuration from environment variables.
//
//	VAULTPULL_SPLIT_ENABLED=true
//	VAULTPULL_SPLIT_SOURCE=MULTI_SECRET
//	VAULTPULL_SPLIT_DELIMITER=,
//	VAULTPULL_SPLIT_SEPARATOR==
func SplitConfigFromEnv() SplitConfig {
	return SplitConfig{
		Enabled:   isTruthy(os.Getenv("VAULTPULL_SPLIT_ENABLED")),
		SourceKey: os.Getenv("VAULTPULL_SPLIT_SOURCE"),
		Delimiter: envOrDefault("VAULTPULL_SPLIT_DELIMITER", ","),
		Separator: envOrDefault("VAULTPULL_SPLIT_SEPARATOR", "="),
	}
}

func envOrDefault(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

// ApplySplit reads cfg.SourceKey from secrets, splits its value by
// cfg.Delimiter, and injects the resulting key/value pairs back into a copy
// of the map. The original source key is preserved.
func ApplySplit(secrets map[string]string, cfg SplitConfig) (map[string]string, error) {
	if !cfg.Enabled || cfg.SourceKey == "" {
		return secrets, nil
	}

	// Case-insensitive lookup for source key
	var rawValue string
	found := false
	for k, v := range secrets {
		if strings.EqualFold(k, cfg.SourceKey) {
			rawValue = v
			found = true
			break
		}
	}
	if !found {
		return secrets, nil
	}

	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		out[k] = v
	}

	segments := strings.Split(rawValue, cfg.Delimiter)
	for _, seg := range segments {
		seg = strings.TrimSpace(seg)
		if seg == "" {
			continue
		}
		parts := strings.SplitN(seg, cfg.Separator, 2)
		if len(parts) != 2 {
			continue
		}
		k := strings.TrimSpace(parts[0])
		v := strings.TrimSpace(parts[1])
		if k != "" {
			out[k] = v
		}
	}

	return out, nil
}
