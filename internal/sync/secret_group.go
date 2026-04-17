package sync

import (
	"os"
	"strings"
)

// GroupConfig controls how secrets are grouped into separate output files.
type GroupConfig struct {
	Enabled  bool
	GroupKey string // metadata key used to determine group name
	OutDir   string // directory to write grouped .env files
}

// GroupConfigFromEnv reads group config from environment variables.
func GroupConfigFromEnv() GroupConfig {
	enabled := os.Getenv("VAULTPULL_GROUP_ENABLED")
	groupKey := os.Getenv("VAULTPULL_GROUP_KEY")
	outDir := os.Getenv("VAULTPULL_GROUP_OUT_DIR")

	if groupKey == "" {
		groupKey = "group"
	}
	if outDir == "" {
		outDir = "."
	}

	return GroupConfig{
		Enabled:  enabled == "true" || enabled == "1",
		GroupKey: groupKey,
		OutDir:   outDir,
	}
}

// GroupSecrets partitions a flat secrets map into named groups using a label key.
// Secrets without the label are placed into the "default" group.
func GroupSecrets(secrets map[string]string, labels map[string]map[string]string, groupKey string) map[string]map[string]string {
	groups := make(map[string]map[string]string)

	for k, v := range secrets {
		groupName := "default"
		if lblMap, ok := labels[k]; ok {
			if g, ok := lblMap[groupKey]; ok && strings.TrimSpace(g) != "" {
				groupName = strings.TrimSpace(g)
			}
		}
		if groups[groupName] == nil {
			groups[groupName] = make(map[string]string)
		}
		groups[groupName][k] = v
	}

	return groups
}
