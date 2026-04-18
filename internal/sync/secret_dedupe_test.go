package sync

import "testing"

func TestDedupeSecrets_Disabled_ReturnsOriginal(t *testing.T) {
	secrets := map[string]string{"KEY": "val1", "key": "val2"}
	cfg := DedupeConfig{Enabled: false}
	out := DedupeSecrets(secrets, cfg)
	if len(out) != 2 {
		t.Errorf("expected 2 keys, got %d", len(out))
	}
}

func TestDedupeSecrets_CaseSensitive_NoDuplicates(t *testing.T) {
	secrets := map[string]string{"KEY": "a", "key": "b"}
	cfg := DedupeConfig{Enabled: true, CaseSensitive: true}
	out := DedupeSecrets(secrets, cfg)
	if len(out) != 2 {
		t.Errorf("expected 2 keys (case-sensitive), got %d", len(out))
	}
}

func TestDedupeSecrets_CaseInsensitive_RemovesDuplicates(t *testing.T) {
	secrets := map[string]string{"KEY": "first", "key": "second"}
	cfg := DedupeConfig{Enabled: true, CaseSensitive: false}
	out := DedupeSecrets(secrets, cfg)
	if len(out) != 1 {
		t.Errorf("expected 1 key after dedup, got %d", len(out))
	}
}

func TestDedupeSecrets_CaseInsensitive_NoConflicts(t *testing.T) {
	secrets := map[string]string{"FOO": "bar", "BAZ": "qux"}
	cfg := DedupeConfig{Enabled: true, CaseSensitive: false}
	out := DedupeSecrets(secrets, cfg)
	if len(out) != 2 {
		t.Errorf("expected 2 keys, got %d", len(out))
	}
}

func TestDedupeSecrets_EmptyMap(t *testing.T) {
	secrets := map[string]string{}
	cfg := DedupeConfig{Enabled: true, CaseSensitive: false}
	out := DedupeSecrets(secrets, cfg)
	if len(out) != 0 {
		t.Errorf("expected empty map, got %d keys", len(out))
	}
}

func TestToLower_ConvertsUppercase(t *testing.T) {
	if got := toLower("HELLO_WORLD"); got != "hello_world" {
		t.Errorf("unexpected result: %s", got)
	}
}
