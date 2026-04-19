package sync

import (
	"testing"
)

func TestNormalizeSecrets_Disabled_ReturnsOriginal(t *testing.T) {
	cfg := NormalizeConfig{Enabled: false}
	input := map[string]string{"myKey": "val"}
	out, err := NormalizeSecrets(cfg, input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["myKey"] != "val" {
		t.Errorf("expected original key preserved")
	}
}

func TestNormalizeSecrets_UpperKeys(t *testing.T) {
	cfg := NormalizeConfig{Enabled: true, UpperKeys: true}
	input := map[string]string{"db_host": "localhost"}
	out, err := NormalizeSecrets(cfg, input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["DB_HOST"] != "localhost" {
		t.Errorf("expected DB_HOST, got keys: %v", out)
	}
}

func TestNormalizeSecrets_SnakeCase(t *testing.T) {
	cfg := NormalizeConfig{Enabled: true, SnakeCase: true}
	input := map[string]string{"DbHost": "localhost"}
	out, err := NormalizeSecrets(cfg, input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["db_host"] != "localhost" {
		t.Errorf("expected db_host, got: %v", out)
	}
}

func TestNormalizeSecrets_StripDashes(t *testing.T) {
	cfg := NormalizeConfig{Enabled: true, StripDashes: true}
	input := map[string]string{"my-key": "value"}
	out, err := NormalizeSecrets(cfg, input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["my_key"] != "value" {
		t.Errorf("expected my_key, got: %v", out)
	}
}

func TestNormalizeSecrets_CombinedRules(t *testing.T) {
	cfg := NormalizeConfig{Enabled: true, UpperKeys: true, StripDashes: true}
	input := map[string]string{"my-secret-key": "abc"}
	out, err := NormalizeSecrets(cfg, input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["MY_SECRET_KEY"] != "abc" {
		t.Errorf("expected MY_SECRET_KEY, got: %v", out)
	}
}

func TestNormalizeConfigFromEnv_Defaults(t *testing.T) {
	cfg := NormalizeConfigFromEnv()
	if cfg.Enabled {
		t.Errorf("expected Enabled=false by default")
	}
	if cfg.UpperKeys || cfg.SnakeCase || cfg.StripDashes {
		t.Errorf("expected all options false by default")
	}
}
