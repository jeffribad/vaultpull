package sync

import (
	"testing"
)

func TestApplyPromotion_Disabled_ReturnsOriginal(t *testing.T) {
	secrets := map[string]string{"STAGING_DB": "val"}
	cfg := PromoteConfig{Enabled: false, FromPrefix: "STAGING_", ToPrefix: "PROD_"}
	out := ApplyPromotion(secrets, cfg)
	if _, ok := out["PROD_DB"]; ok {
		t.Error("expected no promotion when disabled")
	}
}

func TestApplyPromotion_EmptyFromPrefix_ReturnsOriginal(t *testing.T) {
	secrets := map[string]string{"STAGING_DB": "val"}
	cfg := PromoteConfig{Enabled: true, FromPrefix: "", ToPrefix: "PROD_"}
	out := ApplyPromotion(secrets, cfg)
	if _, ok := out["PROD_DB"]; ok {
		t.Error("expected no promotion with empty from prefix")
	}
}

func TestApplyPromotion_CopiesMatchingKeys(t *testing.T) {
	secrets := map[string]string{"STAGING_DB": "postgres", "OTHER": "x"}
	cfg := PromoteConfig{Enabled: true, FromPrefix: "STAGING_", ToPrefix: "PROD_", Overwrite: true}
	out := ApplyPromotion(secrets, cfg)
	if out["PROD_DB"] != "postgres" {
		t.Errorf("expected PROD_DB=postgres, got %q", out["PROD_DB"])
	}
	if out["STAGING_DB"] != "postgres" {
		t.Error("expected original key to be preserved")
	}
	if out["OTHER"] != "x" {
		t.Error("expected non-matching key to be preserved")
	}
}

func TestApplyPromotion_NoOverwrite_KeepsExisting(t *testing.T) {
	secrets := map[string]string{"STAGING_DB": "new", "PROD_DB": "existing"}
	cfg := PromoteConfig{Enabled: true, FromPrefix: "STAGING_", ToPrefix: "PROD_", Overwrite: false}
	out := ApplyPromotion(secrets, cfg)
	if out["PROD_DB"] != "existing" {
		t.Errorf("expected PROD_DB=existing, got %q", out["PROD_DB"])
	}
}

func TestApplyPromotion_DoesNotMutateInput(t *testing.T) {
	secrets := map[string]string{"STAGING_KEY": "val"}
	cfg := PromoteConfig{Enabled: true, FromPrefix: "STAGING_", ToPrefix: "PROD_", Overwrite: true}
	_ = ApplyPromotion(secrets, cfg)
	if _, ok := secrets["PROD_KEY"]; ok {
		t.Error("input map should not be mutated")
	}
}
