package sync

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// FlattenConfig controls nested key flattening behavior.
type FlattenConfig struct {
	Enabled   bool
	Separator string
	MaxDepth  int
}

// FlattenConfigFromEnv loads flatten config from environment variables.
func FlattenConfigFromEnv() FlattenConfig {
	enabled := false
	if v := os.Getenv("VAULTPULL_FLATTEN_ENABLED"); v != "" {
		if b, err := strconv.ParseBool(v); err == nil {
			enabled = b
		} else if v == "1" {
			enabled = true
		}
	}

	separator := "__"
	if v := os.Getenv("VAULTPULL_FLATTEN_SEPARATOR"); v != "" {
		separator = v
	}

	maxDepth := 5
	if v := os.Getenv("VAULTPULL_FLATTEN_MAX_DEPTH"); v != "" {
		if d, err := strconv.Atoi(v); err == nil && d > 0 {
			maxDepth = d
		}
	}

	return FlattenConfig{
		Enabled:   enabled,
		Separator: separator,
		MaxDepth:  maxDepth,
	}
}

// FlattenSecrets expands dot-notation or nested keys into flat keys using
// the configured separator. For example, "db.host" becomes "DB__HOST".
func FlattenSecrets(secrets map[string]string, cfg FlattenConfig) (map[string]string, error) {
	if !cfg.Enabled {
		return secrets, nil
	}

	result := make(map[string]string, len(secrets))
	for k, v := range secrets {
		flat := flattenKey(k, cfg.Separator, cfg.MaxDepth)
		if _, exists := result[flat]; exists {
			return nil, fmt.Errorf("flatten: key collision after flattening %q to %q", k, flat)
		}
		result[flat] = v
	}
	return result, nil
}

func flattenKey(key, separator string, maxDepth int) string {
	parts := strings.FieldsFunc(key, func(r rune) bool {
		return r == '.' || r == '/'
	})
	if len(parts) > maxDepth {
		parts = parts[:maxDepth]
	}
	return strings.ToUpper(strings.Join(parts, separator))
}
