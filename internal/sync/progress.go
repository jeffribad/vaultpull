package sync

import (
	"fmt"
	"io"
	"os"
	"sync/atomic"
)

// ProgressReporter reports sync progress to a writer.
type ProgressReporter struct {
	out     io.Writer
	total   int
	current int32
	quiet   bool
}

// NewProgressReporter creates a new ProgressReporter.
func NewProgressReporter(total int, quiet bool) *ProgressReporter {
	return &ProgressReporter{
		out:   os.Stderr,
		total: total,
		quiet: quiet,
	}
}

// NewProgressReporterWithWriter creates a ProgressReporter writing to w.
func NewProgressReporterWithWriter(w io.Writer, total int, quiet bool) *ProgressReporter {
	return &ProgressReporter{out: w, total: total, quiet: quiet}
}

// Advance increments progress and prints the current state.
func (p *ProgressReporter) Advance(key string) {
	n := int(atomic.AddInt32(&p.current, 1))
	if p.quiet {
		return
	}
	fmt.Fprintf(p.out, "[%d/%d] syncing %s\n", n, p.total, key)
}

// Done prints a completion summary.
func (p *ProgressReporter) Done(written int, skipped int) {
	if p.quiet {
		return
	}
	fmt.Fprintf(p.out, "done: %d written, %d skipped\n", written, skipped)
}
