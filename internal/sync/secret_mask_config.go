package sync

import (
	"fmt"
	"os"
	"strings"
)

// MaskConfig controls how sensitive secret values are masked in output.
type MaskConfig struct {
	Enabled      bool
	CustomKeys   []string
	RevealLength int
}

// MaskConfigFromEnv loads MaskConfig from environment variables.
//
// VAULTPULL_MASK_ENABLED=true        enables masking (default: true)
// VAULTPULL_MASK_KEYS=FOO,BAR        additional keys to treat as sensitive
// VAULTPULL_MASK_REVEAL_LENGTH=4     characters to reveal at end (default: 4)
func MaskConfigFromEnv() MaskConfig {
	enabled := true
	if v := os.Getenv("VAULTPULL_MASK_ENABLED"); v != "" {
		v = strings.ToLower(strings.TrimSpace(v))
		enabled = v == "true" || v == "1" || v == "yes"
	}

	var customKeys []string
	if v := os.Getenv("VAULTPULL_MASK_KEYS"); v != "" {
		for _, k := range strings.Split(v, ",") {
			if t := strings.TrimSpace(k); t != "" {
				customKeys = append(customKeys, strings.ToUpper(t))
			}
		}
	}

	revealLength := 4
	if v := os.Getenv("VAULTPULL_MASK_REVEAL_LENGTH"); v != "" {
		var n int
		if _, err := fmt.Sscanf(v, "%d", &n); err == nil && n >= 0 {
			revealLength = n
		}
	}

	return MaskConfig{
		Enabled:      enabled,
		CustomKeys:   customKeys,
		RevealLength: revealLength,
	}
}

// ApplyMaskConfig returns a DefaultMaskConfig augmented with custom keys and reveal length.
func ApplyMaskConfig(cfg MaskConfig) DefaultMaskConfig {
	return DefaultMaskConfig{
		Enabled:      cfg.Enabled,
		ExtraKeys:    cfg.CustomKeys,
		RevealLength: cfg.RevealLength,
	}
}
