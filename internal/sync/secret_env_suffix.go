package sync

import (
	"os"
	"strings"
)

// SuffixAddConfig controls whether key suffixes are appended to secrets.
type SuffixAddConfig struct {
	Enabled bool
	Suffix  string
	Keys    []string // if empty, applies to all keys
}

// SuffixAddConfigFromEnv loads SuffixAddConfig from environment variables.
//
//	VAULTPULL_SUFFIX_ADD_ENABLED=true
//	VAULTPULL_SUFFIX_ADD_SUFFIX=_V2
//	VAULTPULL_SUFFIX_ADD_KEYS=DB_HOST,DB_PORT  (optional)
func SuffixAddConfigFromEnv() SuffixAddConfig {
	enabled := isTruthy(os.Getenv("VAULTPULL_SUFFIX_ADD_ENABLED"))
	suffix := os.Getenv("VAULTPULL_SUFFIX_ADD_SUFFIX")
	keys := splitTrimmed(os.Getenv("VAULTPULL_SUFFIX_ADD_KEYS"), ",")
	return SuffixAddConfig{
		Enabled: enabled,
		Suffix:  suffix,
		Keys:    keys,
	}
}

// AddKeySuffix returns a new map with the configured suffix appended to
// matching key names. If Keys is empty, all keys are suffixed.
// Returns the original map unchanged when disabled or suffix is empty.
func AddKeySuffix(secrets map[string]string, cfg SuffixAddConfig) map[string]string {
	if !cfg.Enabled || cfg.Suffix == "" {
		return secrets
	}

	applyAll := len(cfg.Keys) == 0
	allowed := make(map[string]bool, len(cfg.Keys))
	for _, k := range cfg.Keys {
		allowed[strings.ToUpper(k)] = true
	}

	result := make(map[string]string, len(secrets))
	for k, v := range secrets {
		if applyAll || allowed[strings.ToUpper(k)] {
			result[k+cfg.Suffix] = v
		} else {
			result[k] = v
		}
	}
	return result
}
