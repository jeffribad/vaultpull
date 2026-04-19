package sync

import (
	"testing"
)

func TestApplyFormat_Disabled_ReturnsOriginal(t *testing.T) {
	input := map[string]string{"KEY": "  value  "}
	cfg := FormatConfig{Enabled: false, TrimSpace: true}
	out := ApplyFormat(input, cfg)
	if out["KEY"] != "  value  " {
		t.Errorf("expected original value, got %q", out["KEY"])
	}
}

func TestApplyFormat_TrimSpace(t *testing.T) {
	input := map[string]string{"KEY": "  hello world  "}
	cfg := FormatConfig{Enabled: true, TrimSpace: true}
	out := ApplyFormat(input, cfg)
	if out["KEY"] != "hello world" {
		t.Errorf("expected trimmed value, got %q", out["KEY"])
	}
}

func TestApplyFormat_NormalizeNewlines(t *testing.T) {
	input := map[string]string{"KEY": "line1\r\nline2\rline3"}
	cfg := FormatConfig{Enabled: true, NormalizeNewlines: true}
	out := ApplyFormat(input, cfg)
	expected := "line1\nline2\nline3"
	if out["KEY"] != expected {
		t.Errorf("expected %q, got %q", expected, out["KEY"])
	}
}

func TestApplyFormat_StripNulls(t *testing.T) {
	input := map[string]string{"KEY": "val\x00ue"}
	cfg := FormatConfig{Enabled: true, StripNulls: true}
	out := ApplyFormat(input, cfg)
	if out["KEY"] != "value" {
		t.Errorf("expected null-stripped value, got %q", out["KEY"])
	}
}

func TestApplyFormat_DoesNotMutateInput(t *testing.T) {
	input := map[string]string{"KEY": "  spaced  "}
	cfg := FormatConfig{Enabled: true, TrimSpace: true}
	_ = ApplyFormat(input, cfg)
	if input["KEY"] != "  spaced  " {
		t.Error("input map was mutated")
	}
}

func TestFormatConfigFromEnv_Defaults(t *testing.T) {
	cfg := FormatConfigFromEnv()
	if cfg.Enabled {
		t.Error("expected Enabled to be false by default")
	}
	if cfg.TrimSpace {
		t.Error("expected TrimSpace to be false by default")
	}
}

func TestIsTruthy(t *testing.T) {
	cases := map[string]bool{
		"true": true, "1": true, "yes": true,
		"false": false, "0": false, "": false, "no": false,
	}
	for input, want := range cases {
		got := isTruthy(input)
		if got != want {
			t.Errorf("isTruthy(%q) = %v, want %v", input, got, want)
		}
	}
}
