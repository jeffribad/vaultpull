package sync

import (
	"testing"
)

func TestApplyConditional_Disabled_ReturnsOriginal(t *testing.T) {
	cfg := ConditionalConfig{Enabled: false}
	secrets := map[string]string{"DB_HOST": "localhost", "API_KEY": "abc"}
	result := ApplyConditional(cfg, secrets)
	if len(result) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(result))
	}
}

func TestApplyConditional_NoRules_ReturnsAll(t *testing.T) {
	cfg := ConditionalConfig{Enabled: true, Rules: map[string]conditionalRule{}}
	secrets := map[string]string{"DB_HOST": "localhost"}
	result := ApplyConditional(cfg, secrets)
	if len(result) != 1 {
		t.Fatalf("expected 1 key, got %d", len(result))
	}
}

func TestApplyConditional_MatchingRule_IncludesKey(t *testing.T) {
	t.Setenv("APP_ENV", "production")
	cfg := ConditionalConfig{
		Enabled: true,
		Rules: map[string]conditionalRule{
			"PROD_SECRET": {EnvVar: "APP_ENV", Expected: "production"},
		},
	}
	secrets := map[string]string{"PROD_SECRET": "s3cr3t", "OTHER": "val"}
	result := ApplyConditional(cfg, secrets)
	if _, ok := result["PROD_SECRET"]; !ok {
		t.Error("expected PROD_SECRET to be included")
	}
	if _, ok := result["OTHER"]; !ok {
		t.Error("expected OTHER to be included")
	}
}

func TestApplyConditional_NonMatchingRule_ExcludesKey(t *testing.T) {
	t.Setenv("APP_ENV", "staging")
	cfg := ConditionalConfig{
		Enabled: true,
		Rules: map[string]conditionalRule{
			"PROD_SECRET": {EnvVar: "APP_ENV", Expected: "production"},
		},
	}
	secrets := map[string]string{"PROD_SECRET": "s3cr3t", "OTHER": "val"}
	result := ApplyConditional(cfg, secrets)
	if _, ok := result["PROD_SECRET"]; ok {
		t.Error("expected PROD_SECRET to be excluded")
	}
	if _, ok := result["OTHER"]; !ok {
		t.Error("expected OTHER to be included")
	}
}

func TestApplyConditional_CaseInsensitiveKeyMatch(t *testing.T) {
	t.Setenv("FEATURE_FLAG", "on")
	cfg := ConditionalConfig{
		Enabled: true,
		Rules: map[string]conditionalRule{
			"FEATURE_KEY": {EnvVar: "FEATURE_FLAG", Expected: "on"},
		},
	}
	secrets := map[string]string{"feature_key": "value"}
	result := ApplyConditional(cfg, secrets)
	if _, ok := result["feature_key"]; !ok {
		t.Error("expected feature_key to be included via case-insensitive rule match")
	}
}

func TestApplyConditional_DoesNotMutateInput(t *testing.T) {
	t.Setenv("ENV", "dev")
	cfg := ConditionalConfig{
		Enabled: true,
		Rules: map[string]conditionalRule{
			"PROD_ONLY": {EnvVar: "ENV", Expected: "prod"},
		},
	}
	secrets := map[string]string{"PROD_ONLY": "secret", "SHARED": "val"}
	origLen := len(secrets)
	ApplyConditional(cfg, secrets)
	if len(secrets) != origLen {
		t.Error("input map was mutated")
	}
}
