package sync

import (
	"testing"
)

func TestSchemaConfigFromEnv_Defaults(t *testing.T) {
	t.Setenv("VAULTPULL_SCHEMA_ENABLED", "")
	t.Setenv("VAULTPULL_SCHEMA_RULES", "")

	cfg := SchemaConfigFromEnv()
	if cfg.Enabled {
		t.Error("expected Enabled=false by default")
	}
	if len(cfg.Rules) != 0 {
		t.Errorf("expected no rules, got %d", len(cfg.Rules))
	}
}

func TestSchemaConfigFromEnv_Enabled(t *testing.T) {
	t.Setenv("VAULTPULL_SCHEMA_ENABLED", "1")
	t.Setenv("VAULTPULL_SCHEMA_RULES", "")

	cfg := SchemaConfigFromEnv()
	if !cfg.Enabled {
		t.Error("expected Enabled=true")
	}
}

func TestSchemaConfigFromEnv_ParsesRules(t *testing.T) {
	t.Setenv("VAULTPULL_SCHEMA_ENABLED", "1")
	t.Setenv("VAULTPULL_SCHEMA_RULES", "API_KEY:^[A-Za-z0-9]{32}$,DB_URL:^postgres://")

	cfg := SchemaConfigFromEnv()
	if len(cfg.Rules) != 2 {
		t.Fatalf("expected 2 rules, got %d", len(cfg.Rules))
	}
	if cfg.Rules[0].Key != "API_KEY" {
		t.Errorf("unexpected key: %s", cfg.Rules[0].Key)
	}
	if cfg.Rules[1].Key != "DB_URL" {
		t.Errorf("unexpected key: %s", cfg.Rules[1].Key)
	}
}

func TestSchemaConfigFromEnv_SkipsMalformedRules(t *testing.T) {
	t.Setenv("VAULTPULL_SCHEMA_ENABLED", "1")
	t.Setenv("VAULTPULL_SCHEMA_RULES", "NOPATTERN,VALID:^ok$")

	cfg := SchemaConfigFromEnv()
	if len(cfg.Rules) != 1 {
		t.Errorf("expected 1 valid rule, got %d", len(cfg.Rules))
	}
}
