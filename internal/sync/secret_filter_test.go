package sync

import (
	"testing"
)

var baseSecrets = map[string]string{
	"DB_HOST":      "localhost",
	"DB_PORT":      "5432",
	"API_KEY":      "abc123",
	"SECRET_TOKEN": "supersecret",
}

func TestApplySecretFilter_NoFilter(t *testing.T) {
	cfg := SecretFilterConfig{}
	out := ApplySecretFilter(baseSecrets, cfg)
	if len(out) != len(baseSecrets) {
		t.Errorf("expected %d keys, got %d", len(baseSecrets), len(out))
	}
}

func TestApplySecretFilter_IncludeKeys(t *testing.T) {
	cfg := SecretFilterConfig{IncludeKeys: []string{"DB_HOST", "DB_PORT"}}
	out := ApplySecretFilter(baseSecrets, cfg)
	if len(out) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(out))
	}
	if _, ok := out["API_KEY"]; ok {
		t.Error("API_KEY should have been excluded")
	}
}

func TestApplySecretFilter_ExcludeKeys(t *testing.T) {
	cfg := SecretFilterConfig{ExcludeKeys: []string{"SECRET_TOKEN"}}
	out := ApplySecretFilter(baseSecrets, cfg)
	if _, ok := out["SECRET_TOKEN"]; ok {
		t.Error("SECRET_TOKEN should have been excluded")
	}
	if len(out) != 3 {
		t.Errorf("expected 3 keys, got %d", len(out))
	}
}

func TestApplySecretFilter_IncludeAndExclude(t *testing.T) {
	cfg := SecretFilterConfig{
		IncludeKeys: []string{"DB_HOST", "DB_PORT", "API_KEY"},
		ExcludeKeys: []string{"API_KEY"},
	}
	out := ApplySecretFilter(baseSecrets, cfg)
	if len(out) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(out))
	}
	if _, ok := out["API_KEY"]; ok {
		t.Error("API_KEY should have been excluded")
	}
}

func TestApplySecretFilter_CaseInsensitiveMatch(t *testing.T) {
	cfg := SecretFilterConfig{ExcludeKeys: []string{"db_host"}}
	out := ApplySecretFilter(baseSecrets, cfg)
	if _, ok := out["DB_HOST"]; ok {
		t.Error("DB_HOST should have been excluded by case-insensitive match")
	}
}
