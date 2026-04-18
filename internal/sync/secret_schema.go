package sync

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

// SchemaRule defines a validation rule for a secret key.
type SchemaRule struct {
	Key     string
	Pattern string
	Regexp  *regexp.Regexp
}

// SchemaConfig controls schema-based secret validation.
type SchemaConfig struct {
	Enabled bool
	Rules   []SchemaRule
}

// SchemaConfigFromEnv reads schema config from environment variables.
// VAULTPULL_SCHEMA_ENABLED=1
// VAULTPULL_SCHEMA_RULES=API_KEY:^[A-Za-z0-9]{32}$,DB_URL:^postgres://
func SchemaConfigFromEnv() SchemaConfig {
	enabled := os.Getenv("VAULTPULL_SCHEMA_ENABLED") == "1" ||
		os.Getenv("VAULTPULL_SCHEMA_ENABLED") == "true"

	var rules []SchemaRule
	raw := strings.TrimSpace(os.Getenv("VAULTPULL_SCHEMA_RULES"))
	if raw != "" {
		for _, pair := range strings.Split(raw, ",") {
			parts := strings.SplitN(strings.TrimSpace(pair), ":", 2)
			if len(parts) != 2 {
				continue
			}
			key := strings.TrimSpace(parts[0])
			pattern := strings.TrimSpace(parts[1])
			re, err := regexp.Compile(pattern)
			if err != nil {
				continue
			}
			rules = append(rules, SchemaRule{Key: key, Pattern: pattern, Regexp: re})
		}
	}

	return SchemaConfig{Enabled: enabled, Rules: rules}
}

// ValidateSchema checks secrets against schema rules.
// Returns a slice of validation error strings.
func ValidateSchema(cfg SchemaConfig, secrets map[string]string) []error {
	if !cfg.Enabled || len(cfg.Rules) == 0 {
		return nil
	}

	var errs []error
	for _, rule := range cfg.Rules {
		val, ok := secrets[rule.Key]
		if !ok {
			continue
		}
		if !rule.Regexp.MatchString(val) {
			errs = append(errs, fmt.Errorf("secret %q does not match schema pattern %q", rule.Key, rule.Pattern))
		}
	}
	return errs
}
