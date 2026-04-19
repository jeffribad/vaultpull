package sync

import (
	"testing"
)

func TestCheckPins_Disabled_ReturnsNil(t *testing.T) {
	cfg := PinConfig{Enabled: false, Pins: map[string]string{"DB_PASS": "abc"}}
	errs := CheckPins(cfg, map[string]string{"DB_PASS": "xyz_value"})
	if errs != nil {
		t.Fatalf("expected nil, got %v", errs)
	}
}

func TestCheckPins_NoPins_ReturnsNil(t *testing.T) {
	cfg := PinConfig{Enabled: true, Pins: map[string]string{}}
	errs := CheckPins(cfg, map[string]string{"DB_PASS": "abc123"})
	if errs != nil {
		t.Fatalf("expected nil, got %v", errs)
	}
}

func TestCheckPins_MatchingPrefix_ReturnsNil(t *testing.T) {
	cfg := PinConfig{Enabled: true, Pins: map[string]string{"API_KEY": "sk-prod"}}
	errs := CheckPins(cfg, map[string]string{"API_KEY": "sk-prod-abc123"})
	if errs != nil {
		t.Fatalf("expected nil, got %v", errs)
	}
}

func TestCheckPins_ViolationDetected(t *testing.T) {
	cfg := PinConfig{Enabled: true, Pins: map[string]string{"API_KEY": "sk-prod"}}
	errs := CheckPins(cfg, map[string]string{"API_KEY": "sk-test-xyz"})
	if len(errs) != 1 {
		t.Fatalf("expected 1 error, got %d", len(errs))
	}
}

func TestCheckPins_FailFast_StopsEarly(t *testing.T) {
	cfg := PinConfig{
		Enabled:  true,
		FailFast: true,
		Pins:     map[string]string{"KEY_A": "aaa", "KEY_B": "bbb"},
	}
	secrets := map[string]string{"KEY_A": "xxx", "KEY_B": "yyy"}
	errs := CheckPins(cfg, secrets)
	if len(errs) != 1 {
		t.Fatalf("expected 1 error with fail-fast, got %d", len(errs))
	}
}

func TestCheckPins_MissingKey_Skipped(t *testing.T) {
	cfg := PinConfig{Enabled: true, Pins: map[string]string{"MISSING": "abc"}}
	errs := CheckPins(cfg, map[string]string{})
	if errs != nil {
		t.Fatalf("expected nil for missing key, got %v", errs)
	}
}

func TestPinViolation_Error(t *testing.T) {
	v := PinViolation{Key: "X", Expected: "abc", Actual: "xyz"}
	if v.Error() == "" {
		t.Fatal("expected non-empty error string")
	}
}
