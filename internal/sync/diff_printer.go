package sync

import (
	"fmt"
	"io"
	"sort"
)

// PrintDiff writes a human-readable diff report to w.
// Values are masked to avoid leaking secrets into terminal output.
func PrintDiff(w io.Writer, d DiffResult) {
	if !d.HasChanges() {
		fmt.Fprintln(w, "  (no changes)")
		return
	}

	if len(d.Added) > 0 {
		keys := sortedKeys(d.Added)
		for _, k := range keys {
			fmt.Fprintf(w, "  + %s\n", k)
		}
	}

	if len(d.Updated) > 0 {
		keys := sortedKeys(d.Updated)
		for _, k := range keys {
			fmt.Fprintf(w, "  ~ %s\n", k)
		}
	}

	for _, k := range d.Removed {
		fmt.Fprintf(w, "  - %s\n", k)
	}
}

func sortedKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
