package sync

import (
	"fmt"
	"sort"
)

// DiffResult holds the changes between existing and incoming secrets.
type DiffResult struct {
	Added     map[string]string
	Updated   map[string]string
	Removed   []string
	Unchanged []string
}

// Diff compares existing key-value pairs with incoming ones and returns
// a categorised DiffResult. existing may be nil (e.g. new file).
func Diff(existing, incoming map[string]string) DiffResult {
	result := DiffResult{
		Added:   make(map[string]string),
		Updated: make(map[string]string),
	}

	for k, v := range incoming {
		oldVal, ok := existing[k]
		if !ok {
			result.Added[k] = v
		} else if oldVal != v {
			result.Updated[k] = v
		} else {
			result.Unchanged = append(result.Unchanged, k)
		}
	}

	for k := range existing {
		if _, ok := incoming[k]; !ok {
			result.Removed = append(result.Removed, k)
		}
	}

	sort.Strings(result.Removed)
	sort.Strings(result.Unchanged)

	return result
}

// HasChanges returns true if there is at least one add, update, or removal.
func (d DiffResult) HasChanges() bool {
	return len(d.Added) > 0 || len(d.Updated) > 0 || len(d.Removed) > 0
}

// Summary returns a human-readable one-liner describing the diff.
func (d DiffResult) Summary() string {
	if !d.HasChanges() {
		return "no changes detected"
	}
	return fmt.Sprintf("%d added, %d updated, %d removed",
		len(d.Added), len(d.Updated), len(d.Removed))
}
