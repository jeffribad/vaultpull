package sync

import (
	"os"
	"strconv"
	"strings"
)

// CommentConfig controls whether section comments are injected into .env output.
type CommentConfig struct {
	Enabled bool
	Prefix  string
}

// CommentConfigFromEnv loads comment config from environment variables.
func CommentConfigFromEnv() CommentConfig {
	enabled := false
	if v := os.Getenv("VAULTPULL_COMMENTS_ENABLED"); v != "" {
		if b, err := strconv.ParseBool(v); err == nil {
			enabled = b
		} else if v == "1" {
			enabled = true
		}
	}
	prefix := os.Getenv("VAULTPULL_COMMENTS_PREFIX")
	if prefix == "" {
		prefix = "# -- %s --"
	}
	return CommentConfig{Enabled: enabled, Prefix: prefix}
}

// InjectComments inserts a header comment above each group of keys sharing a
// common prefix (split on "_"). Only the first segment is used as a group label.
func InjectComments(secrets map[string]string, cfg CommentConfig) []string {
	if !cfg.Enabled || len(secrets) == 0 {
		return nil
	}

	keys := sortedSecretKeys(secrets)
	var lines []string
	lastGroup := ""

	for _, k := range keys {
		group := groupLabel(k)
		if group != lastGroup {
			if lastGroup != "" {
				lines = append(lines, "")
			}
			lines = append(lines, formatComment(cfg.Prefix, group))
			lastGroup = group
		}
		lines = append(lines, k+"="+secrets[k])
	}
	return lines
}

func groupLabel(key string) string {
	parts := strings.SplitN(key, "_", 2)
	return strings.ToUpper(parts[0])
}

func formatComment(prefix, label string) string {
	if strings.Contains(prefix, "%s") {
		return strings.Replace(prefix, "%s", label, 1)
	}
	return prefix + " " + label
}
