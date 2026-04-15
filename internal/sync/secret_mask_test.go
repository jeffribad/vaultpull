package sync

import (
	"strings"
	"testing"
)

func TestIsSensitiveKey_MatchesSensitiveNames(t *testing.T) {
	sensitiveKeys := []string{
		"PASSWORD", "db_password", "API_KEY", "apikey",
		"SECRET", "auth_token", "PRIVATE_KEY", "aws_secret",
		"CREDENTIAL", "passwd", "PWD",
	}
	for _, key := range sensitiveKeys {
		if !IsSensitiveKey(key) {
			t.Errorf("expected %q to be sensitive", key)
		}
	}
}

func TestIsSensitiveKey_IgnoresNonSensitiveNames(t *testing.T) {
	safeKeys := []string{
		"APP_ENV", "PORT", "HOST", "LOG_LEVEL", "REGION",
	}
	for _, key := range safeKeys {
		if IsSensitiveKey(key) {
			t.Errorf("expected %q to not be sensitive", key)
		}
	}
}

func TestMaskValue_MasksPartialValue(t *testing.T) {
	cfg := DefaultMaskConfig()
	result := MaskValue("supersecret", cfg)
	if !strings.HasPrefix(result, "supe") {
		t.Errorf("expected result to start with 'supe', got %q", result)
	}
	if !strings.Contains(result, "***") {
		t.Errorf("expected result to contain mask chars, got %q", result)
	}
	if len(result) != len("supersecret") {
		t.Errorf("expected same length as original, got %d", len(result))
	}
}

func TestMaskValue_ShortValue_FullyMasked(t *testing.T) {
	cfg := DefaultMaskConfig()
	result := MaskValue("abc", cfg)
	if result != "***" {
		t.Errorf("expected full mask for short value, got %q", result)
	}
}

func TestMaskValue_EmptyValue(t *testing.T) {
	cfg := DefaultMaskConfig()
	result := MaskValue("", cfg)
	if result != "" {
		t.Errorf("expected empty string, got %q", result)
	}
}

func TestMaskSecrets_OnlySensitiveKeysMasked(t *testing.T) {
	cfg := DefaultMaskConfig()
	secrets := map[string]string{
		"DB_PASSWORD": "hunter2",
		"APP_ENV":     "production",
		"API_KEY":     "abcdefgh",
		"PORT":        "8080",
	}
	masked := MaskSecrets(secrets, cfg)

	if masked["APP_ENV"] != "production" {
		t.Errorf("APP_ENV should not be masked, got %q", masked["APP_ENV"])
	}
	if masked["PORT"] != "8080" {
		t.Errorf("PORT should not be masked, got %q", masked["PORT"])
	}
	if masked["DB_PASSWORD"] == "hunter2" {
		t.Errorf("DB_PASSWORD should be masked")
	}
	if masked["API_KEY"] == "abcdefgh" {
		t.Errorf("API_KEY should be masked")
	}
}
