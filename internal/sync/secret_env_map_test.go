package sync

import (
	"testing"
)

func TestApplyEnvMap_Disabled_ReturnsOriginal(t *testing.T) {
	cfg := EnvMapConfig{Enabled: false, Mappings: map[string]string{"A": "B"}}
	secrets := map[string]string{"A": "value"}

	result := ApplyEnvMap(cfg, secrets)
	if result["A"] != "value" {
		t.Errorf("expected original key A, got %v", result)
	}
}

func TestApplyEnvMap_NoMappings_ReturnsOriginal(t *testing.T) {
	cfg := EnvMapConfig{Enabled: true, Mappings: map[string]string{}}
	secrets := map[string]string{"A": "value"}

	result := ApplyEnvMap(cfg, secrets)
	if result["A"] != "value" {
		t.Errorf("expected original key A, got %v", result)
	}
}

func TestApplyEnvMap_RenamesSingleKey(t *testing.T) {
	cfg := EnvMapConfig{
		Enabled:  true,
		Mappings: map[string]string{"DB_PASS": "DATABASE_PASSWORD"},
	}
	secrets := map[string]string{"DB_PASS": "secret123", "OTHER": "val"}

	result := ApplyEnvMap(cfg, secrets)
	if _, ok := result["DB_PASS"]; ok {
		t.Error("old key DB_PASS should not exist")
	}
	if result["DATABASE_PASSWORD"] != "secret123" {
		t.Errorf("expected DATABASE_PASSWORD=secret123, got %v", result)
	}
	if result["OTHER"] != "val" {
		t.Error("unmapped key OTHER should be preserved")
	}
}

func TestApplyEnvMap_RenamesMultipleKeys(t *testing.T) {
	cfg := EnvMapConfig{
		Enabled: true,
		Mappings: map[string]string{
			"A": "X",
			"B": "Y",
		},
	}
	secrets := map[string]string{"A": "1", "B": "2", "C": "3"}

	result := ApplyEnvMap(cfg, secrets)
	if result["X"] != "1" || result["Y"] != "2" || result["C"] != "3" {
		t.Errorf("unexpected result: %v", result)
	}
}

func TestApplyEnvMap_DoesNotMutateInput(t *testing.T) {
	cfg := EnvMapConfig{
		Enabled:  true,
		Mappings: map[string]string{"OLD": "NEW"},
	}
	secrets := map[string]string{"OLD": "val"}

	ApplyEnvMap(cfg, secrets)
	if _, ok := secrets["OLD"]; !ok {
		t.Error("original map should not be mutated")
	}
}
