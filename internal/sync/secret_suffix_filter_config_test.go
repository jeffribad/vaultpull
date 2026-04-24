package sync

import (
	"os"
	"testing"
)

func TestSuffixFilterConfigFromEnv_Defaults(t *testing.T) {
	os.Unsetenv("VAULTPULL_SUFFIX_FILTER_ENABLED")
	os.Unsetenv("VAULTPULL_SUFFIX_FILTER_ALLOW")
	os.Unsetenv("VAULTPULL_SUFFIX_FILTER_DENY")

	cfg := SuffixFilterConfigFromEnv()
	if cfg.Enabled {
		t.Error("expected Enabled to be false by default")
	}
	if len(cfg.AllowSuffix) != 0 {
		t.Errorf("expected empty AllowSuffix, got %v", cfg.AllowSuffix)
	}
	if len(cfg.DenySuffix) != 0 {
		t.Errorf("expected empty DenySuffix, got %v", cfg.DenySuffix)
	}
}

func TestSuffixFilterConfigFromEnv_Enabled(t *testing.T) {
	t.Setenv("VAULTPULL_SUFFIX_FILTER_ENABLED", "true")
	cfg := SuffixFilterConfigFromEnv()
	if !cfg.Enabled {
		t.Error("expected Enabled to be true")
	}
}

func TestSuffixFilterConfigFromEnv_NumericEnabled(t *testing.T) {
	t.Setenv("VAULTPULL_SUFFIX_FILTER_ENABLED", "1")
	cfg := SuffixFilterConfigFromEnv()
	if !cfg.Enabled {
		t.Error("expected Enabled to be true when set to '1'")
	}
}

func TestSuffixFilterConfigFromEnv_ParsesAllowAndDeny(t *testing.T) {
	t.Setenv("VAULTPULL_SUFFIX_FILTER_ENABLED", "true")
	t.Setenv("VAULTPULL_SUFFIX_FILTER_ALLOW", "_URL, _HOST")
	t.Setenv("VAULTPULL_SUFFIX_FILTER_DENY", "_SECRET,_KEY")

	cfg := SuffixFilterConfigFromEnv()
	if len(cfg.AllowSuffix) != 2 {
		t.Errorf("expected 2 allow suffixes, got %d", len(cfg.AllowSuffix))
	}
	if cfg.AllowSuffix[0] != "_URL" {
		t.Errorf("expected _URL, got %s", cfg.AllowSuffix[0])
	}
	if len(cfg.DenySuffix) != 2 {
		t.Errorf("expected 2 deny suffixes, got %d", len(cfg.DenySuffix))
	}
	if cfg.DenySuffix[1] != "_KEY" {
		t.Errorf("expected _KEY, got %s", cfg.DenySuffix[1])
	}
}
