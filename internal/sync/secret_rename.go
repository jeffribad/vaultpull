package sync

import (
	"os"
	"strings"
)

// RenameMap holds key rename rules: original -> desired
type RenameMap map[string]string

// RenameConfigFromEnv reads VAULTPULL_RENAME_KEYS=OLD_KEY:NEW_KEY,... and returns a RenameMap.
func RenameConfigFromEnv() RenameMap {
	raw := os.Getenv("VAULTPULL_RENAME_KEYS")
	return ParseRenameMap(raw)
}

// ParseRenameMap parses a comma-separated list of old:new pairs.
func ParseRenameMap(raw string) RenameMap {
	rm := make(RenameMap)
	if raw == "" {
		return rm
	}
	for _, pair := range strings.Split(raw, ",") {
		parts := strings.SplitN(strings.TrimSpace(pair), ":", 2)
		if len(parts) != 2 {
			continue
		}
		old := strings.TrimSpace(parts[0])
		new := strings.TrimSpace(parts[1])
		if old != "" && new != "" {
			rm[old] = new
		}
	}
	return rm
}

// ApplyRenames returns a new map with keys renamed according to rm.
// Keys not in rm are passed through unchanged.
func ApplyRenames(secrets map[string]string, rm RenameMap) map[string]string {
	if len(rm) == 0 {
		return secrets
	}
	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		if renamed, ok := rm[k]; ok {
			out[renamed] = v
		} else {
			out[k] = v
		}
	}
	return out
}
