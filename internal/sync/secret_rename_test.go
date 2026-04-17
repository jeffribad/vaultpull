package sync

import (
	"os"
	"testing"
)

func TestParseRenameMap_Empty(t *testing.T) {
	rm := ParseRenameMap("")
	if len(rm) != 0 {
		t.Fatalf("expected empty map, got %v", rm)
	}
}

func TestParseRenameMap_SinglePair(t *testing.T) {
	rm := ParseRenameMap("OLD_KEY:NEW_KEY")
	if rm["OLD_KEY"] != "NEW_KEY" {
		t.Fatalf("expected NEW_KEY, got %q", rm["OLD_KEY"])
	}
}

func TestParseRenameMap_MultiplePairs(t *testing.T) {
	rm := ParseRenameMap("A:B, C:D")
	if rm["A"] != "B" || rm["C"] != "D" {
		t.Fatalf("unexpected map: %v", rm)
	}
}

func TestParseRenameMap_SkipsMalformed(t *testing.T) {
	rm := ParseRenameMap("NOCODON,X:Y")
	if _, ok := rm["NOCODON"]; ok {
		t.Fatal("malformed entry should be skipped")
	}
	if rm["X"] != "Y" {
		t.Fatalf("expected Y, got %q", rm["X"])
	}
}

func TestApplyRenames_NoRules(t *testing.T) {
	secrets := map[string]string{"FOO": "bar"}
	out := ApplyRenames(secrets, RenameMap{})
	if out["FOO"] != "bar" {
		t.Fatal("expected FOO to pass through")
	}
}

func TestApplyRenames_RenamesKey(t *testing.T) {
	secrets := map[string]string{"OLD": "val", "KEEP": "x"}
	rm := RenameMap{"OLD": "NEW"}
	out := ApplyRenames(secrets, rm)
	if _, ok := out["OLD"]; ok {
		t.Fatal("OLD should have been renamed")
	}
	if out["NEW"] != "val" {
		t.Fatalf("expected val, got %q", out["NEW"])
	}
	if out["KEEP"] != "x" {
		t.Fatal("KEEP should be unchanged")
	}
}

func TestRenameConfigFromEnv_ReadsEnv(t *testing.T) {
	t.Setenv("VAULTPULL_RENAME_KEYS", "DB_PASS:DATABASE_PASSWORD")
	rm := RenameConfigFromEnv()
	if rm["DB_PASS"] != "DATABASE_PASSWORD" {
		t.Fatalf("unexpected map: %v", rm)
	}
	os.Unsetenv("VAULTPULL_RENAME_KEYS")
}
