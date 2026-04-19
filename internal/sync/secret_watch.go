package sync

import (
	"os"
	"strconv"
	"strings"
	"time"
)

// WatchConfig controls periodic secret re-sync behaviour.
type WatchConfig struct {
	Enabled  bool
	Interval time.Duration
	Keys     []string
}

// WatchConfigFromEnv loads WatchConfig from environment variables.
func WatchConfigFromEnv() WatchConfig {
	enabled := false
	if v := strings.TrimSpace(os.Getenv("VAULTPULL_WATCH_ENABLED")); v != "" {
		if b, err := strconv.ParseBool(v); err == nil {
			enabled = b
		} else if v == "1" {
			enabled = true
		}
	}

	interval := 60 * time.Second
	if v := strings.TrimSpace(os.Getenv("VAULTPULL_WATCH_INTERVAL_SECONDS")); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			interval = time.Duration(n) * time.Second
		}
	}

	var keys []string
	if v := os.Getenv("VAULTPULL_WATCH_KEYS"); v != "" {
		for _, k := range strings.Split(v, ",") {
			if t := strings.TrimSpace(k); t != "" {
				keys = append(keys, t)
			}
		}
	}

	return WatchConfig{
		Enabled:  enabled,
		Interval: interval,
		Keys:     keys,
	}
}

// WatchSecrets calls onChange whenever any watched key's value differs from
// the previous snapshot. It blocks until ctx is cancelled (via done channel).
func WatchSecrets(cfg WatchConfig, current map[string]string, fetch func() (map[string]string, error), onChange func(changed map[string]string), done <-chan struct{}) {
	if !cfg.Enabled {
		return
	}

	prev := copyMap(current)
	ticker := time.NewTicker(cfg.Interval)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			fresh, err := fetch()
			if err != nil {
				continue
			}
			changed := detectChanges(cfg.Keys, prev, fresh)
			if len(changed) > 0 {
				onChange(changed)
				prev = copyMap(fresh)
			}
		}
	}
}

func detectChanges(keys []string, prev, next map[string]string) map[string]string {
	changed := map[string]string{}
	watch := keys
	if len(watch) == 0 {
		for k := range next {
			watch = append(watch, k)
		}
	}
	for _, k := range watch {
		if next[k] != prev[k] {
			changed[k] = next[k]
		}
	}
	return changed
}

func copyMap(m map[string]string) map[string]string {
	out := make(map[string]string, len(m))
	for k, v := range m {
		out[k] = v
	}
	return out
}
