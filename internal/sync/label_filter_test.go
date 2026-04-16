package sync

import (
	"os"
	"testing"
)

func TestLabelFilterConfigFromEnv_Empty(t *testing.T) {
	os.Unsetenv("VAULTPULL_LABELS")
	cfg := LabelFilterConfigFromEnv()
	if len(cfg.Labels) != 0 {
		t.Fatalf("expected empty labels, got %v", cfg.Labels)
	}
}

func TestLabelFilterConfigFromEnv_ParsesPairs(t *testing.T) {
	os.Setenv("VAULTPULL_LABELS", "env=prod,team=backend")
	defer os.Unsetenv("VAULTPULL_LABELS")
	cfg := LabelFilterConfigFromEnv()
	if cfg.Labels["env"] != "prod" || cfg.Labels["team"] != "backend" {
		t.Fatalf("unexpected labels: %v", cfg.Labels)
	}
}

func TestApplyLabelFilter_NoFilter(t *testing.T) {
	secrets := map[string]string{"DB_PASS": "secret", "API_KEY": "key"}
	cfg := LabelFilterConfig{Labels: map[string]string{}}
	out := ApplyLabelFilter(secrets, nil, cfg)
	if len(out) != 2 {
		t.Fatalf("expected all secrets, got %d", len(out))
	}
}

func TestApplyLabelFilter_MatchingLabel(t *testing.T) {
	secrets := map[string]string{"DB_PASS": "s", "API_KEY": "k"}
	labels := map[string]map[string]string{
		"DB_PASS": {"env": "prod"},
		"API_KEY": {"env": "staging"},
	}
	cfg := LabelFilterConfig{Labels: map[string]string{"env": "prod"}}
	out := ApplyLabelFilter(secrets, labels, cfg)
	if len(out) != 1 || out["DB_PASS"] != "s" {
		t.Fatalf("expected only DB_PASS, got %v", out)
	}
}

func TestApplyLabelFilter_MultipleSelectors(t *testing.T) {
	secrets := map[string]string{"A": "1", "B": "2", "C": "3"}
	labels := map[string]map[string]string{
		"A": {"env": "prod", "team": "backend"},
		"B": {"env": "prod", "team": "frontend"},
		"C": {"env": "staging", "team": "backend"},
	}
	cfg := LabelFilterConfig{Labels: map[string]string{"env": "prod", "team": "backend"}}
	out := ApplyLabelFilter(secrets, labels, cfg)
	if len(out) != 1 || out["A"] != "1" {
		t.Fatalf("expected only A, got %v", out)
	}
}

func TestApplyLabelFilter_NoMatches(t *testing.T) {
	secrets := map[string]string{"X": "val"}
	labels := map[string]map[string]string{
		"X": {"env": "dev"},
	}
	cfg := LabelFilterConfig{Labels: map[string]string{"env": "prod"}}
	out := ApplyLabelFilter(secrets, labels, cfg)
	if len(out) != 0 {
		t.Fatalf("expected no secrets, got %v", out)
	}
}
