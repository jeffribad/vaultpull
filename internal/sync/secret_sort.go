package sync

import (
	"os"
	"sort"
	"strings"
)

// SortConfig controls how secrets are sorted before writing.
type SortConfig struct {
	Enabled   bool
	Field     string // "key" or "value"
	Direction string // "asc" or "desc"
}

// SortConfigFromEnv loads sort configuration from environment variables.
func SortConfigFromEnv() SortConfig {
	enabled := os.Getenv("VAULTPULL_SORT_ENABLED")
	field := os.Getenv("VAULTPULL_SORT_FIELD")
	direction := os.Getenv("VAULTPULL_SORT_DIRECTION")

	if field == "" {
		field = "key"
	}
	if direction == "" {
		direction = "asc"
	}

	return SortConfig{
		Enabled:   enabled == "true" || enabled == "1",
		Field:     strings.ToLower(field),
		Direction: strings.ToLower(direction),
	}
}

// ApplySort returns a new map with keys sorted into a slice of key-value pairs.
// Since maps are unordered, this returns an ordered []string of keys.
func ApplySort(secrets map[string]string, cfg SortConfig) []string {
	keys := make([]string, 0, len(secrets))
	for k := range secrets {
		keys = append(keys, k)
	}

	if !cfg.Enabled {
		sort.Strings(keys)
		return keys
	}

	switch cfg.Field {
	case "value":
		sort.Slice(keys, func(i, j int) bool {
			if cfg.Direction == "desc" {
				return secrets[keys[i]] > secrets[keys[j]]
			}
			return secrets[keys[i]] < secrets[keys[j]]
		})
	default: // "key"
		sort.Slice(keys, func(i, j int) bool {
			if cfg.Direction == "desc" {
				return keys[i] > keys[j]
			}
			return keys[i] < keys[j]
		})
	}

	return keys
}
