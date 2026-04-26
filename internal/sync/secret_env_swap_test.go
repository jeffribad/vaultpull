package sync

import (
	"testing"
)

func TestApplySwap_Disabled_ReturnsOriginal(t *testing.T) {
	cfg := SwapConfig{Enabled: false, Pairs: [][2]string{{"A", "B"}}}
	secrets := map[string]string{"A": "alpha", "B": "beta"}

	out := ApplySwap(cfg, secrets)

	if out["A"] != "alpha" || out["B"] != "beta" {
		t.Errorf("expected original values, got A=%s B=%s", out["A"], out["B"])
	}
}

func TestApplySwap_NoPairs_ReturnsOriginal(t *testing.T) {
	cfg := SwapConfig{Enabled: true}
	secrets := map[string]string{"X": "x-val"}

	out := ApplySwap(cfg, secrets)

	if out["X"] != "x-val" {
		t.Errorf("expected x-val, got %s", out["X"])
	}
}

func TestApplySwap_SwapsTwoKeys(t *testing.T) {
	cfg := SwapConfig{Enabled: true, Pairs: [][2]string{{"DB_HOST", "CACHE_HOST"}}}
	secrets := map[string]string{
		"DB_HOST":    "db.internal",
		"CACHE_HOST": "cache.internal",
	}

	out := ApplySwap(cfg, secrets)

	if out["DB_HOST"] != "cache.internal" {
		t.Errorf("expected cache.internal, got %s", out["DB_HOST"])
	}
	if out["CACHE_HOST"] != "db.internal" {
		t.Errorf("expected db.internal, got %s", out["CACHE_HOST"])
	}
}

func TestApplySwap_MissingKey_SkipsPair(t *testing.T) {
	cfg := SwapConfig{Enabled: true, Pairs: [][2]string{{"A", "MISSING"}}}
	secrets := map[string]string{"A": "aval"}

	out := ApplySwap(cfg, secrets)

	if out["A"] != "aval" {
		t.Errorf("expected aval unchanged, got %s", out["A"])
	}
}

func TestApplySwap_DoesNotMutateInput(t *testing.T) {
	cfg := SwapConfig{Enabled: true, Pairs: [][2]string{{"P", "Q"}}}
	secrets := map[string]string{"P": "pval", "Q": "qval"}

	_ = ApplySwap(cfg, secrets)

	if secrets["P"] != "pval" || secrets["Q"] != "qval" {
		t.Error("input map was mutated")
	}
}

func TestSwapConfigFromEnv_Defaults(t *testing.T) {
	t.Setenv("VAULTPULL_SWAP_ENABLED", "")
	t.Setenv("VAULTPULL_SWAP_PAIRS", "")

	cfg := SwapConfigFromEnv()

	if cfg.Enabled {
		t.Error("expected Enabled=false by default")
	}
	if len(cfg.Pairs) != 0 {
		t.Errorf("expected no pairs, got %d", len(cfg.Pairs))
	}
}

func TestSwapConfigFromEnv_ParsesPairs(t *testing.T) {
	t.Setenv("VAULTPULL_SWAP_ENABLED", "true")
	t.Setenv("VAULTPULL_SWAP_PAIRS", "FOO:BAR, BAZ:QUX ")

	cfg := SwapConfigFromEnv()

	if !cfg.Enabled {
		t.Error("expected Enabled=true")
	}
	if len(cfg.Pairs) != 2 {
		t.Fatalf("expected 2 pairs, got %d", len(cfg.Pairs))
	}
	if cfg.Pairs[0] != [2]string{"FOO", "BAR"} {
		t.Errorf("unexpected first pair: %v", cfg.Pairs[0])
	}
	if cfg.Pairs[1] != [2]string{"BAZ", "QUX"} {
		t.Errorf("unexpected second pair: %v", cfg.Pairs[1])
	}
}
