package sync

import (
	"testing"
)

func TestConditionalConfigFromEnv_Defaults(t *testing.T) {
	t.Setenv("VAULTPULL_CONDITIONAL_ENABLED", "")
	t.Setenv("VAULTPULL_CONDITIONAL_RULES", "")
	cfg := ConditionalConfigFromEnv()
	if cfg.Enabled {
		t.Error("expected Enabled=false by default")
	}
	if len(cfg.Rules) != 0 {
		t.Errorf("expected 0 rules, got %d", len(cfg.Rules))
	}
}

func TestConditionalConfigFromEnv_Enabled(t *testing.T) {
	t.Setenv("VAULTPULL_CONDITIONAL_ENABLED", "true")
	cfg := ConditionalConfigFromEnv()
	if !cfg.Enabled {
		t.Error("expected Enabled=true")
	}
}

func TestConditionalConfigFromEnv_NumericEnabled(t *testing.T) {
	t.Setenv("VAULTPULL_CONDITIONAL_ENABLED", "1")
	cfg := ConditionalConfigFromEnv()
	if !cfg.Enabled {
		t.Error("expected Enabled=true for value '1'")
	}
}

func TestConditionalConfigFromEnv_ParsesRules(t *testing.T) {
	t.Setenv("VAULTPULL_CONDITIONAL_ENABLED", "true")
	t.Setenv("VAULTPULL_CONDITIONAL_RULES", "DB_PASS:APP_ENV=production, API_KEY : TIER = premium")
	cfg := ConditionalConfigFromEnv()
	if len(cfg.Rules) != 2 {
		t.Fatalf("expected 2 rules, got %d", len(cfg.Rules))
	}
	r, ok := cfg.Rules["DB_PASS"]
	if !ok {
		t.Fatal("expected rule for DB_PASS")
	}
	if r.EnvVar != "APP_ENV" || r.Expected != "production" {
		t.Errorf("unexpected rule: %+v", r)
	}
	r2, ok := cfg.Rules["API_KEY"]
	if !ok {
		t.Fatal("expected rule for API_KEY")
	}
	if r2.EnvVar != "TIER" || r2.Expected != "premium" {
		t.Errorf("unexpected rule: %+v", r2)
	}
}

func TestConditionalConfigFromEnv_SkipsMalformedRules(t *testing.T) {
	t.Setenv("VAULTPULL_CONDITIONAL_ENABLED", "true")
	t.Setenv("VAULTPULL_CONDITIONAL_RULES", "NOCORON,NOEQ:MISSINGEQ,VALID:ENV=val")
	cfg := ConditionalConfigFromEnv()
	if len(cfg.Rules) != 1 {
		t.Errorf("expected 1 valid rule, got %d", len(cfg.Rules))
	}
}
