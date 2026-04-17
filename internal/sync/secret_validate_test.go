package sync

import (
	"testing"
)

func TestValidateSecrets_Disabled_AlwaysPasses(t *testing.T) {
	cfg := ValidateConfig{Enabled: false, RequiredKeys: []string{"MISSING_KEY"}}
	if err := ValidateSecrets(map[string]string{}, cfg); err != nil {
		t.Errorf("expected no error when disabled, got: %v", err)
	}
}

func TestValidateSecrets_RequiredKey_Present(t *testing.T) {
	cfg := ValidateConfig{Enabled: true, RequiredKeys: []string{"DB_HOST"}}
	err := ValidateSecrets(map[string]string{"DB_HOST": "localhost"}, cfg)
	if err != nil {
		t.Errorf("expected no error, got: %v", err)
	}
}

func TestValidateSecrets_RequiredKey_Missing(t *testing.T) {
	cfg := ValidateConfig{Enabled: true, RequiredKeys: []string{"DB_HOST"}}
	err := ValidateSecrets(map[string]string{}, cfg)
	if err == nil {
		t.Fatal("expected error for missing required key")
	}
	ve, ok := err.(*ValidationError)
	if !ok || !ve.HasErrors() {
		t.Errorf("expected ValidationError, got: %v", err)
	}
}

func TestValidateSecrets_NonemptyKey_EmptyValue(t *testing.T) {
	cfg := ValidateConfig{Enabled: true, NonemptyKeys: []string{"API_KEY"}}
	err := ValidateSecrets(map[string]string{"API_KEY": "   "}, cfg)
	if err == nil {
		t.Fatal("expected error for empty nonempty key")
	}
}

func TestValidateSecrets_NonemptyKey_ValidValue(t *testing.T) {
	cfg := ValidateConfig{Enabled: true, NonemptyKeys: []string{"API_KEY"}}
	err := ValidateSecrets(map[string]string{"API_KEY": "abc123"}, cfg)
	if err != nil {
		t.Errorf("expected no error, got: %v", err)
	}
}

func TestValidateSecrets_MultipleErrors(t *testing.T) {
	cfg := ValidateConfig{
		Enabled:      true,
		RequiredKeys: []string{"DB_HOST", "DB_PORT"},
		NonemptyKeys: []string{"TOKEN"},
	}
	err := ValidateSecrets(map[string]string{"TOKEN": ""}, cfg)
	ve, ok := err.(*ValidationError)
	if !ok {
		t.Fatalf("expected *ValidationError, got %T", err)
	}
	if len(ve.Errors) != 3 {
		t.Errorf("expected 3 errors, got %d: %v", len(ve.Errors), ve.Errors)
	}
}
