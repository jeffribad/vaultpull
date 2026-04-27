package sync

import (
	"os"
	"regexp"
	"strings"
)

// MaskPatternConfig controls regex-based value masking for secrets whose keys
// match a configured pattern. Matched values are replaced with a fixed mask.
type MaskPatternConfig struct {
	Enabled     bool
	Pattern     string
	Mask        string
	CompiledRe  *regexp.Regexp
}

// MaskPatternConfigFromEnv loads MaskPatternConfig from environment variables.
//
//	VAULTPULL_MASK_PATTERN_ENABLED  – enable pattern masking (default: false)
//	VAULTPULL_MASK_PATTERN_REGEX    – regex applied to key names
//	VAULTPULL_MASK_PATTERN_MASK     – replacement string (default: "***")
func MaskPatternConfigFromEnv() MaskPatternConfig {
	cfg := MaskPatternConfig{
		Mask: "***",
	}

	if v := os.Getenv("VAULTPULL_MASK_PATTERN_ENABLED"); isTruthy(v) {
		cfg.Enabled = true
	}
	if v := os.Getenv("VAULTPULL_MASK_PATTERN_REGEX"); v != "" {
		cfg.Pattern = v
	}
	if v := os.Getenv("VAULTPULL_MASK_PATTERN_MASK"); v != "" {
		cfg.Mask = v
	}

	if cfg.Enabled && cfg.Pattern != "" {
		re, err := regexp.Compile(cfg.Pattern)
		if err == nil {
			cfg.CompiledRe = re
		}
	}

	return cfg
}

// ApplyMaskPattern replaces the values of any secrets whose keys match the
// configured regex with the mask string. Returns a new map; input is not
// mutated. If the config is disabled or no pattern is set, the original map
// is returned unchanged.
func ApplyMaskPattern(cfg MaskPatternConfig, secrets map[string]string) map[string]string {
	if !cfg.Enabled || cfg.CompiledRe == nil {
		return secrets
	}

	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		if cfg.CompiledRe.MatchString(strings.ToLower(k)) || cfg.CompiledRe.MatchString(k) {
			out[k] = cfg.Mask
		} else {
			out[k] = v
		}
	}
	return out
}
