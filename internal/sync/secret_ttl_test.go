package sync

import (
	"testing"
	"time"
)

func TestCheckSecretTTL_Disabled_ReturnsNil(t *testing.T) {
	cfg := TTLConfig{Enabled: false, MaxAgeDays: 30}
	secrets := map[string]string{"OLD_KEY": "value"}
	createdAt := map[string]time.Time{
		"OLD_KEY": time.Now().AddDate(0, 0, -60),
	}
	violations := CheckSecretTTL(cfg, secrets, createdAt)
	if violations != nil {
		t.Errorf("expected nil violations when disabled, got %v", violations)
	}
}

func TestCheckSecretTTL_NoViolations(t *testing.T) {
	cfg := TTLConfig{Enabled: true, MaxAgeDays: 30}
	secrets := map[string]string{"FRESH_KEY": "value"}
	createdAt := map[string]time.Time{
		"FRESH_KEY": time.Now().AddDate(0, 0, -5),
	}
	violations := CheckSecretTTL(cfg, secrets, createdAt)
	if len(violations) != 0 {
		t.Errorf("expected no violations, got %v", violations)
	}
}

func TestCheckSecretTTL_DetectsViolation(t *testing.T) {
	cfg := TTLConfig{Enabled: true, MaxAgeDays: 30}
	secrets := map[string]string{"OLD_KEY": "value"}
	createdAt := map[string]time.Time{
		"OLD_KEY": time.Now().AddDate(0, 0, -45),
	}
	violations := CheckSecretTTL(cfg, secrets, createdAt)
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
	if violations[0].Key != "OLD_KEY" {
		t.Errorf("expected violation key OLD_KEY, got %s", violations[0].Key)
	}
	if violations[0].AgeDays < 44 {
		t.Errorf("expected AgeDays >= 44, got %d", violations[0].AgeDays)
	}
}

func TestCheckSecretTTL_MissingCreatedAt_Skipped(t *testing.T) {
	cfg := TTLConfig{Enabled: true, MaxAgeDays: 30}
	secrets := map[string]string{"NO_META": "value"}
	createdAt := map[string]time.Time{}
	violations := CheckSecretTTL(cfg, secrets, createdAt)
	if len(violations) != 0 {
		t.Errorf("expected no violations for missing metadata, got %v", violations)
	}
}
