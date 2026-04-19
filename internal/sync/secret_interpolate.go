package sync

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

// InterpolateConfig controls variable interpolation within secret values.
type InterpolateConfig struct {
	Enabled bool
	AllowEnv bool
}

// InterpolateConfigFromEnv loads interpolation config from environment variables.
func InterpolateConfigFromEnv() InterpolateConfig {
	return InterpolateConfig{
		Enabled:  isTruthy(os.Getenv("VAULTPULL_INTERPOLATE_ENABLED")),
		AllowEnv: isTruthy(os.Getenv("VAULTPULL_INTERPOLATE_ALLOW_ENV")),
	}
}

var interpolatePattern = regexp.MustCompile(`\$\{([^}]+)\}`)

// InterpolateSecrets replaces ${KEY} references in values with other secret values
// or environment variables (if AllowEnv is set).
func InterpolateSecrets(secrets map[string]string, cfg InterpolateConfig) (map[string]string, error) {
	if !cfg.Enabled {
		return secrets, nil
	}

	result := make(map[string]string, len(secrets))
	for k, v := range secrets {
		result[k] = v
	}

	for key, val := range result {
		interpolated, err := interpolateValue(val, result, cfg.AllowEnv)
		if err != nil {
			return nil, fmt.Errorf("interpolating %q: %w", key, err)
		}
		result[key] = interpolated
	}

	return result, nil
}

func interpolateValue(val string, secrets map[string]string, allowEnv bool) (string, error) {
	var resolveErr error
	result := interpolatePattern.ReplaceAllStringFunc(val, func(match string) string {
		if resolveErr != nil {
			return match
		}
		ref := strings.TrimSpace(match[2 : len(match)-1])
		if v, ok := secrets[ref]; ok {
			return v
		}
		if allowEnv {
			if v, ok := os.LookupEnv(ref); ok {
				return v
			}
		}
		resolveErr = fmt.Errorf("unresolved reference: %q", ref)
		return match
	})
	return result, resolveErr
}
