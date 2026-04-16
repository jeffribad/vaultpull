package sync

import (
	"os"
	"strings"
)

// LabelFilterConfig holds key=value label selectors used to filter secrets.
type LabelFilterConfig struct {
	Labels map[string]string
}

// LabelFilterConfigFromEnv reads VAULTPULL_LABELS (e.g. "env=prod,team=backend").
func LabelFilterConfigFromEnv() LabelFilterConfig {
	raw := os.Getenv("VAULTPULL_LABELS")
	return LabelFilterConfig{Labels: parseLabels(raw)}
}

func parseLabels(raw string) map[string]string {
	labels := map[string]string{}
	for _, pair := range splitTrimmed(raw, ",") {
		parts := strings.SplitN(pair, "=", 2)
		if len(parts) == 2 {
			k := strings.TrimSpace(parts[0])
			v := strings.TrimSpace(parts[1])
			if k != "" {
				labels[k] = v
			}
		}
	}
	return labels
}

// ApplyLabelFilter returns only secrets whose metadata labels match ALL selectors.
// secretLabels is a map of secret key -> its labels. If cfg has no labels, all secrets pass.
func ApplyLabelFilter(secrets map[string]string, secretLabels map[string]map[string]string, cfg LabelFilterConfig) map[string]string {
	if len(cfg.Labels) == 0 {
		return secrets
	}
	result := map[string]string{}
	for k, v := range secrets {
		labels := secretLabels[k]
		if matchesAllLabels(labels, cfg.Labels) {
			result[k] = v
		}
	}
	return result
}

func matchesAllLabels(have, want map[string]string) bool {
	for k, v := range want {
		if have[k] != v {
			return false
		}
	}
	return true
}
