package sync

import (
	"testing"
)

func TestValidateSchema_Integration_ConfigDriven(t *testing.T) {
	t.Setenv("VAULTPULL_SCHEMA_ENABLED", "1")
	t.Setenv("VAULTPULL_SCHEMA_RULES", "DB_PASSWORD:^.{12,}$,API_TOKEN:^tok_[a-z]+$")

	cfg := SchemaConfigFromEnv()

	secrets := map[string]string{
		"DB_PASSWORD": "supersecretpass",
		"API_TOKEN":   "tok_abcdef",
	}
	if errs := ValidateSchema(cfg, secrets); len(errs) != 0 {
		t.Errorf("expected no errors, got: %v", errs)
	}
}

func TestValidateSchema_Integration_DetectsViolations(t *testing.T) {
	t.Setenv("VAULTPULL_SCHEMA_ENABLED", "1")
	t.Setenv("VAULTPULL_SCHEMA_RULES", "DB_PASSWORD:^.{12,}$,API_TOKEN:^tok_[a-z]+$")

	cfg := SchemaConfigFromEnv()

	secrets := map[string]string{
		"DB_PASSWORD": "short",
		"API_TOKEN":   "invalid_token",
	}
	errs := ValidateSchema(cfg, secrets)
	if len(errs) != 2 {
		t.Errorf("expected 2 violations, got %d: %v", len(errs), errs)
	}
}

func TestValidateSchema_Integration_DisabledSkipsAll(t *testing.T) {
	t.Setenv("VAULTPULL_SCHEMA_ENABLED", "0")
	t.Setenv("VAULTPULL_SCHEMA_RULES", "API_TOKEN:^tok_[a-z]+$")

	cfg := SchemaConfigFromEnv()

	secrets := map[string]string{"API_TOKEN": "completely_wrong"}
	if errs := ValidateSchema(cfg, secrets); len(errs) != 0 {
		t.Errorf("expected disabled schema to skip validation, got %d errors", len(errs))
	}
}
