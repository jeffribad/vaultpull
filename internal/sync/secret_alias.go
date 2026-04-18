package sync

import (
	"os"
	"strings"
)

// AliasConfig holds configuration for secret key aliasing.
// Aliases create additional keys pointing to the same value without removing the original.
type AliasConfig struct {
	Enabled bool
	Aliases map[string][]string // original key -> list of alias keys
}

// AliasConfigFromEnv loads alias configuration from environment variables.
// VAULTPULL_ALIAS_ENABLED=1
// VAULTPULL_ALIASES=DB_HOST:DATABASE_HOST,REDIS_PASS:CACHE_PASSWORD
func AliasConfigFromEnv() AliasConfig {
	enabled := os.Getenv("VAULTPULL_ALIAS_ENABLED")
	raw := os.Getenv("VAULTPULL_ALIASES")
	return AliasConfig{
		Enabled: enabled == "1" || strings.EqualFold(enabled, "true"),
		Aliases: ParseAliasMap(raw),
	}
}

// ParseAliasMap parses a comma-separated list of original:alias pairs.
// Multiple aliases for one key can be specified as DB_HOST:HOST1,DB_HOST:HOST2.
func ParseAliasMap(raw string) map[string][]string {
	result := make(map[string][]string)
	if raw == "" {
		return result
	}
	for _, part := range strings.Split(raw, ",") {
		part = strings.TrimSpace(part)
		idx := strings.Index(part, ":")
		if idx <= 0 || idx == len(part)-1 {
			continue
		}
		orig := strings.TrimSpace(part[:idx])
		alias := strings.TrimSpace(part[idx+1:])
		if orig != "" && alias != "" {
			result[orig] = append(result[orig], alias)
		}
	}
	return result
}

// ApplyAliases returns a new map with alias keys injected alongside originals.
func ApplyAliases(secrets map[string]string, cfg AliasConfig) map[string]string {
	if !cfg.Enabled || len(cfg.Aliases) == 0 {
		return secrets
	}
	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		out[k] = v
	}
	for orig, aliases := range cfg.Aliases {
		val, ok := out[orig]
		if !ok {
			continue
		}
		for _, alias := range aliases {
			if _, exists := out[alias]; !exists {
				out[alias] = val
			}
		}
	}
	return out
}
