package sync

import (
	"testing"
)

func TestApplyPrefixFilter_Disabled_ReturnsOriginal(t *testing.T) {
	secrets := map[string]string{"APP_KEY": "val", "DB_HOST": "host"}
	cfg := PrefixFilterConfig{Enabled: false}
	result := ApplyPrefixFilter(secrets, cfg)
	if len(result) != 2 {
		t.Errorf("expected 2 keys, got %d", len(result))
	}
}

func TestApplyPrefixFilter_AllowPrefix(t *testing.T) {
	secrets := map[string]string{"APP_KEY": "val", "DB_HOST": "host", "APP_SECRET": "s"}
	cfg := PrefixFilterConfig{Enabled: true, AllowPrefixes: []string{"APP_"}}
	result := ApplyPrefixFilter(secrets, cfg)
	if len(result) != 2 {
		t.Errorf("expected 2 keys, got %d", len(result))
	}
	if _, ok := result["DB_HOST"]; ok {
		t.Error("DB_HOST should have been filtered out")
	}
}

func TestApplyPrefixFilter_DenyPrefix(t *testing.T) {
	secrets := map[string]string{"APP_KEY": "val", "INTERNAL_TOKEN": "tok", "APP_ID": "id"}
	cfg := PrefixFilterConfig{Enabled: true, DenyPrefixes: []string{"INTERNAL_"}}
	result := ApplyPrefixFilter(secrets, cfg)
	if len(result) != 2 {
		t.Errorf("expected 2 keys, got %d", len(result))
	}
	if _, ok := result["INTERNAL_TOKEN"]; ok {
		t.Error("INTERNAL_TOKEN should have been denied")
	}
}

func TestApplyPrefixFilter_AllowAndDeny(t *testing.T) {
	secrets := map[string]string{"APP_KEY": "val", "APP_INTERNAL_SECRET": "s", "DB_HOST": "h"}
	cfg := PrefixFilterConfig{Enabled: true, AllowPrefixes: []string{"APP_"}, DenyPrefixes: []string{"APP_INTERNAL_"}}
	result := ApplyPrefixFilter(secrets, cfg)
	if len(result) != 1 {
		t.Errorf("expected 1 key, got %d", len(result))
	}
	if _, ok := result["APP_KEY"]; !ok {
		t.Error("APP_KEY should be present")
	}
}

func TestApplyPrefixFilter_CaseInsensitive(t *testing.T) {
	secrets := map[string]string{"app_key": "val", "db_host": "host"}
	cfg := PrefixFilterConfig{Enabled: true, AllowPrefixes: []string{"APP_"}}
	result := ApplyPrefixFilter(secrets, cfg)
	if len(result) != 1 {
		t.Errorf("expected 1 key, got %d", len(result))
	}
}

func TestPrefixFilterConfigFromEnv_Defaults(t *testing.T) {
	cfg := PrefixFilterConfigFromEnv()
	if cfg.Enabled {
		t.Error("expected disabled by default")
	}
	if len(cfg.AllowPrefixes) != 0 || len(cfg.DenyPrefixes) != 0 {
		t.Error("expected empty prefix lists by default")
	}
}
