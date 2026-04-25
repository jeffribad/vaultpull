package sync

import (
	"testing"
)

func TestAddKeyPrefix_Disabled_ReturnsOriginal(t *testing.T) {
	cfg := PrefixAddConfig{Enabled: false, Prefix: "APP_"}
	secrets := map[string]string{"DB_HOST": "localhost"}
	result := AddKeyPrefix(cfg, secrets)
	if _, ok := result["APP_DB_HOST"]; ok {
		t.Fatal("expected no prefix when disabled")
	}
	if result["DB_HOST"] != "localhost" {
		t.Fatal("expected original key to be preserved")
	}
}

func TestAddKeyPrefix_EmptyPrefix_ReturnsOriginal(t *testing.T) {
	cfg := PrefixAddConfig{Enabled: true, Prefix: ""}
	secrets := map[string]string{"FOO": "bar"}
	result := AddKeyPrefix(cfg, secrets)
	if _, ok := result["FOO"]; !ok {
		t.Fatal("expected original key when prefix is empty")
	}
}

func TestAddKeyPrefix_AllKeys(t *testing.T) {
	cfg := PrefixAddConfig{Enabled: true, Prefix: "APP_"}
	secrets := map[string]string{"HOST": "localhost", "PORT": "5432"}
	result := AddKeyPrefix(cfg, secrets)
	if result["APP_HOST"] != "localhost" {
		t.Errorf("expected APP_HOST=localhost, got %q", result["APP_HOST"])
	}
	if result["APP_PORT"] != "5432" {
		t.Errorf("expected APP_PORT=5432, got %q", result["APP_PORT"])
	}
	if len(result) != 2 {
		t.Errorf("expected 2 keys, got %d", len(result))
	}
}

func TestAddKeyPrefix_SpecificKeys_OnlyPrefixesListed(t *testing.T) {
	cfg := PrefixAddConfig{Enabled: true, Prefix: "APP_", Keys: []string{"HOST"}}
	secrets := map[string]string{"HOST": "localhost", "PORT": "5432"}
	result := AddKeyPrefix(cfg, secrets)
	if result["APP_HOST"] != "localhost" {
		t.Errorf("expected APP_HOST=localhost")
	}
	if result["PORT"] != "5432" {
		t.Errorf("expected PORT to remain unprefixed")
	}
	if _, ok := result["APP_PORT"]; ok {
		t.Error("expected PORT not to be prefixed")
	}
}

func TestAddKeyPrefix_DoesNotMutateInput(t *testing.T) {
	cfg := PrefixAddConfig{Enabled: true, Prefix: "X_"}
	secrets := map[string]string{"KEY": "val"}
	_ = AddKeyPrefix(cfg, secrets)
	if _, ok := secrets["X_KEY"]; ok {
		t.Fatal("original map was mutated")
	}
}

func TestAddKeyPrefix_CaseInsensitiveKeyMatch(t *testing.T) {
	cfg := PrefixAddConfig{Enabled: true, Prefix: "PRE_", Keys: []string{"host"}}
	secrets := map[string]string{"HOST": "localhost", "PORT": "3306"}
	result := AddKeyPrefix(cfg, secrets)
	if result["PRE_HOST"] != "localhost" {
		t.Errorf("expected case-insensitive key match to prefix HOST")
	}
	if result["PORT"] != "3306" {
		t.Errorf("expected PORT to remain unprefixed")
	}
}
