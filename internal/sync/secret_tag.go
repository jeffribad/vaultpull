package sync

import (
	"os"
	"strings"
)

// TagConfig controls tag-based filtering of secrets.
type TagConfig struct {
	Enabled      bool
	RequiredTags []string
	ExcludeTags  []string
}

// TagConfigFromEnv loads tag filter config from environment variables.
func TagConfigFromEnv() TagConfig {
	enabled := strings.TrimSpace(os.Getenv("VAULTPULL_TAG_FILTER_ENABLED"))
	cfg := TagConfig{
		Enabled:      enabled == "true" || enabled == "1",
		RequiredTags: splitTrimmed(os.Getenv("VAULTPULL_REQUIRED_TAGS"), ","),
		ExcludeTags:  splitTrimmed(os.Getenv("VAULTPULL_EXCLUDE_TAGS"), ","),
	}
	return cfg
}

// ApplyTagFilter filters secrets based on their "_tags" metadata field.
// Secrets are expected to optionally carry a "_tags" key with comma-separated tag values.
func ApplyTagFilter(secrets map[string]string, cfg TagConfig) map[string]string {
	if !cfg.Enabled {
		return secrets
	}

	result := make(map[string]string)
	for k, v := range secrets {
		if k == "_tags" {
			continue
		}
		result[k] = v
	}

	rawTags, ok := secrets["_tags"]
	if !ok {
		if len(cfg.RequiredTags) > 0 {
			return map[string]string{}
		}
		return result
	}

	presentTags := splitTrimmed(rawTags, ",")
	tagSet := make(map[string]struct{}, len(presentTags))
	for _, t := range presentTags {
		tagSet[strings.ToLower(t)] = struct{}{}
	}

	for _, req := range cfg.RequiredTags {
		if _, found := tagSet[strings.ToLower(req)]; !found {
			return map[string]string{}
		}
	}

	for _, ex := range cfg.ExcludeTags {
		if _, found := tagSet[strings.ToLower(ex)]; found {
			return map[string]string{}
		}
	}

	return result
}
