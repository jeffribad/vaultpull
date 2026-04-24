package sync

import (
	"os"
	"regexp"
	"strings"
)

// RegexFilterConfig controls key-based regex filtering of secrets.
type RegexFilterConfig struct {
	Enabled      bool
	AllowPattern string
	DenyPattern  string
}

// RegexFilterConfigFromEnv loads regex filter config from environment variables.
func RegexFilterConfigFromEnv() RegexFilterConfig {
	return RegexFilterConfig{
		Enabled:      isTruthy(os.Getenv("VAULTPULL_REGEX_FILTER_ENABLED")),
		AllowPattern: strings.TrimSpace(os.Getenv("VAULTPULL_REGEX_FILTER_ALLOW")),
		DenyPattern:  strings.TrimSpace(os.Getenv("VAULTPULL_REGEX_FILTER_DENY")),
	}
}

// ApplyRegexFilter filters secrets by matching keys against allow/deny regex patterns.
// If an allow pattern is set, only keys matching it are kept.
// If a deny pattern is set, keys matching it are removed.
// Allow is evaluated before deny.
func ApplyRegexFilter(cfg RegexFilterConfig, secrets map[string]string) (map[string]string, error) {
	if !cfg.Enabled {
		return secrets, nil
	}

	var allowRe, denyRe *regexp.Regexp
	var err error

	if cfg.AllowPattern != "" {
		allowRe, err = regexp.Compile(cfg.AllowPattern)
		if err != nil {
			return nil, err
		}
	}

	if cfg.DenyPattern != "" {
		denyRe, err = regexp.Compile(cfg.DenyPattern)
		if err != nil {
			return nil, err
		}
	}

	result := make(map[string]string, len(secrets))
	for k, v := range secrets {
		if allowRe != nil && !allowRe.MatchString(k) {
			continue
		}
		if denyRe != nil && denyRe.MatchString(k) {
			continue
		}
		result[k] = v
	}
	return result, nil
}
