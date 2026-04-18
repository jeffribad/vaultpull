package sync

import (
	"testing"
)

func TestRedactSecrets_Disabled_ReturnsOriginal(t *testing.T) {
	secrets := map[string]string{"API_KEY": "secret123", "HOST": "localhost"}
	cfg := RedactConfig{Enabled: false, Keys: []string{"API_KEY"}}

	result := RedactSecrets(secrets, cfg)
	if result["API_KEY"] != "secret123" {
		t.Errorf("expected original value, got %s", result["API_KEY"])
	}
}

func TestRedactSecrets_RedactsMatchingKey(t *testing.T) {
	secrets := map[string]string{"API_KEY": "secret123", "HOST": "localhost"}
	cfg := RedactConfig{Enabled: true, Keys: []string{"API_KEY"}}

	result := RedactSecrets(secrets, cfg)
	if result["API_KEY"] != "[REDACTED]" {
		t.Errorf("expected [REDACTED], got %s", result["API_KEY"])
	}
	if result["HOST"] != "localhost" {
		t.Errorf("expected HOST to be unchanged, got %s", result["HOST"])
	}
}

func TestRedactSecrets_CaseInsensitiveMatch(t *testing.T) {
	secrets := map[string]string{"api_key": "topsecret"}
	cfg := RedactConfig{Enabled: true, Keys: []string{"API_KEY"}}

	result := RedactSecrets(secrets, cfg)
	if result["api_key"] != "[REDACTED]" {
		t.Errorf("expected case-insensitive redaction, got %s", result["api_key"])
	}
}

func TestRedactSecrets_NoKeys_ReturnsOriginal(t *testing.T) {
	secrets := map[string]string{"TOKEN": "abc"}
	cfg := RedactConfig{Enabled: true, Keys: []string{}}

	result := RedactSecrets(secrets, cfg)
	if result["TOKEN"] != "abc" {
		t.Errorf("expected original value when no keys configured, got %s", result["TOKEN"])
	}
}

func TestRedactSecrets_DoesNotMutateInput(t *testing.T) {
	secrets := map[string]string{"SECRET": "value"}
	cfg := RedactConfig{Enabled: true, Keys: []string{"SECRET"}}

	RedactSecrets(secrets, cfg)
	if secrets["SECRET"] != "value" {
		t.Error("expected original map to be unmodified")
	}
}
