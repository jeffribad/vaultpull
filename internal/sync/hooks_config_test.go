package sync

import (
	"testing"
)

func TestHooksFromEnv_Empty(t *testing.T) {
	t.Setenv(envPreHooks, "")
	t.Setenv(envPostHooks, "")

	cfg := HooksFromEnv()
	if len(cfg.PreHooks) != 0 || len(cfg.PostHooks) != 0 {
		t.Error("expected empty hooks when env vars are unset")
	}
}

func TestHooksFromEnv_ParsesPreHooks(t *testing.T) {
	t.Setenv(envPreHooks, "echo a, echo b")
	t.Setenv(envPostHooks, "")

	cfg := HooksFromEnv()
	if len(cfg.PreHooks) != 2 {
		t.Fatalf("expected 2 pre hooks, got %d", len(cfg.PreHooks))
	}
	if cfg.PreHooks[0] != "echo a" || cfg.PreHooks[1] != "echo b" {
		t.Errorf("unexpected pre hooks: %v", cfg.PreHooks)
	}
}

func TestHooksFromEnv_ParsesPostHooks(t *testing.T) {
	t.Setenv(envPreHooks, "")
	t.Setenv(envPostHooks, "echo done")

	cfg := HooksFromEnv()
	if len(cfg.PostHooks) != 1 || cfg.PostHooks[0] != "echo done" {
		t.Errorf("unexpected post hooks: %v", cfg.PostHooks)
	}
}

func TestHooksConfig_ToHooks_Types(t *testing.T) {
	cfg := HooksConfig{
		PreHooks:  []string{"echo pre"},
		PostHooks: []string{"echo post"},
	}
	hooks := cfg.ToHooks()
	if len(hooks) != 2 {
		t.Fatalf("expected 2 hooks, got %d", len(hooks))
	}
	if hooks[0].Type != HookPreSync {
		t.Errorf("expected first hook to be pre-sync")
	}
	if hooks[1].Type != HookPostSync {
		t.Errorf("expected second hook to be post-sync")
	}
}

func TestParseHookList_TrimsWhitespace(t *testing.T) {
	result := parseHookList("  echo hello ,  echo world  ")
	if len(result) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(result))
	}
	if result[0] != "echo hello" || result[1] != "echo world" {
		t.Errorf("unexpected result: %v", result)
	}
}
