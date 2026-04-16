package sync

import (
	"context"
	"fmt"
	"time"
)

// TimeoutConfig holds configuration for operation timeouts.
type TimeoutConfig struct {
	VaultTimeout  time.Duration
	WriteTimeout  time.Duration
	GlobalTimeout time.Duration
}

// DefaultTimeoutConfig returns sensible defaults.
func DefaultTimeoutConfig() TimeoutConfig {
	return TimeoutConfig{
		VaultTimeout:  10 * time.Second,
		WriteTimeout:  5 * time.Second,
		GlobalTimeout: 30 * time.Second,
	}
}

// WithTimeout wraps a function call with a context deadline.
func WithTimeout(ctx context.Context, d time.Duration, fn func(ctx context.Context) error) error {
	if d <= 0 {
		return fn(ctx)
	}
	tctx, cancel := context.WithTimeout(ctx, d)
	defer cancel()
	err := fn(tctx)
	if err != nil {
		if tctx.Err() == context.DeadlineExceeded {
			return fmt.Errorf("operation timed out after %s", d)
		}
		return err
	}
	return nil
}
