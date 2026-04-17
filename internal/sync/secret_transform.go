package sync

import (
	"fmt"
	"os"
	"strings"
)

// TransformConfig holds rules for transforming secret values.
type TransformConfig struct {
	Prefix    string
	Suffix    string
	UpperCase bool
	LowerCase bool
}

// TransformConfigFromEnv loads transform config from environment variables.
func TransformConfigFromEnv() TransformConfig {
	return TransformConfig{
		Prefix:    os.Getenv("VAULTPULL_KEY_PREFIX"),
		Suffix:    os.Getenv("VAULTPULL_KEY_SUFFIX"),
		UpperCase: os.Getenv("VAULTPULL_KEY_UPPERCASE") == "true" || os.Getenv("VAULTPULL_KEY_UPPERCASE") == "1",
		LowerCase: os.Getenv("VAULTPULL_KEY_LOWERCASE") == "true" || os.Getenv("VAULTPULL_KEY_LOWERCASE") == "1",
	}
}

// ApplyTransforms returns a new map with keys transformed according to config.
func ApplyTransforms(secrets map[string]string, cfg TransformConfig) (map[string]string, error) {
	if cfg.UpperCase && cfg.LowerCase {
		return nil, fmt.Errorf("transform: UPPERCASE and LOWERCASE are mutually exclusive")
	}

	result := make(map[string]string, len(secrets))
	for k, v := range secrets {
		newKey := k
		if cfg.Prefix != "" {
			newKey = cfg.Prefix + newKey
		}
		if cfg.Suffix != "" {
			newKey = newKey + cfg.Suffix
		}
		if cfg.UpperCase {
			newKey = strings.ToUpper(newKey)
		} else if cfg.LowerCase {
			newKey = strings.ToLower(newKey)
		}
		result[newKey] = v
	}
	return result, nil
}
