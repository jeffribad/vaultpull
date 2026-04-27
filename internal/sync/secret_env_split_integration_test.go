package sync

import (
	"testing"
)

func TestApplySplit_Integration_ConfigDriven(t *testing.T) {
	t.Setenv("VAULTPULL_SPLIT_ENABLED", "true")
	t.Setenv("VAULTPULL_SPLIT_SOURCE", "PACKED")
	t.Setenv("VAULTPULL_SPLIT_DELIMITER", ";")
	t.Setenv("VAULTPULL_SPLIT_SEPARATOR", "=")

	cfg := SplitConfigFromEnv()
	secrets := map[string]string{
		"PACKED": "API_URL=https://example.com;API_KEY=secret123",
		"OTHER":  "value",
	}

	out, err := ApplySplit(secrets, cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["API_URL"] != "https://example.com" {
		t.Errorf("expected API_URL=https://example.com, got %q", out["API_URL"])
	}
	if out["API_KEY"] != "secret123" {
		t.Errorf("expected API_KEY=secret123, got %q", out["API_KEY"])
	}
	if out["OTHER"] != "value" {
		t.Error("expected OTHER to be preserved")
	}
	if out["PACKED"] == "" {
		t.Error("expected PACKED source key to be preserved")
	}
}

func TestApplySplit_Integration_DisabledPassthrough(t *testing.T) {
	t.Setenv("VAULTPULL_SPLIT_ENABLED", "false")
	t.Setenv("VAULTPULL_SPLIT_SOURCE", "PACKED")

	cfg := SplitConfigFromEnv()
	secrets := map[string]string{
		"PACKED": "A=1,B=2",
	}

	out, err := ApplySplit(secrets, cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 1 {
		t.Errorf("expected 1 key when disabled, got %d", len(out))
	}
}
