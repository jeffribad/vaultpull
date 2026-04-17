package sync

import (
	"strings"
	"testing"
	"time"
)

func TestCheckSecretExpiry_Disabled_ReturnsNil(t *testing.T) {
	cfg := ExpiryConfig{Enabled: false}
	secrets := map[string]string{
		"EXPIRES_AT": time.Now().Add(-time.Hour).Format(time.RFC3339),
	}
	warnings := CheckSecretExpiry(cfg, secrets)
	if len(warnings) != 0 {
		t.Errorf("expected no warnings when disabled, got %d", len(warnings))
	}
}

func TestCheckSecretExpiry_NoViolations(t *testing.T) {
	cfg := ExpiryConfig{Enabled: true, WarnWithinDays: 7}
	secrets := map[string]string{
		"EXPIRES_AT": time.Now().Add(30 * 24 * time.Hour).Format(time.RFC3339),
	}
	warnings := CheckSecretExpiry(cfg, secrets)
	if len(warnings) != 0 {
		t.Errorf("expected no warnings, got %d", len(warnings))
	}
}

func TestCheckSecretExpiry_DetectsExpiringSoon(t *testing.T) {
	cfg := ExpiryConfig{Enabled: true, WarnWithinDays: 7}
	secrets := map[string]string{
		"TOKEN_EXPIRES": time.Now().Add(2 * 24 * time.Hour).Format(time.RFC3339),
	}
	warnings := CheckSecretExpiry(cfg, secrets)
	if len(warnings) != 1 {
		t.Fatalf("expected 1 warning, got %d", len(warnings))
	}
	if warnings[0].Expired {
		t.Error("expected not expired, only expiring soon")
	}
}

func TestCheckSecretExpiry_DetectsExpired(t *testing.T) {
	cfg := ExpiryConfig{Enabled: true, WarnWithinDays: 7}
	secrets := map[string]string{
		"OLD_TOKEN": time.Now().Add(-24 * time.Hour).Format(time.RFC3339),
	}
	warnings := CheckSecretExpiry(cfg, secrets)
	if len(warnings) != 1 {
		t.Fatalf("expected 1 warning, got %d", len(warnings))
	}
	if !warnings[0].Expired {
		t.Error("expected Expired=true")
	}
}

func TestCheckSecretExpiry_SkipsNonDateValues(t *testing.T) {
	cfg := ExpiryConfig{Enabled: true, WarnWithinDays: 7}
	secrets := map[string]string{
		"DB_PASSWORD": "supersecret",
		"API_KEY":     "abc123",
	}
	warnings := CheckSecretExpiry(cfg, secrets)
	if len(warnings) != 0 {
		t.Errorf("expected no warnings for non-date values, got %d", len(warnings))
	}
}

func TestExpiryWarning_String_Expired(t *testing.T) {
	w := ExpiryWarning{Key: "MY_TOKEN", ExpiresAt: time.Now().Add(-time.Hour), Expired: true}
	s := w.String()
	if !strings.Contains(s, "expired") {
		t.Errorf("expected 'expired' in string, got: %s", s)
	}
}
