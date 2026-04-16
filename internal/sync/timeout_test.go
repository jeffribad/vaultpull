package sync

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"
)

func TestWithTimeout_CompletesInTime(t *testing.T) {
	err := WithTimeout(context.Background(), time.Second, func(ctx context.Context) error {
		return nil
	})
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestWithTimeout_TimesOut(t *testing.T) {
	err := WithTimeout(context.Background(), 10*time.Millisecond, func(ctx context.Context) error {
		time.Sleep(100 * time.Millisecond)
		return ctx.Err()
	})
	if err == nil {
		t.Fatal("expected timeout error")
	}
	if !strings.Contains(err.Error(), "timed out") {
		t.Errorf("unexpected error message: %s", err)
	}
}

func TestWithTimeout_PropagatesError(t *testing.T) {
	sentinel := errors.New("vault unavailable")
	err := WithTimeout(context.Background(), time.Second, func(ctx context.Context) error {
		return sentinel
	})
	if !errors.Is(err, sentinel) {
		t.Errorf("expected sentinel error, got %v", err)
	}
}

func TestWithTimeout_ZeroDuration_NoDeadline(t *testing.T) {
	called := false
	err := WithTimeout(context.Background(), 0, func(ctx context.Context) error {
		called = true
		return nil
	})
	if err != nil || !called {
		t.Fatalf("expected fn to be called without error")
	}
}

func TestDefaultTimeoutConfig(t *testing.T) {
	cfg := DefaultTimeoutConfig()
	if cfg.VaultTimeout != 10*time.Second {
		t.Errorf("unexpected VaultTimeout: %v", cfg.VaultTimeout)
	}
	if cfg.GlobalTimeout != 30*time.Second {
		t.Errorf("unexpected GlobalTimeout: %v", cfg.GlobalTimeout)
	}
}
