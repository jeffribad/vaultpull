package sync

import (
	"testing"
)

func TestSanitizeSecrets_Disabled_ReturnsOriginal(t *testing.T) {
	input := map[string]string{"KEY": "  value  "}
	cfg := SanitizeConfig{Enabled: false}
	out := SanitizeSecrets(input, cfg)
	if out["KEY"] != "  value  " {
		t.Errorf("expected original value, got %q", out["KEY"])
	}
}

func TestSanitizeSecrets_TrimWhitespace(t *testing.T) {
	input := map[string]string{"KEY": "  hello  "}
	cfg := SanitizeConfig{Enabled: true, TrimWhitespace: true}
	out := SanitizeSecrets(input, cfg)
	if out["KEY"] != "hello" {
		t.Errorf("expected trimmed value, got %q", out["KEY"])
	}
}

func TestSanitizeSecrets_ReplaceNewlines(t *testing.T) {
	input := map[string]string{"KEY": "line1\nline2"}
	cfg := SanitizeConfig{Enabled: true, ReplaceNewlines: true, NewlineReplacement: " "}
	out := SanitizeSecrets(input, cfg)
	if out["KEY"] != "line1 line2" {
		t.Errorf("expected newline replaced, got %q", out["KEY"])
	}
}

func TestSanitizeSecrets_StripControlChars(t *testing.T) {
	input := map[string]string{"KEY": "val\x01ue\x1f"}
	cfg := SanitizeConfig{Enabled: true, StripControl: true}
	out := SanitizeSecrets(input, cfg)
	if out["KEY"] != "value" {
		t.Errorf("expected control chars stripped, got %q", out["KEY"])
	}
}

func TestSanitizeSecrets_DoesNotMutateInput(t *testing.T) {
	input := map[string]string{"KEY": "  val\x00ue  "}
	cfg := SanitizeConfig{Enabled: true, TrimWhitespace: true, StripControl: true}
	_ = SanitizeSecrets(input, cfg)
	if input["KEY"] != "  val\x00ue  " {
		t.Error("input map was mutated")
	}
}

func TestSanitizeSecrets_CRLFReplaced(t *testing.T) {
	input := map[string]string{"KEY": "a\r\nb"}
	cfg := SanitizeConfig{Enabled: true, ReplaceNewlines: true, NewlineReplacement: "|"}
	out := SanitizeSecrets(input, cfg)
	if out["KEY"] != "a|b" {
		t.Errorf("expected CRLF replaced, got %q", out["KEY"])
	}
}

func TestStripControlChars_PreservesTab(t *testing.T) {
	result := stripControlChars("col1\tcol2")
	if result != "col1\tcol2" {
		t.Errorf("expected tab preserved, got %q", result)
	}
}
