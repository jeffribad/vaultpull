package sync

import (
	"testing"
)

func TestTagConfigFromEnv_Defaults(t *testing.T) {
	t.Setenv("VAULTPULL_TAG_FILTER_ENABLED", "")
	t.Setenv("VAULTPULL_REQUIRED_TAGS", "")
	t.Setenv("VAULTPULL_EXCLUDE_TAGS", "")

	cfg := TagConfigFromEnv()
	if cfg.Enabled {
		t.Error("expected Enabled to be false by default")
	}
	if len(cfg.RequiredTags) != 0 {
		t.Errorf("expected no required tags, got %v", cfg.RequiredTags)
	}
	if len(cfg.ExcludeTags) != 0 {
		t.Errorf("expected no exclude tags, got %v", cfg.ExcludeTags)
	}
}

func TestTagConfigFromEnv_Enabled(t *testing.T) {
	t.Setenv("VAULTPULL_TAG_FILTER_ENABLED", "true")
	cfg := TagConfigFromEnv()
	if !cfg.Enabled {
		t.Error("expected Enabled to be true")
	}
}

func TestTagConfigFromEnv_NumericEnabled(t *testing.T) {
	t.Setenv("VAULTPULL_TAG_FILTER_ENABLED", "1")
	cfg := TagConfigFromEnv()
	if !cfg.Enabled {
		t.Error("expected Enabled=true for value '1'")
	}
}

func TestTagConfigFromEnv_ParsesTags(t *testing.T) {
	t.Setenv("VAULTPULL_TAG_FILTER_ENABLED", "true")
	t.Setenv("VAULTPULL_REQUIRED_TAGS", "prod, internal")
	t.Setenv("VAULTPULL_EXCLUDE_TAGS", "deprecated")

	cfg := TagConfigFromEnv()
	if len(cfg.RequiredTags) != 2 {
		t.Errorf("expected 2 required tags, got %v", cfg.RequiredTags)
	}
	if cfg.RequiredTags[0] != "prod" || cfg.RequiredTags[1] != "internal" {
		t.Errorf("unexpected required tags: %v", cfg.RequiredTags)
	}
	if len(cfg.ExcludeTags) != 1 || cfg.ExcludeTags[0] != "deprecated" {
		t.Errorf("unexpected exclude tags: %v", cfg.ExcludeTags)
	}
}
