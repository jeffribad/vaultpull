package sync

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// HookType represents the lifecycle point at which a hook runs.
type HookType string

const (
	HookPreSync  HookType = "pre"
	HookPostSync HookType = "post"
)

// Hook defines a shell command to run at a lifecycle point.
type Hook struct {
	Type    HookType
	Command string
}

// HookRunner executes lifecycle hooks around a sync operation.
type HookRunner struct {
	hooks []Hook
	out   *os.File
}

// NewHookRunner creates a HookRunner with the provided hooks.
func NewHookRunner(hooks []Hook, out *os.File) *HookRunner {
	if out == nil {
		out = os.Stdout
	}
	return &HookRunner{hooks: hooks, out: out}
}

// RunPre executes all pre-sync hooks in order.
func (h *HookRunner) RunPre() error {
	return h.run(HookPreSync)
}

// RunPost executes all post-sync hooks in order.
func (h *HookRunner) RunPost() error {
	return h.run(HookPostSync)
}

func (h *HookRunner) run(t HookType) error {
	for _, hook := range h.hooks {
		if hook.Type != t {
			continue
		}
		if err := h.execute(hook.Command); err != nil {
			return fmt.Errorf("hook %q failed: %w", hook.Command, err)
		}
	}
	return nil
}

func (h *HookRunner) execute(command string) error {
	parts := strings.Fields(command)
	if len(parts) == 0 {
		return fmt.Errorf("empty hook command")
	}
	cmd := exec.Command(parts[0], parts[1:]...) //nolint:gosec
	cmd.Stdout = h.out
	cmd.Stderr = h.out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("command %q exited with error: %w", parts[0], err)
	}
	return nil
}

// Count returns the number of hooks registered for the given type.
func (h *HookRunner) Count(t HookType) int {
	n := 0
	for _, hook := range h.hooks {
		if hook.Type == t {
			n++
		}
	}
	return n
}
