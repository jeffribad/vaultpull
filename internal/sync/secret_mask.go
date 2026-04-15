package sync

import (
	"regexp"
	"strings"
)

// MaskConfig controls how secret values are masked in output.
type MaskConfig struct {
	// ShowChars is the number of leading characters to reveal.
	ShowChars int
	// MaskChar is the character used for masking.
	MaskChar string
}

// DefaultMaskConfig returns a sensible default masking configuration.
func DefaultMaskConfig() MaskConfig {
	return MaskConfig{
		ShowChars: 4,
		MaskChar:  "*",
	}
}

// sensitiveKeyPattern matches common secret key names.
var sensitiveKeyPattern = regexp.MustCompile(
	`(?i)(password|secret|token|key|api_key|apikey|credential|private|auth|passwd|pwd)`,
)

// IsSensitiveKey returns true if the key name looks like it holds a secret value.
func IsSensitiveKey(key string) bool {
	return sensitiveKeyPattern.MatchString(key)
}

// MaskValue masks a secret value according to the given config.
// If the value is shorter than ShowChars, the entire value is masked.
func MaskValue(value string, cfg MaskConfig) string {
	if len(value) == 0 {
		return ""
	}
	show := cfg.ShowChars
	if show >= len(value) {
		return strings.Repeat(cfg.MaskChar, len(value))
	}
	masked := len(value) - show
	return value[:show] + strings.Repeat(cfg.MaskChar, masked)
}

// MaskSecrets returns a copy of the secrets map with sensitive values masked.
func MaskSecrets(secrets map[string]string, cfg MaskConfig) map[string]string {
	result := make(map[string]string, len(secrets))
	for k, v := range secrets {
		if IsSensitiveKey(k) {
			result[k] = MaskValue(v, cfg)
		} else {
			result[k] = v
		}
	}
	return result
}
