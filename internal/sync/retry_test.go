package sync

import (
	"errors"
	"testing"
	"time"
)

func TestWithRetry_SucceedsFirstAttempt(t *testing.T) {
	calls := 0
	err := WithRetry(DefaultRetryConfig(), func() error {
		calls++
		return nil
	})
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
	if calls != 1 {
		t.Fatalf("expected 1 call, got %d", calls)
	}
}

func TestWithRetry_RetriesOnRetryableError(t *testing.T) {
	calls := 0
	cfg := RetryConfig{MaxAttempts: 3, Delay: time.Millisecond, Multiplier: 1.0}
	err := WithRetry(cfg, func() error {
		calls++
		if calls < 3 {
			return &RetryableError{Cause: errors.New("transient")}
		}
		return nil
	})
	if err != nil {
		t.Fatalf("expected nil after retries, got %v", err)
	}
	if calls != 3 {
		t.Fatalf("expected 3 calls, got %d", calls)
	}
}

func TestWithRetry_StopsOnNonRetryableError(t *testing.T) {
	calls := 0
	cfg := RetryConfig{MaxAttempts: 5, Delay: time.Millisecond, Multiplier: 1.0}
	permanent := errors.New("permanent failure")
	err := WithRetry(cfg, func() error {
		calls++
		return permanent
	})
	if !errors.Is(err, permanent) {
		t.Fatalf("expected permanent error, got %v", err)
	}
	if calls != 1 {
		t.Fatalf("expected 1 call, got %d", calls)
	}
}

func TestWithRetry_ExhaustsAttempts(t *testing.T) {
	calls := 0
	cfg := RetryConfig{MaxAttempts: 3, Delay: time.Millisecond, Multiplier: 1.0}
	err := WithRetry(cfg, func() error {
		calls++
		return &RetryableError{Cause: errors.New("always fails")}
	})
	if err == nil {
		t.Fatal("expected error after exhausting attempts")
	}
	if calls != 3 {
		t.Fatalf("expected 3 calls, got %d", calls)
	}
}

func TestIsRetryable_True(t *testing.T) {
	err := &RetryableError{Cause: errors.New("oops")}
	if !IsRetryable(err) {
		t.Fatal("expected IsRetryable to be true")
	}
}

func TestIsRetryable_False(t *testing.T) {
	err := errors.New("plain error")
	if IsRetryable(err) {
		t.Fatal("expected IsRetryable to be false")
	}
}

func TestWithRetry_ZeroMaxAttempts_RunsOnce(t *testing.T) {
	calls := 0
	cfg := RetryConfig{MaxAttempts: 0, Delay: time.Millisecond, Multiplier: 1.0}
	_ = WithRetry(cfg, func() error {
		calls++
		return &RetryableError{Cause: errors.New("err")}
	})
	if calls != 1 {
		t.Fatalf("expected 1 call with zero MaxAttempts, got %d", calls)
	}
}

func TestRetryableError_Unwrap(t *testing.T) {
	cause := errors.New("underlying cause")
	err := &RetryableError{Cause: cause}
	if !errors.Is(err, cause) {
		t.Fatal("expected errors.Is to find the wrapped cause via Unwrap")
	}
}
