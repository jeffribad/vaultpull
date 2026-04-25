package sync

import (
	"testing"
)

func TestApplyEnvCopy_Disabled_ReturnsOriginal(t *testing.T) {
	cfg := CopyConfig{Enabled: false, Pairs: map[string]string{"A": "B"}}
	secrets := map[string]string{"A": "hello"}
	out := ApplyEnvCopy(cfg, secrets)
	if _, ok := out["B"]; ok {
		t.Error("expected B to not exist when disabled")
	}
}

func TestApplyEnvCopy_NoPairs_ReturnsOriginal(t *testing.T) {
	cfg := CopyConfig{Enabled: true, Pairs: map[string]string{}}
	secrets := map[string]string{"A": "hello"}
	out := ApplyEnvCopy(cfg, secrets)
	if len(out) != 1 {
		t.Errorf("expected 1 key, got %d", len(out))
	}
}

func TestApplyEnvCopy_CopiesSrcToDst(t *testing.T) {
	cfg := CopyConfig{Enabled: true, Pairs: map[string]string{"DB_PASS": "DATABASE_PASSWORD"}}
	secrets := map[string]string{"DB_PASS": "secret123"}
	out := ApplyEnvCopy(cfg, secrets)
	if out["DATABASE_PASSWORD"] != "secret123" {
		t.Errorf("expected DATABASE_PASSWORD=secret123, got %q", out["DATABASE_PASSWORD"])
	}
	if out["DB_PASS"] != "secret123" {
		t.Error("expected source key DB_PASS to be preserved")
	}
}

func TestApplyEnvCopy_CaseInsensitiveSrcLookup(t *testing.T) {
	cfg := CopyConfig{Enabled: true, Pairs: map[string]string{"db_pass": "DB_PASSWORD"}}
	secrets := map[string]string{"DB_PASS": "topsecret"}
	out := ApplyEnvCopy(cfg, secrets)
	if out["DB_PASSWORD"] != "topsecret" {
		t.Errorf("expected DB_PASSWORD=topsecret, got %q", out["DB_PASSWORD"])
	}
}

func TestApplyEnvCopy_DoesNotMutateInput(t *testing.T) {
	cfg := CopyConfig{Enabled: true, Pairs: map[string]string{"X": "Y"}}
	secrets := map[string]string{"X": "val"}
	ApplyEnvCopy(cfg, secrets)
	if _, ok := secrets["Y"]; ok {
		t.Error("expected input map not to be mutated")
	}
}

func TestParseCopyPairs_SinglePair(t *testing.T) {
	pairs := parseCopyPairs("FOO:BAR")
	if pairs["FOO"] != "BAR" {
		t.Errorf("expected FOO->BAR, got %v", pairs)
	}
}

func TestParseCopyPairs_MultiplePairs(t *testing.T) {
	pairs := parseCopyPairs("A:B, C:D")
	if pairs["A"] != "B" || pairs["C"] != "D" {
		t.Errorf("unexpected pairs: %v", pairs)
	}
}

func TestParseCopyPairs_SkipsMalformed(t *testing.T) {
	pairs := parseCopyPairs("NOCOLON, :NODST, SRC:")
	if len(pairs) != 0 {
		t.Errorf("expected 0 valid pairs, got %v", pairs)
	}
}
