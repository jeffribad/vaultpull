package sync

import (
	"os"
	"strings"
)

// CopyConfigFromEnv loads the env-copy configuration from environment variables.
// VAULTPULL_ENV_COPY_ENABLED: enable copying secrets to new keys
// VAULTPULL_ENV_COPY_PAIRS: comma-separated list of src:dst key pairs
type CopyConfig struct {
	Enabled bool
	Pairs   map[string]string // src -> dst
}

func CopyConfigFromEnv() CopyConfig {
	enabled := strings.TrimSpace(os.Getenv("VAULTPULL_ENV_COPY_ENABLED"))
	pairs := strings.TrimSpace(os.Getenv("VAULTPULL_ENV_COPY_PAIRS"))

	cfg := CopyConfig{
		Enabled: enabled == "true" || enabled == "1",
		Pairs:   make(map[string]string),
	}

	if pairs != "" {
		cfg.Pairs = parseCopyPairs(pairs)
	}

	return cfg
}

func parseCopyPairs(raw string) map[string]string {
	result := make(map[string]string)
	for _, part := range strings.Split(raw, ",") {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		idx := strings.Index(part, ":")
		if idx <= 0 || idx == len(part)-1 {
			continue
		}
		src := strings.TrimSpace(part[:idx])
		dst := strings.TrimSpace(part[idx+1:])
		if src != "" && dst != "" {
			result[src] = dst
		}
	}
	return result
}

// ApplyEnvCopy copies secret values from source keys to destination keys.
// Existing destination keys are overwritten. The source key is preserved.
func ApplyEnvCopy(cfg CopyConfig, secrets map[string]string) map[string]string {
	if !cfg.Enabled || len(cfg.Pairs) == 0 {
		return secrets
	}

	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		out[k] = v
	}

	for src, dst := range cfg.Pairs {
		val, ok := secrets[src]
		if !ok {
			// try case-insensitive lookup
			for k, v := range secrets {
				if strings.EqualFold(k, src) {
					val = v
					ok = true
					break
				}
			}
		}
		if ok {
			out[dst] = val
		}
	}

	return out
}
