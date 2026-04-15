package sync

import (
	"errors"
	"time"
)

// RetryConfig holds configuration for retry behaviour.
type RetryConfig struct {
	MaxAttempts int
	Delay       time.Duration
	Multiplier  float64
}

// DefaultRetryConfig returns a sensible default retry configuration.
func DefaultRetryConfig() RetryConfig {
	return RetryConfig{
		MaxAttempts: 3,
		Delay:       500 * time.Millisecond,
		Multiplier:  2.0,
	}
}

// RetryableError wraps an error to indicate the operation can be retried.
type RetryableError struct {
	Cause error
}

func (e *RetryableError) Error() string {
	return e.Cause.Error()
}

func (e *RetryableError) Unwrap() error {
	return e.Cause
}

// IsRetryable reports whether err is a RetryableError.
func IsRetryable(err error) bool {
	var r *RetryableError
	return errors.As(err, &r)
}

// WithRetry executes fn up to cfg.MaxAttempts times, retrying only when
// the returned error satisfies IsRetryable. It returns the last error
// encountered, or nil on success.
func WithRetry(cfg RetryConfig, fn func() error) error {
	if cfg.MaxAttempts < 1 {
		cfg.MaxAttempts = 1
	}
	delay := cfg.Delay
	var err error
	for attempt := 1; attempt <= cfg.MaxAttempts; attempt++ {
		err = fn()
		if err == nil {
			return nil
		}
		if !IsRetryable(err) {
			return err
		}
		if attempt < cfg.MaxAttempts {
			time.Sleep(delay)
			delay = time.Duration(float64(delay) * cfg.Multiplier)
		}
	}
	return err
}
