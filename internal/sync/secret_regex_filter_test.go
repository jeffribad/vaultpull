package sync

import (
	"testing"
)

func TestApplyRegexFilter_Disabled_ReturnsOriginal(t *testing.T) {
	cfg := RegexFilterConfig{Enabled: false, AllowPattern: "^DB_"}
	secrets := map[string]string{"DB_HOST": "localhost", "API_KEY": "abc"}
	out, err := ApplyRegexFilter(cfg, secrets)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != len(secrets) {
		t.Errorf("expected %d keys, got %d", len(secrets), len(out))
	}
}

func TestApplyRegexFilter_AllowPattern_FiltersKeys(t *testing.T) {
	cfg := RegexFilterConfig{Enabled: true, AllowPattern: "^DB_"}
	secrets := map[string]string{"DB_HOST": "localhost", "DB_PORT": "5432", "API_KEY": "abc"}
	out, err := ApplyRegexFilter(cfg, secrets)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 2 {
		t.Errorf("expected 2 keys, got %d", len(out))
	}
	if _, ok := out["API_KEY"]; ok {
		t.Error("API_KEY should have been filtered out")
	}
}

func TestApplyRegexFilter_DenyPattern_RemovesKeys(t *testing.T) {
	cfg := RegexFilterConfig{Enabled: true, DenyPattern: "SECRET|TOKEN"}
	secrets := map[string]string{"DB_HOST": "localhost", "API_SECRET": "xyz", "AUTH_TOKEN": "tok"}
	out, err := ApplyRegexFilter(cfg, secrets)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 1 {
		t.Errorf("expected 1 key, got %d", len(out))
	}
	if _, ok := out["DB_HOST"]; !ok {
		t.Error("DB_HOST should be present")
	}
}

func TestApplyRegexFilter_AllowAndDeny_AllowFirst(t *testing.T) {
	cfg := RegexFilterConfig{Enabled: true, AllowPattern: "^DB_", DenyPattern: "_SECRET$"}
	secrets := map[string]string{"DB_HOST": "localhost", "DB_SECRET": "s3cr3t", "API_KEY": "abc"}
	out, err := ApplyRegexFilter(cfg, secrets)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 1 {
		t.Errorf("expected 1 key, got %d", len(out))
	}
	if _, ok := out["DB_HOST"]; !ok {
		t.Error("DB_HOST should be present")
	}
}

func TestApplyRegexFilter_InvalidAllowPattern_ReturnsError(t *testing.T) {
	cfg := RegexFilterConfig{Enabled: true, AllowPattern: "[invalid"}
	secrets := map[string]string{"KEY": "val"}
	_, err := ApplyRegexFilter(cfg, secrets)
	if err == nil {
		t.Error("expected error for invalid allow pattern")
	}
}

func TestApplyRegexFilter_InvalidDenyPattern_ReturnsError(t *testing.T) {
	cfg := RegexFilterConfig{Enabled: true, DenyPattern: "[invalid"}
	secrets := map[string]string{"KEY": "val"}
	_, err := ApplyRegexFilter(cfg, secrets)
	if err == nil {
		t.Error("expected error for invalid deny pattern")
	}
}

func TestApplyRegexFilter_DoesNotMutateInput(t *testing.T) {
	cfg := RegexFilterConfig{Enabled: true, AllowPattern: "^KEEP_"}
	secrets := map[string]string{"KEEP_ME": "yes", "DROP_ME": "no"}
	orig := map[string]string{"KEEP_ME": "yes", "DROP_ME": "no"}
	_, _ = ApplyRegexFilter(cfg, secrets)
	if len(secrets) != len(orig) {
		t.Error("input map was mutated")
	}
}
