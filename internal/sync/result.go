package sync

import (
	"fmt"
	"io"
)

// Summary prints a human-readable summary of a sync Result to w.
func (r *Result) Summary(w io.Writer) {
	if r.Written == 0 && r.Skipped > 0 {
		fmt.Fprintf(w, "dry-run: %d secret(s) would be written to %s\n", r.Skipped, r.FilePath)
		return
	}
	fmt.Fprintf(w, "synced %d secret(s) → %s\n", r.Written, r.FilePath)
}

// IsEmpty reports whether no secrets were written or previewed.
func (r *Result) IsEmpty() bool {
	return r.Written == 0 && r.Skipped == 0
}
