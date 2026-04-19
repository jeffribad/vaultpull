package sync

import (
	"testing"
	"time"
)

func TestWatchSecrets_Disabled_NeverCallsOnChange(t *testing.T) {
	cfg := WatchConfig{Enabled: false, Interval: 10 * time.Millisecond}
	called := false
	done := make(chan struct{})
	close(done)
	WatchSecrets(cfg, map[string]string{"K": "v"}, func() (map[string]string, error) {
		return map[string]string{"K": "changed"}, nil
	}, func(_ map[string]string) {
		called = true
	}, done)
	if called {
		t.Error("onChange should not be called when disabled")
	}
}

func TestWatchSecrets_DetectsChange(t *testing.T) {
	cfg := WatchConfig{Enabled: true, Interval: 10 * time.Millisecond}
	initial := map[string]string{"KEY": "old"}
	changedCh := make(chan map[string]string, 1)
	done := make(chan struct{})

	go WatchSecrets(cfg, initial, func() (map[string]string, error) {
		return map[string]string{"KEY": "new"}, nil
	}, func(changed map[string]string) {
		changedCh <- changed
		close(done)
	}, done)

	select {
	case got := <-changedCh:
		if got["KEY"] != "new" {
			t.Errorf("expected 'new', got %q", got["KEY"])
		}
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for change detection")
	}
}

func TestDetectChanges_AllKeys(t *testing.T) {
	prev := map[string]string{"A": "1", "B": "2"}
	next := map[string]string{"A": "1", "B": "99"}
	got := detectChanges(nil, prev, next)
	if len(got) != 1 || got["B"] != "99" {
		t.Errorf("unexpected changes: %v", got)
	}
}

func TestDetectChanges_FilteredKeys(t *testing.T) {
	prev := map[string]string{"A": "1", "B": "2"}
	next := map[string]string{"A": "changed", "B": "changed"}
	got := detectChanges([]string{"B"}, prev, next)
	if _, ok := got["A"]; ok {
		t.Error("A should not be in changed set when not watched")
	}
	if got["B"] != "changed" {
		t.Errorf("expected B=changed, got %v", got)
	}
}

func TestCopyMap_DoesNotShareReference(t *testing.T) {
	orig := map[string]string{"x": "1"}
	copy := copyMap(orig)
	copy["x"] = "mutated"
	if orig["x"] != "1" {
		t.Error("copyMap should not share underlying map")
	}
}
