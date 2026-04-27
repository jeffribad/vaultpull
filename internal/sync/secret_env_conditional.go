package sync

import (
	"os"
	"strings"
)

// ConditionalConfigFromEnv loads conditional inclusion config from environment variables.
// VAULTPULL_CONDITIONAL_ENABLED=true enables the feature.
// VAULTPULL_CONDITIONAL_RULES=KEY:ENV_VAR=VALUE,... defines rules.
// A secret key is included only if the associated env var matches the expected value.
type ConditionalConfig struct {
	Enabled bool
	Rules   map[string]conditionalRule
}

type conditionalRule struct {
	EnvVar   string
	Expected string
}

func ConditionalConfigFromEnv() ConditionalConfig {
	cfg := ConditionalConfig{
		Rules: make(map[string]conditionalRule),
	}

	raw := os.Getenv("VAULTPULL_CONDITIONAL_ENABLED")
	cfg.Enabled = raw == "true" || raw == "1"

	rulesRaw := os.Getenv("VAULTPULL_CONDITIONAL_RULES")
	if rulesRaw == "" {
		return cfg
	}

	for _, entry := range strings.Split(rulesRaw, ",") {
		entry = strings.TrimSpace(entry)
		// format: SECRET_KEY:ENV_VAR=EXPECTED_VALUE
		colon := strings.Index(entry, ":")
		if colon < 0 {
			continue
		}
		secretKey := strings.TrimSpace(entry[:colon])
		rest := strings.TrimSpace(entry[colon+1:])
		eq := strings.Index(rest, "=")
		if eq < 0 {
			continue
		}
		envVar := strings.TrimSpace(rest[:eq])
		expected := strings.TrimSpace(rest[eq+1:])
		if secretKey == "" || envVar == "" {
			continue
		}
		cfg.Rules[strings.ToUpper(secretKey)] = conditionalRule{
			EnvVar:   envVar,
			Expected: expected,
		}
	}

	return cfg
}

// ApplyConditional filters secrets based on conditional env var rules.
// If a secret key has a rule and the env var does not match, the key is excluded.
func ApplyConditional(cfg ConditionalConfig, secrets map[string]string) map[string]string {
	if !cfg.Enabled || len(cfg.Rules) == 0 {
		return secrets
	}

	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		rule, hasRule := cfg.Rules[strings.ToUpper(k)]
		if !hasRule {
			out[k] = v
			continue
		}
		actual := os.Getenv(rule.EnvVar)
		if actual == rule.Expected {
			out[k] = v
		}
	}
	return out
}
