package sync

import (
	"testing"
)

func TestApplySort_DefaultAscByKey(t *testing.T) {
	secrets := map[string]string{"ZEBRA": "1", "APPLE": "2", "MANGO": "3"}
	cfg := SortConfig{Enabled: true, Field: "key", Direction: "asc"}
	keys := ApplySort(secrets, cfg)
	if keys[0] != "APPLE" || keys[1] != "MANGO" || keys[2] != "ZEBRA" {
		t.Errorf("unexpected order: %v", keys)
	}
}

func TestApplySort_DescByKey(t *testing.T) {
	secrets := map[string]string{"ZEBRA": "1", "APPLE": "2", "MANGO": "3"}
	cfg := SortConfig{Enabled: true, Field: "key", Direction: "desc"}
	keys := ApplySort(secrets, cfg)
	if keys[0] != "ZEBRA" || keys[1] != "MANGO" || keys[2] != "APPLE" {
		t.Errorf("unexpected order: %v", keys)
	}
}

func TestApplySort_AscByValue(t *testing.T) {
	secrets := map[string]string{"A": "zebra", "B": "apple", "C": "mango"}
	cfg := SortConfig{Enabled: true, Field: "value", Direction: "asc"}
	keys := ApplySort(secrets, cfg)
	if keys[0] != "B" || keys[1] != "C" || keys[2] != "A" {
		t.Errorf("unexpected order by value asc: %v", keys)
	}
}

func TestApplySort_DescByValue(t *testing.T) {
	secrets := map[string]string{"A": "zebra", "B": "apple", "C": "mango"}
	cfg := SortConfig{Enabled: true, Field: "value", Direction: "desc"}
	keys := ApplySort(secrets, cfg)
	if keys[0] != "A" || keys[1] != "C" || keys[2] != "B" {
		t.Errorf("unexpected order by value desc: %v", keys)
	}
}

func TestApplySort_Disabled_StillSortsAlphabetically(t *testing.T) {
	secrets := map[string]string{"Z": "1", "A": "2"}
	cfg := SortConfig{Enabled: false}
	keys := ApplySort(secrets, cfg)
	if keys[0] != "A" || keys[1] != "Z" {
		t.Errorf("expected alphabetical fallback, got %v", keys)
	}
}

func TestApplySort_EmptySecrets(t *testing.T) {
	cfg := SortConfig{Enabled: true, Field: "key", Direction: "asc"}
	keys := ApplySort(map[string]string{}, cfg)
	if len(keys) != 0 {
		t.Errorf("expected empty slice, got %v", keys)
	}
}
