package sync

import (
	"testing"
)

func TestApplyBlacklist_Disabled_ReturnsOriginal(t *testing.T) {
	cfg := BlacklistConfig{Enabled: false, Keys: []string{"SECRET"}}
	input := map[string]string{"SECRET": "value", "OTHER": "ok"}
	out := ApplyBlacklist(cfg, input)
	if len(out) != 2 {
		t.Errorf("expected 2 keys, got %d", len(out))
	}
}

func TestApplyBlacklist_NoKeys_ReturnsOriginal(t *testing.T) {
	cfg := BlacklistConfig{Enabled: true, Keys: nil}
	input := map[string]string{"SECRET": "value"}
	out := ApplyBlacklist(cfg, input)
	if len(out) != 1 {
		t.Errorf("expected 1 key, got %d", len(out))
	}
}

func TestApplyBlacklist_RemovesMatchingKey(t *testing.T) {
	cfg := BlacklistConfig{Enabled: true, Keys: []string{"SECRET_KEY"}}
	input := map[string]string{"SECRET_KEY": "s3cr3t", "SAFE_KEY": "public"}
	out := ApplyBlacklist(cfg, input)
	if _, ok := out["SECRET_KEY"]; ok {
		t.Error("expected SECRET_KEY to be removed")
	}
	if _, ok := out["SAFE_KEY"]; !ok {
		t.Error("expected SAFE_KEY to be present")
	}
}

func TestApplyBlacklist_CaseInsensitiveMatch(t *testing.T) {
	cfg := BlacklistConfig{Enabled: true, Keys: []string{"secret_key"}}
	input := map[string]string{"SECRET_KEY": "value", "OTHER": "ok"}
	out := ApplyBlacklist(cfg, input)
	if _, ok := out["SECRET_KEY"]; ok {
		t.Error("expected SECRET_KEY to be removed via case-insensitive match")
	}
}

func TestApplyBlacklist_DoesNotMutateInput(t *testing.T) {
	cfg := BlacklistConfig{Enabled: true, Keys: []string{"TOKEN"}}
	input := map[string]string{"TOKEN": "abc", "NAME": "app"}
	_ = ApplyBlacklist(cfg, input)
	if _, ok := input["TOKEN"]; !ok {
		t.Error("input map was mutated")
	}
}

func TestApplyBlacklist_MultipleKeys_RemovesAll(t *testing.T) {
	cfg := BlacklistConfig{Enabled: true, Keys: []string{"A", "B", "C"}}
	input := map[string]string{"A": "1", "B": "2", "C": "3", "D": "4"}
	out := ApplyBlacklist(cfg, input)
	if len(out) != 1 {
		t.Errorf("expected 1 key remaining, got %d", len(out))
	}
	if _, ok := out["D"]; !ok {
		t.Error("expected D to remain")
	}
}
