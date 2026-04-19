package sync

import (
	"testing"
)

func TestCheckDeprecated_Disabled_ReturnsNil(t *testing.T) {
	cfg := DeprecateConfig{Enabled: false, Deprecated: map[string]string{"OLD_KEY": "NEW_KEY"}}
	secrets := map[string]string{"OLD_KEY": "value"}
	violations, err := CheckDeprecated(cfg, secrets)
	if err != nil || len(violations) != 0 {
		t.Fatalf("expected no violations when disabled, got %v, %v", violations, err)
	}
}

func TestCheckDeprecated_NoDeprecatedKeys_ReturnsNil(t *testing.T) {
	cfg := DeprecateConfig{Enabled: true, Deprecated: map[string]string{}}
	secrets := map[string]string{"SOME_KEY": "value"}
	violations, err := CheckDeprecated(cfg, secrets)
	if err != nil || len(violations) != 0 {
		t.Fatalf("expected no violations, got %v", violations)
	}
}

func TestCheckDeprecated_DetectsViolation(t *testing.T) {
	cfg := DeprecateConfig{
		Enabled:     true,
		Deprecated:  map[string]string{"OLD_KEY": "NEW_KEY"},
		FailOnUsage: false,
	}
	secrets := map[string]string{"OLD_KEY": "val", "OTHER": "x"}
	violations, err := CheckDeprecated(cfg, secrets)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
	if violations[0].Key != "OLD_KEY" || violations[0].Replacement != "NEW_KEY" {
		t.Errorf("unexpected violation: %+v", violations[0])
	}
}

func TestCheckDeprecated_FailOnUsage_ReturnsError(t *testing.T) {
	cfg := DeprecateConfig{
		Enabled:     true,
		Deprecated:  map[string]string{"OLD_KEY": "NEW_KEY"},
		FailOnUsage: true,
	}
	secrets := map[string]string{"OLD_KEY": "val"}
	_, err := CheckDeprecated(cfg, secrets)
	if err == nil {
		t.Fatal("expected error when FailOnUsage is true")
	}
}

func TestCheckDeprecated_CaseInsensitiveMatch(t *testing.T) {
	cfg := DeprecateConfig{
		Enabled:    true,
		Deprecated: map[string]string{"OLD_KEY": ""},
	}
	secrets := map[string]string{"old_key": "value"}
	violations, err := CheckDeprecated(cfg, secrets)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation for case-insensitive match, got %d", len(violations))
	}
}

func TestDeprecationViolation_Error_WithReplacement(t *testing.T) {
	v := DeprecationViolation{Key: "OLD", Replacement: "NEW"}
	if msg := v.Error(); msg == "" {
		t.Error("expected non-empty error message")
	}
}

func TestDeprecationViolation_Error_NoReplacement(t *testing.T) {
	v := DeprecationViolation{Key: "OLD", Replacement: ""}
	if msg := v.Error(); msg == "" {
		t.Error("expected non-empty error message")
	}
}
