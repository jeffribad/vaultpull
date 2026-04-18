package sync

import (
	"regexp"
	"testing"
)

func makeSchemaRule(key, pattern string) SchemaRule {
	return SchemaRule{Key: key, Pattern: pattern, Regexp: regexp.MustCompile(pattern)}
}

func TestValidateSchema_Disabled_ReturnsNil(t *testing.T) {
	cfg := SchemaConfig{Enabled: false, Rules: []SchemaRule{makeSchemaRule("API_KEY", `^\w{32}$`)}}
	secrets := map[string]string{"API_KEY": "bad"}
	if errs := ValidateSchema(cfg, secrets); len(errs) != 0 {
		t.Errorf("expected no errors when disabled, got %d", len(errs))
	}
}

func TestValidateSchema_NoRules_ReturnsNil(t *testing.T) {
	cfg := SchemaConfig{Enabled: true, Rules: nil}
	secrets := map[string]string{"FOO": "bar"}
	if errs := ValidateSchema(cfg, secrets); len(errs) != 0 {
		t.Errorf("expected no errors with no rules, got %d", len(errs))
	}
}

func TestValidateSchema_PassingValue(t *testing.T) {
	cfg := SchemaConfig{
		Enabled: true,
		Rules:   []SchemaRule{makeSchemaRule("API_KEY", `^[A-Za-z0-9]{8}$`)},
	}
	secrets := map[string]string{"API_KEY": "Ab3dEf7H"}
	if errs := ValidateSchema(cfg, secrets); len(errs) != 0 {
		t.Errorf("expected no errors, got: %v", errs)
	}
}

func TestValidateSchema_FailingValue(t *testing.T) {
	cfg := SchemaConfig{
		Enabled: true,
		Rules:   []SchemaRule{makeSchemaRule("API_KEY", `^[A-Za-z0-9]{8}$`)},
	}
	secrets := map[string]string{"API_KEY": "tooshort"}
	// "tooshort" is 8 chars but all lowercase — actually matches; use a stricter pattern
	cfg.Rules[0] = makeSchemaRule("API_KEY", `^[0-9]{8}$`)
	errs := ValidateSchema(cfg, secrets)
	if len(errs) != 1 {
		t.Errorf("expected 1 error, got %d", len(errs))
	}
}

func TestValidateSchema_MissingKeySkipped(t *testing.T) {
	cfg := SchemaConfig{
		Enabled: true,
		Rules:   []SchemaRule{makeSchemaRule("REQUIRED_KEY", `^.+$`)},
	}
	secrets := map[string]string{"OTHER_KEY": "value"}
	if errs := ValidateSchema(cfg, secrets); len(errs) != 0 {
		t.Errorf("expected missing key to be skipped, got %d errors", len(errs))
	}
}

func TestValidateSchema_MultipleFailures(t *testing.T) {
	cfg := SchemaConfig{
		Enabled: true,
		Rules: []SchemaRule{
			makeSchemaRule("KEY1", `^[0-9]+$`),
			makeSchemaRule("KEY2", `^[0-9]+$`),
		},
	}
	secrets := map[string]string{"KEY1": "abc", "KEY2": "xyz"}
	if errs := ValidateSchema(cfg, secrets); len(errs) != 2 {
		t.Errorf("expected 2 errors, got %d", len(errs))
	}
}
