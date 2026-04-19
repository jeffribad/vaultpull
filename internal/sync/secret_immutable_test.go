package sync

import (
	"testing"
)

func TestEnforceImmutable_Disabled_ReturnsNil(t *testing.T) {
	cfg := ImmutableConfig{Enabled: false, Keys: []string{"DB_PASS"}}
	current := map[string]string{"DB_PASS": "old"}
	incoming := map[string]string{"DB_PASS": "new"}
	if err := EnforceImmutable(cfg, current, incoming); err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestEnforceImmutable_NoKeys_ReturnsNil(t *testing.T) {
	cfg := ImmutableConfig{Enabled: true, Keys: nil}
	current := map[string]string{"DB_PASS": "old"}
	incoming := map[string]string{"DB_PASS": "new"}
	if err := EnforceImmutable(cfg, current, incoming); err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestEnforceImmutable_NoExistingValue_Allowed(t *testing.T) {
	cfg := ImmutableConfig{Enabled: true, Keys: []string{"DB_PASS"}}
	current := map[string]string{}
	incoming := map[string]string{"DB_PASS": "new"}
	if err := EnforceImmutable(cfg, current, incoming); err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestEnforceImmutable_SameValue_Allowed(t *testing.T) {
	cfg := ImmutableConfig{Enabled: true, Keys: []string{"DB_PASS"}}
	current := map[string]string{"DB_PASS": "secret"}
	incoming := map[string]string{"DB_PASS": "secret"}
	if err := EnforceImmutable(cfg, current, incoming); err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestEnforceImmutable_ChangedValue_ReturnsError(t *testing.T) {
	cfg := ImmutableConfig{Enabled: true, Keys: []string{"DB_PASS"}}
	current := map[string]string{"DB_PASS": "old"}
	incoming := map[string]string{"DB_PASS": "new"}
	if err := EnforceImmutable(cfg, current, incoming); err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestEnforceImmutable_CaseInsensitiveKey(t *testing.T) {
	cfg := ImmutableConfig{Enabled: true, Keys: []string{"db_pass"}}
	current := map[string]string{"DB_PASS": "old"}
	incoming := map[string]string{"DB_PASS": "changed"}
	if err := EnforceImmutable(cfg, current, incoming); err == nil {
		t.Fatal("expected error for immutable key change")
	}
}
