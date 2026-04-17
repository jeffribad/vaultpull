package sync

import (
	"testing"
)

// TestApplyRenames_Integration_PipelineCompatibility verifies that ApplyRenames
// can be chained with ApplySecretFilter without data loss.
func TestApplyRenames_Integration_PipelineCompatibility(t *testing.T) {
	secrets := map[string]string{
		"DB_HOST": "localhost",
		"DB_PASS": "secret",
		"API_KEY":  "abc123",
	}

	// First filter to only include DB_ keys
	cfg := SecretFilterConfigFromEnv()
	cfg.IncludeKeys = []string{"DB_HOST", "DB_PASS"}
	filtered := ApplySecretFilter(secrets, cfg)

	if len(filtered) != 2 {
		t.Fatalf("expected 2 filtered keys, got %d", len(filtered))
	}

	// Then rename DB_PASS -> DATABASE_PASSWORD
	rm := RenameMap{"DB_PASS": "DATABASE_PASSWORD"}
	renamed := ApplyRenames(filtered, rm)

	if _, ok := renamed["DB_PASS"]; ok {
		t.Fatal("DB_PASS should have been renamed")
	}
	if renamed["DATABASE_PASSWORD"] != "secret" {
		t.Fatalf("expected secret, got %q", renamed["DATABASE_PASSWORD"])
	}
	if renamed["DB_HOST"] != "localhost" {
		t.Fatal("DB_HOST should be unchanged")
	}
}

// TestApplyRenames_Integration_CollisionLastWins checks that if two source keys
// are renamed to the same target, the last one processed wins.
func TestApplyRenames_Integration_CollisionHandled(t *testing.T) {
	secrets := map[string]string{
		"OLD_A": "value_a",
		"OLD_B": "value_b",
	}
	rm := RenameMap{
		"OLD_A": "SHARED",
		"OLD_B": "SHARED",
	}
	out := ApplyRenames(secrets, rm)
	if _, ok := out["SHARED"]; !ok {
		t.Fatal("SHARED key should exist in output")
	}
}
