package sync

import (
	"fmt"
	"os"
	"strings"
)

// SubstituteConfig controls environment variable substitution in secret values.
type SubstituteConfig struct {
	Enabled    bool
	AllowEmpty bool   // if true, unresolved ${VAR} becomes empty string instead of error
	Prefix     string // only substitute vars with this prefix, e.g. "APP_"
}

// SubstituteConfigFromEnv loads SubstituteConfig from environment variables.
func SubstituteConfigFromEnv() SubstituteConfig {
	return SubstituteConfig{
		Enabled:    isTruthy(os.Getenv("VAULTPULL_SUBSTITUTE_ENABLED")),
		AllowEmpty: isTruthy(os.Getenv("VAULTPULL_SUBSTITUTE_ALLOW_EMPTY")),
		Prefix:     strings.TrimSpace(os.Getenv("VAULTPULL_SUBSTITUTE_PREFIX")),
	}
}

// SubstituteSecrets replaces ${VAR} placeholders in secret values with
// corresponding values from the secrets map itself or the OS environment.
func SubstituteSecrets(cfg SubstituteConfig, secrets map[string]string) (map[string]string, error) {
	if !cfg.Enabled {
		return secrets, nil
	}

	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		resolved, err := substituteValue(v, cfg, secrets)
		if err != nil {
			return nil, fmt.Errorf("substitute: key %q: %w", k, err)
		}
		out[k] = resolved
	}
	return out, nil
}

func substituteValue(val string, cfg SubstituteConfig, secrets map[string]string) (string, error) {
	var sb strings.Builder
	s := val
	for {
		start := strings.Index(s, "${")
		if start == -1 {
			sb.WriteString(s)
			break
		}
		end := strings.Index(s[start:], "}")
		if end == -1 {
			sb.WriteString(s)
			break
		}
		end += start
		sb.WriteString(s[:start])
		varName := s[start+2 : end]

		if cfg.Prefix != "" && !strings.HasPrefix(varName, cfg.Prefix) {
			sb.WriteString(s[start : end+1])
			s = s[end+1:]
			continue
		}

		replacement, ok := secrets[varName]
		if !ok {
			replacement = os.Getenv(varName)
			if replacement == "" && !ok {
				if !cfg.AllowEmpty {
					return "", fmt.Errorf("unresolved variable: %s", varName)
				}
			}
		}
		sb.WriteString(replacement)
		s = s[end+1:]
	}
	return sb.String(), nil
}
