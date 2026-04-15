package sync

import (
	"os"
	"strings"
)

const (
	envPreHooks  = "VAULTPULL_PRE_HOOKS"
	envPostHooks = "VAULTPULL_POST_HOOKS"
)

// HooksConfig holds pre- and post-sync hook commands.
type HooksConfig struct {
	PreHooks  []string
	PostHooks []string
}

// HooksFromEnv reads hook commands from environment variables.
// VAULTPULL_PRE_HOOKS and VAULTPULL_POST_HOOKS accept comma-separated commands.
func HooksFromEnv() HooksConfig {
	return HooksConfig{
		PreHooks:  parseHookList(os.Getenv(envPreHooks)),
		PostHooks: parseHookList(os.Getenv(envPostHooks)),
	}
}

// ToHooks converts a HooksConfig into a slice of Hook values.
func (c HooksConfig) ToHooks() []Hook {
	var hooks []Hook
	for _, cmd := range c.PreHooks {
		hooks = append(hooks, Hook{Type: HookPreSync, Command: cmd})
	}
	for _, cmd := range c.PostHooks {
		hooks = append(hooks, Hook{Type: HookPostSync, Command: cmd})
	}
	return hooks
}

func parseHookList(raw string) []string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil
	}
	parts := strings.Split(raw, ",")
	var out []string
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}
