package sync

import (
	"os"
	"path/filepath"
	"testing"
)

func TestHooks_Integration_PreAndPostRunInOrder(t *testing.T) {
	dir := t.TempDir()
	markerPre := filepath.Join(dir, "pre.txt")
	markerPost := filepath.Join(dir, "post.txt")

	hooks := []Hook{
		{Type: HookPreSync, Command: "touch " + markerPre},
		{Type: HookPostSync, Command: "touch " + markerPost},
	}

	runner := NewHookRunner(hooks, os.Stdout)

	if err := runner.RunPre(); err != nil {
		t.Fatalf("RunPre failed: %v", err)
	}
	if _, err := os.Stat(markerPre); os.IsNotExist(err) {
		t.Error("pre hook did not create marker file")
	}
	if _, err := os.Stat(markerPost); !os.IsNotExist(err) {
		t.Error("post hook should not have run during RunPre")
	}

	if err := runner.RunPost(); err != nil {
		t.Fatalf("RunPost failed: %v", err)
	}
	if _, err := os.Stat(markerPost); os.IsNotExist(err) {
		t.Error("post hook did not create marker file")
	}
}

func TestHooks_Integration_FromEnv(t *testing.T) {
	dir := t.TempDir()
	marker := filepath.Join(dir, "env-hook.txt")

	t.Setenv(envPreHooks, "touch "+marker)
	t.Setenv(envPostHooks, "")

	cfg := HooksFromEnv()
	hooks := cfg.ToHooks()
	runner := NewHookRunner(hooks, os.Stdout)

	if err := runner.RunPre(); err != nil {
		t.Fatalf("RunPre from env failed: %v", err)
	}
	if _, err := os.Stat(marker); os.IsNotExist(err) {
		t.Error("env-configured pre hook did not run")
	}
}
