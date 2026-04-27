package sync

import (
	"testing"
)

func TestConditional_Integration_ConfigDriven(t *testing.T) {
	t.Setenv("VAULTPULL_CONDITIONAL_ENABLED", "true")
	t.Setenv("VAULTPULL_CONDITIONAL_RULES", "PROD_DB:DEPLOY_ENV=production")
	t.Setenv("DEPLOY_ENV", "production")

	cfg := ConditionalConfigFromEnv()
	secrets := map[string]string{
		"PROD_DB": "prod-connection-string",
		"SHARED":  "shared-value",
	}

	result := ApplyConditional(cfg, secrets)

	if _, ok := result["PROD_DB"]; !ok {
		t.Error("expected PROD_DB to be present when DEPLOY_ENV=production")
	}
	if _, ok := result["SHARED"]; !ok {
		t.Error("expected SHARED to always be present")
	}
}

func TestConditional_Integration_DisabledPassthrough(t *testing.T) {
	t.Setenv("VAULTPULL_CONDITIONAL_ENABLED", "false")
	t.Setenv("VAULTPULL_CONDITIONAL_RULES", "PROD_DB:DEPLOY_ENV=production")
	t.Setenv("DEPLOY_ENV", "staging")

	cfg := ConditionalConfigFromEnv()
	secrets := map[string]string{
		"PROD_DB": "prod-connection-string",
		"SHARED":  "shared-value",
	}

	result := ApplyConditional(cfg, secrets)

	if len(result) != len(secrets) {
		t.Errorf("expected all secrets to pass through when disabled, got %d/%d", len(result), len(secrets))
	}
}

func TestConditional_Integration_MultiRuleExclusion(t *testing.T) {
	t.Setenv("VAULTPULL_CONDITIONAL_ENABLED", "true")
	t.Setenv("VAULTPULL_CONDITIONAL_RULES", "PROD_DB:DEPLOY_ENV=production,PREMIUM_KEY:TIER=premium")
	t.Setenv("DEPLOY_ENV", "staging")
	t.Setenv("TIER", "free")

	cfg := ConditionalConfigFromEnv()
	secrets := map[string]string{
		"PROD_DB":     "conn",
		"PREMIUM_KEY": "key",
		"COMMON":      "common",
	}

	result := ApplyConditional(cfg, secrets)

	if _, ok := result["PROD_DB"]; ok {
		t.Error("expected PROD_DB to be excluded")
	}
	if _, ok := result["PREMIUM_KEY"]; ok {
		t.Error("expected PREMIUM_KEY to be excluded")
	}
	if _, ok := result["COMMON"]; !ok {
		t.Error("expected COMMON to be present")
	}
}
