package sync

import (
	"testing"
)

func TestDiff_AllAdded(t *testing.T) {
	incoming := map[string]string{"FOO": "bar", "BAZ": "qux"}
	result := Diff(nil, incoming)

	if len(result.Added) != 2 {
		t.Errorf("expected 2 added, got %d", len(result.Added))
	}
	if len(result.Updated) != 0 {
		t.Errorf("expected 0 updated, got %d", len(result.Updated))
	}
	if len(result.Removed) != 0 {
		t.Errorf("expected 0 removed, got %d", len(result.Removed))
	}
}

func TestDiff_UpdatedKey(t *testing.T) {
	existing := map[string]string{"FOO": "old"}
	incoming := map[string]string{"FOO": "new"}
	result := Diff(existing, incoming)

	if len(result.Updated) != 1 {
		t.Errorf("expected 1 updated, got %d", len(result.Updated))
	}
	if result.Updated["FOO"] != "new" {
		t.Errorf("expected updated value 'new', got %q", result.Updated["FOO"])
	}
}

func TestDiff_RemovedKey(t *testing.T) {
	existing := map[string]string{"FOO": "bar", "GONE": "bye"}
	incoming := map[string]string{"FOO": "bar"}
	result := Diff(existing, incoming)

	if len(result.Removed) != 1 || result.Removed[0] != "GONE" {
		t.Errorf("expected GONE in removed, got %v", result.Removed)
	}
}

func TestDiff_UnchangedKey(t *testing.T) {
	existing := map[string]string{"FOO": "same"}
	incoming := map[string]string{"FOO": "same"}
	result := Diff(existing, incoming)

	if len(result.Unchanged) != 1 || result.Unchanged[0] != "FOO" {
		t.Errorf("expected FOO in unchanged, got %v", result.Unchanged)
	}
	if result.HasChanges() {
		t.Error("expected HasChanges to be false")
	}
}

func TestDiff_Summary_NoChanges(t *testing.T) {
	existing := map[string]string{"A": "1"}
	incoming := map[string]string{"A": "1"}
	result := Diff(existing, incoming)

	if result.Summary() != "no changes detected" {
		t.Errorf("unexpected summary: %s", result.Summary())
	}
}

func TestDiff_Summary_WithChanges(t *testing.T) {
	existing := map[string]string{"OLD": "x"}
	incoming := map[string]string{"NEW": "y"}
	result := Diff(existing, incoming)

	expected := "1 added, 0 updated, 1 removed"
	if result.Summary() != expected {
		t.Errorf("expected %q, got %q", expected, result.Summary())
	}
}
