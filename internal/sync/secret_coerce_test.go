package sync

import (
	"testing"
)

func TestCoerceSecrets_Disabled_ReturnsOriginal(t *testing.T) {
	input := map[string]string{"ACTIVE": "1", "COUNT": "  42  "}
	cfg := CoerceConfig{Enabled: false, BoolKeys: []string{"ACTIVE"}}
	out := CoerceSecrets(input, cfg)
	if out["ACTIVE"] != "1" {
		t.Errorf("expected original value, got %q", out["ACTIVE"])
	}
}

func TestCoerceSecrets_BoolTrue(t *testing.T) {
	for _, val := range []string{"1", "true", "yes", "on", "TRUE", "YES"} {
		out := CoerceSecrets(
			map[string]string{"FLAG": val},
			CoerceConfig{Enabled: true, BoolKeys: []string{"FLAG"}},
		)
		if out["FLAG"] != "true" {
			t.Errorf("expected true for input %q, got %q", val, out["FLAG"])
		}
	}
}

func TestCoerceSecrets_BoolFalse(t *testing.T) {
	for _, val := range []string{"0", "false", "no", "off", ""} {
		out := CoerceSecrets(
			map[string]string{"FLAG": val},
			CoerceConfig{Enabled: true, BoolKeys: []string{"FLAG"}},
		)
		if out["FLAG"] != "false" {
			t.Errorf("expected false for input %q, got %q", val, out["FLAG"])
		}
	}
}

func TestCoerceSecrets_NumberTrimsSpace(t *testing.T) {
	out := CoerceSecrets(
		map[string]string{"PORT": "  8080  "},
		CoerceConfig{Enabled: true, NumberKeys: []string{"PORT"}},
	)
	if out["PORT"] != "8080" {
		t.Errorf("expected trimmed number, got %q", out["PORT"])
	}
}

func TestCoerceSecrets_JSONCompactsWhitespace(t *testing.T) {
	out := CoerceSecrets(
		map[string]string{"CONFIG": "{  \"a\":   1  }"},
		CoerceConfig{Enabled: true, JSONKeys: []string{"CONFIG"}},
	)
	if out["CONFIG"] != `{"a": 1}` {
		t.Errorf("unexpected JSON coercion result: %q", out["CONFIG"])
	}
}

func TestCoerceSecrets_DoesNotMutateInput(t *testing.T) {
	input := map[string]string{"FLAG": "yes"}
	CoerceSecrets(input, CoerceConfig{Enabled: true, BoolKeys: []string{"FLAG"}})
	if input["FLAG"] != "yes" {
		t.Error("input map was mutated")
	}
}

func TestCoerceSecrets_MissingKey_Skipped(t *testing.T) {
	input := map[string]string{"OTHER": "value"}
	out := CoerceSecrets(input, CoerceConfig{Enabled: true, BoolKeys: []string{"FLAG"}})
	if _, ok := out["FLAG"]; ok {
		t.Error("expected missing key to be skipped")
	}
	if out["OTHER"] != "value" {
		t.Error("unrelated key should be preserved")
	}
}
