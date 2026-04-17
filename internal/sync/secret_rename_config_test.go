package sync

import (
	"testing"
)

func TestRenameConfigFromEnv_Empty(t *testing.T) {
	t.Setenv("VAULTPULL_RENAME_KEYS", "")
	rm := RenameConfigFromEnv()
	if len(rm) != 0 {
		t.Fatalf("expected empty rename map, got %v", rm)
	}
}

func TestRenameConfigFromEnv_MultipleKeys(t *testing.T) {
	t.Setenv("VAULTPULL_RENAME_KEYS", "X:A,Y:B,Z:C")
	rm := RenameConfigFromEnv()
	expected := map[string]string{"X": "A", "Y": "B", "Z": "C"}
	for k, v := range expected {
		if rm[k] != v {
			t.Fatalf("expected %s=%s, got %s", k, v, rm[k])
		}
	}
}

func TestRenameConfigFromEnv_TrimsWhitespace(t *testing.T) {
	t.Setenv("VAULTPULL_RENAME_KEYS", " FOO : BAR ")
	rm := RenameConfigFromEnv()
	if rm["FOO"] != "BAR" {
		t.Fatalf("expected BAR, got %q", rm["FOO"])
	}
}
