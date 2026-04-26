package sync

import (
	"os"
	"strings"
)

// SwapConfig controls key-swap behaviour: swapping the values of two secrets.
type SwapConfig struct {
	Enabled bool
	Pairs   [][2]string // each pair is [keyA, keyB] whose values should be swapped
}

// SwapConfigFromEnv reads swap configuration from environment variables.
//
//	VAULTPULL_SWAP_ENABLED=true
//	VAULTPULL_SWAP_PAIRS=KEY_A:KEY_B,KEY_C:KEY_D
func SwapConfigFromEnv() SwapConfig {
	cfg := SwapConfig{}

	raw := os.Getenv("VAULTPULL_SWAP_ENABLED")
	cfg.Enabled = raw == "true" || raw == "1"

	pairsRaw := os.Getenv("VAULTPULL_SWAP_PAIRS")
	if pairsRaw == "" {
		return cfg
	}

	for _, token := range strings.Split(pairsRaw, ",") {
		token = strings.TrimSpace(token)
		parts := strings.SplitN(token, ":", 2)
		if len(parts) != 2 {
			continue
		}
		a := strings.TrimSpace(parts[0])
		b := strings.TrimSpace(parts[1])
		if a != "" && b != "" {
			cfg.Pairs = append(cfg.Pairs, [2]string{a, b})
		}
	}

	return cfg
}

// ApplySwap swaps the values of configured key pairs in the secrets map.
// If either key in a pair is absent the pair is skipped.
// The original map is never mutated; a new map is returned.
func ApplySwap(cfg SwapConfig, secrets map[string]string) map[string]string {
	if !cfg.Enabled || len(cfg.Pairs) == 0 {
		return secrets
	}

	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		out[k] = v
	}

	for _, pair := range cfg.Pairs {
		a, b := pair[0], pair[1]
		valA, okA := out[a]
		valB, okB := out[b]
		if !okA || !okB {
			continue
		}
		out[a] = valB
		out[b] = valA
	}

	return out
}
