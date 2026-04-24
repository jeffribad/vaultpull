package sync

import (
	"testing"
)

func TestApplySuffixFilter_Disabled_ReturnsOriginal(t *testing.T) {
	secrets := map[string]string{"DB_URL": "postgres://localhost", "API_SECRET": "s3cr3t"}
	cfg := SuffixFilterConfig{Enabled: false}
	result := ApplySuffixFilter(secrets, cfg)
	if len(result) != 2 {
		t.Errorf("expected 2 secrets, got %d", len(result))
	}
}

func TestApplySuffixFilter_AllowSuffix(t *testing.T) {
	secrets := map[string]string{
		"DB_URL":     "postgres://localhost",
		"API_SECRET": "s3cr3t",
		"APP_HOST":   "localhost",
	}
	cfg := SuffixFilterConfig{Enabled: true, AllowSuffix: []string{"_URL", "_HOST"}}
	result := ApplySuffixFilter(secrets, cfg)
	if len(result) != 2 {
		t.Errorf("expected 2 secrets, got %d", len(result))
	}
	if _, ok := result["API_SECRET"]; ok {
		t.Error("API_SECRET should have been filtered out")
	}
}

func TestApplySuffixFilter_DenySuffix(t *testing.T) {
	secrets := map[string]string{
		"DB_URL":     "postgres://localhost",
		"API_SECRET": "s3cr3t",
		"JWT_KEY":    "abc123",
	}
	cfg := SuffixFilterConfig{Enabled: true, DenySuffix: []string{"_SECRET", "_KEY"}}
	result := ApplySuffixFilter(secrets, cfg)
	if len(result) != 1 {
		t.Errorf("expected 1 secret, got %d", len(result))
	}
	if _, ok := result["DB_URL"]; !ok {
		t.Error("DB_URL should be present")
	}
}

func TestApplySuffixFilter_AllowTakesPrecedenceOverDeny(t *testing.T) {
	secrets := map[string]string{
		"DB_URL":    "postgres://localhost",
		"APP_HOST":  "localhost",
		"JWT_TOKEN": "abc",
	}
	cfg := SuffixFilterConfig{
		Enabled:     true,
		AllowSuffix: []string{"_URL"},
		DenySuffix:  []string{"_URL"},
	}
	result := ApplySuffixFilter(secrets, cfg)
	if _, ok := result["DB_URL"]; !ok {
		t.Error("DB_URL should be kept: allow takes precedence")
	}
	if len(result) != 1 {
		t.Errorf("expected 1 secret, got %d", len(result))
	}
}

func TestApplySuffixFilter_CaseInsensitiveMatch(t *testing.T) {
	secrets := map[string]string{"db_url": "postgres://localhost", "api_secret": "s3cr3t"}
	cfg := SuffixFilterConfig{Enabled: true, AllowSuffix: []string{"_URL"}}
	result := ApplySuffixFilter(secrets, cfg)
	if _, ok := result["db_url"]; !ok {
		t.Error("db_url should match _URL suffix case-insensitively")
	}
}

func TestApplySuffixFilter_DoesNotMutateInput(t *testing.T) {
	secrets := map[string]string{"DB_URL": "postgres://localhost", "API_SECRET": "s3cr3t"}
	cfg := SuffixFilterConfig{Enabled: true, DenySuffix: []string{"_SECRET"}}
	ApplySuffixFilter(secrets, cfg)
	if len(secrets) != 2 {
		t.Error("original secrets map should not be mutated")
	}
}
