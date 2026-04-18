package sync

import (
	"testing"
)

func TestParseAliasMap_Empty(t *testing.T) {
	m := ParseAliasMap("")
	if len(m) != 0 {
		t.Fatalf("expected empty map, got %v", m)
	}
}

func TestParseAliasMap_SinglePair(t *testing.T) {
	m := ParseAliasMap("DB_HOST:DATABASE_HOST")
	aliases, ok := m["DB_HOST"]
	if !ok || len(aliases) != 1 || aliases[0] != "DATABASE_HOST" {
		t.Fatalf("unexpected map: %v", m)
	}
}

func TestParseAliasMap_MultiplePairs(t *testing.T) {
	m := ParseAliasMap("DB_HOST:DATABASE_HOST,REDIS_PASS:CACHE_PASSWORD")
	if len(m) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(m))
	}
}

func TestParseAliasMap_MultipleAliasesForSameKey(t *testing.T) {
	m := ParseAliasMap("DB_HOST:HOST_A,DB_HOST:HOST_B")
	aliases := m["DB_HOST"]
	if len(aliases) != 2 {
		t.Fatalf("expected 2 aliases for DB_HOST, got %v", aliases)
	}
}

func TestParseAliasMap_SkipsMalformed(t *testing.T) {
	m := ParseAliasMap("BADENTRY,:NOKEY,NOVAL:")
	if len(m) != 0 {
		t.Fatalf("expected empty map for malformed input, got %v", m)
	}
}

func TestApplyAliases_Disabled_ReturnsOriginal(t *testing.T) {
	secrets := map[string]string{"DB_HOST": "localhost"}
	cfg := AliasConfig{Enabled: false, Aliases: map[string][]string{"DB_HOST": {"DATABASE_HOST"}}}
	out := ApplyAliases(secrets, cfg)
	if _, ok := out["DATABASE_HOST"]; ok {
		t.Fatal("alias should not be injected when disabled")
	}
}

func TestApplyAliases_InjectsAlias(t *testing.T) {
	secrets := map[string]string{"DB_HOST": "localhost"}
	cfg := AliasConfig{Enabled: true, Aliases: map[string][]string{"DB_HOST": {"DATABASE_HOST"}}}
	out := ApplyAliases(secrets, cfg)
	if out["DATABASE_HOST"] != "localhost" {
		t.Fatalf("expected alias to have value 'localhost', got %q", out["DATABASE_HOST"])
	}
	if out["DB_HOST"] != "localhost" {
		t.Fatal("original key should be preserved")
	}
}

func TestApplyAliases_DoesNotOverwriteExisting(t *testing.T) {
	secrets := map[string]string{"DB_HOST": "localhost", "DATABASE_HOST": "other"}
	cfg := AliasConfig{Enabled: true, Aliases: map[string][]string{"DB_HOST": {"DATABASE_HOST"}}}
	out := ApplyAliases(secrets, cfg)
	if out["DATABASE_HOST"] != "other" {
		t.Fatal("existing key should not be overwritten by alias")
	}
}

func TestApplyAliases_SkipsMissingOriginal(t *testing.T) {
	secrets := map[string]string{"OTHER": "val"}
	cfg := AliasConfig{Enabled: true, Aliases: map[string][]string{"DB_HOST": {"DATABASE_HOST"}}}
	out := ApplyAliases(secrets, cfg)
	if _, ok := out["DATABASE_HOST"]; ok {
		t.Fatal("alias should not be created when original is missing")
	}
}
