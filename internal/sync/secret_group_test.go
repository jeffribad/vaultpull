package sync

import (
	"testing"
)

func TestGroupSecrets_NoLabels(t *testing.T) {
	secrets := map[string]string{"DB_URL": "postgres://", "API_KEY": "abc"}
	groups := GroupSecrets(secrets, nil, "group")
	if len(groups["default"]) != 2 {
		t.Errorf("expected 2 secrets in default group, got %d", len(groups["default"]))
	}
}

func TestGroupSecrets_SingleGroup(t *testing.T) {
	secrets := map[string]string{"DB_URL": "postgres://"}
	labels := map[string]map[string]string{
		"DB_URL": {"group": "backend"},
	}
	groups := GroupSecrets(secrets, labels, "group")
	if _, ok := groups["backend"]; !ok {
		t.Error("expected backend group")
	}
	if groups["backend"]["DB_URL"] != "postgres://" {
		t.Error("expected DB_URL in backend group")
	}
}

func TestGroupSecrets_MultipleGroups(t *testing.T) {
	secrets := map[string]string{"DB_URL": "pg", "REDIS_URL": "redis", "STRIPE_KEY": "sk"}
	labels := map[string]map[string]string{
		"DB_URL":    {"group": "backend"},
		"REDIS_URL": {"group": "backend"},
		"STRIPE_KEY": {"group": "payments"},
	}
	groups := GroupSecrets(secrets, labels, "group")
	if len(groups["backend"]) != 2 {
		t.Errorf("expected 2 in backend, got %d", len(groups["backend"]))
	}
	if len(groups["payments"]) != 1 {
		t.Errorf("expected 1 in payments, got %d", len(groups["payments"]))
	}
}

func TestGroupSecrets_MixedLabeled(t *testing.T) {
	secrets := map[string]string{"A": "1", "B": "2"}
	labels := map[string]map[string]string{
		"A": {"group": "alpha"},
	}
	groups := GroupSecrets(secrets, labels, "group")
	if groups["alpha"]["A"] != "1" {
		t.Error("expected A in alpha")
	}
	if groups["default"]["B"] != "2" {
		t.Error("expected B in default")
	}
}

func TestGroupSecrets_EmptyGroupLabel_FallsToDefault(t *testing.T) {
	secrets := map[string]string{"X": "val"}
	labels := map[string]map[string]string{
		"X": {"group": "   "},
	}
	groups := GroupSecrets(secrets, labels, "group")
	if groups["default"]["X"] != "val" {
		t.Error("expected X in default when label is blank")
	}
}
