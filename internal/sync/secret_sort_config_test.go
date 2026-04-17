package sync

import (
	"testing"
)

func TestSortConfigFromEnv_Defaults(t *testing.T) {
	t.Setenv("VAULTPULL_SORT_ENABLED", "")
	t.Setenv("VAULTPULL_SORT_FIELD", "")
	t.Setenv("VAULTPULL_SORT_DIRECTION", "")

	cfg := SortConfigFromEnv()
	if cfg.Enabled {
		t.Error("expected Enabled=false by default")
	}
	if cfg.Field != "key" {
		t.Errorf("expected Field=key, got %s", cfg.Field)
	}
	if cfg.Direction != "asc" {
		t.Errorf("expected Direction=asc, got %s", cfg.Direction)
	}
}

func TestSortConfigFromEnv_Enabled(t *testing.T) {
	t.Setenv("VAULTPULL_SORT_ENABLED", "true")
	cfg := SortConfigFromEnv()
	if !cfg.Enabled {
		t.Error("expected Enabled=true")
	}
}

func TestSortConfigFromEnv_NumericEnabled(t *testing.T) {
	t.Setenv("VAULTPULL_SORT_ENABLED", "1")
	cfg := SortConfigFromEnv()
	if !cfg.Enabled {
		t.Error("expected Enabled=true for value '1'")
	}
}

func TestSortConfigFromEnv_CustomValues(t *testing.T) {
	t.Setenv("VAULTPULL_SORT_ENABLED", "true")
	t.Setenv("VAULTPULL_SORT_FIELD", "value")
	t.Setenv("VAULTPULL_SORT_DIRECTION", "desc")

	cfg := SortConfigFromEnv()
	if cfg.Field != "value" {
		t.Errorf("expected Field=value, got %s", cfg.Field)
	}
	if cfg.Direction != "desc" {
		t.Errorf("expected Direction=desc, got %s", cfg.Direction)
	}
}
