package sync

import (
	"os"
	"testing"
)

func TestApplyJoin_Disabled_ReturnsOriginal(t *testing.T) {
	secrets := map[string]string{"A": "1", "B": "2"}
	cfg := JoinConfig{Enabled: false, SourceKeys: []string{"A", "B"}, DestKey: "C", Separator: "-"}
	out, err := ApplyJoin(secrets, cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := out["C"]; ok {
		t.Error("expected C to be absent when disabled")
	}
}

func TestApplyJoin_NoSourceKeys_ReturnsOriginal(t *testing.T) {
	secrets := map[string]string{"A": "1"}
	cfg := JoinConfig{Enabled: true, SourceKeys: nil, DestKey: "C", Separator: ","}
	out, err := ApplyJoin(secrets, cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != len(secrets) {
		t.Errorf("expected map unchanged, got %v", out)
	}
}

func TestApplyJoin_JoinsValues(t *testing.T) {
	secrets := map[string]string{"DB_HOST": "localhost", "DB_PORT": "5432"}
	cfg := JoinConfig{Enabled: true, SourceKeys: []string{"DB_HOST", "DB_PORT"}, DestKey: "DB_ADDR", Separator: ":"}
	out, err := ApplyJoin(secrets, cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := out["DB_ADDR"]; got != "localhost:5432" {
		t.Errorf("expected 'localhost:5432', got %q", got)
	}
}

func TestApplyJoin_CaseInsensitiveSourceKey(t *testing.T) {
	secrets := map[string]string{"db_host": "127.0.0.1", "db_port": "3306"}
	cfg := JoinConfig{Enabled: true, SourceKeys: []string{"DB_HOST", "DB_PORT"}, DestKey: "DB_ADDR", Separator: ":"}
	out, err := ApplyJoin(secrets, cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := out["DB_ADDR"]; got != "127.0.0.1:3306" {
		t.Errorf("expected '127.0.0.1:3306', got %q", got)
	}
}

func TestApplyJoin_MissingSourceKey_ReturnsError(t *testing.T) {
	secrets := map[string]string{"A": "1"}
	cfg := JoinConfig{Enabled: true, SourceKeys: []string{"A", "MISSING"}, DestKey: "C", Separator: ","}
	_, err := ApplyJoin(secrets, cfg)
	if err == nil {
		t.Error("expected error for missing source key")
	}
}

func TestApplyJoin_DoesNotMutateInput(t *testing.T) {
	secrets := map[string]string{"X": "foo", "Y": "bar"}
	copy := map[string]string{"X": "foo", "Y": "bar"}
	cfg := JoinConfig{Enabled: true, SourceKeys: []string{"X", "Y"}, DestKey: "Z", Separator: "-"}
	_, _ = ApplyJoin(secrets, cfg)
	for k, v := range copy {
		if secrets[k] != v {
			t.Errorf("input mutated at key %q", k)
		}
	}
	if _, ok := secrets["Z"]; ok {
		t.Error("input mutated: Z should not exist in original")
	}
}

func TestJoinConfigFromEnv_Defaults(t *testing.T) {
	os.Unsetenv("VAULTPULL_JOIN_ENABLED")
	os.Unsetenv("VAULTPULL_JOIN_SOURCE_KEYS")
	os.Unsetenv("VAULTPULL_JOIN_DEST_KEY")
	os.Unsetenv("VAULTPULL_JOIN_SEPARATOR")
	cfg := JoinConfigFromEnv()
	if cfg.Enabled {
		t.Error("expected Enabled=false by default")
	}
	if cfg.Separator != "," {
		t.Errorf("expected default separator ',', got %q", cfg.Separator)
	}
}

func TestJoinConfigFromEnv_ParsesValues(t *testing.T) {
	t.Setenv("VAULTPULL_JOIN_ENABLED", "true")
	t.Setenv("VAULTPULL_JOIN_SOURCE_KEYS", "HOST,PORT")
	t.Setenv("VAULTPULL_JOIN_DEST_KEY", "ADDR")
	t.Setenv("VAULTPULL_JOIN_SEPARATOR", ":")
	cfg := JoinConfigFromEnv()
	if !cfg.Enabled {
		t.Error("expected Enabled=true")
	}
	if len(cfg.SourceKeys) != 2 || cfg.SourceKeys[0] != "HOST" {
		t.Errorf("unexpected SourceKeys: %v", cfg.SourceKeys)
	}
	if cfg.DestKey != "ADDR" {
		t.Errorf("expected DestKey 'ADDR', got %q", cfg.DestKey)
	}
	if cfg.Separator != ":" {
		t.Errorf("expected separator ':', got %q", cfg.Separator)
	}
}
