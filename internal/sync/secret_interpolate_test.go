package sync

import (
	"testing"
)

func TestInterpolateSecrets_Disabled_ReturnsOriginal(t *testing.T) {
	secrets := map[string]string{"A": "${B}", "B": "hello"}
	cfg := InterpolateConfig{Enabled: false}
	out, err := InterpolateSecrets(secrets, cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["A"] != "${B}" {
		t.Errorf("expected original value, got %q", out["A"])
	}
}

func TestInterpolateSecrets_ResolvesReference(t *testing.T) {
	secrets := map[string]string{"BASE_URL": "https://example.com", "API_URL": "${BASE_URL}/api"}
	cfg := InterpolateConfig{Enabled: true}
	out, err := InterpolateSecrets(secrets, cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["API_URL"] != "https://example.com/api" {
		t.Errorf("expected interpolated value, got %q", out["API_URL"])
	}
}

func TestInterpolateSecrets_UnresolvedRef_ReturnsError(t *testing.T) {
	secrets := map[string]string{"A": "${MISSING}"}
	cfg := InterpolateConfig{Enabled: true, AllowEnv: false}
	_, err := InterpolateSecrets(secrets, cfg)
	if err == nil {
		t.Fatal("expected error for unresolved reference")
	}
}

func TestInterpolateSecrets_AllowEnv_ResolvesFromOS(t *testing.T) {
	t.Setenv("MY_HOST", "localhost")
	secrets := map[string]string{"DSN": "postgres://${MY_HOST}/db"}
	cfg := InterpolateConfig{Enabled: true, AllowEnv: true}
	out, err := InterpolateSecrets(secrets, cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["DSN"] != "postgres://localhost/db" {
		t.Errorf("got %q", out["DSN"])
	}
}

func TestInterpolateSecrets_DoesNotMutateInput(t *testing.T) {
	original := map[string]string{"A": "val", "B": "${A}"}
	cfg := InterpolateConfig{Enabled: true}
	_, _ = InterpolateSecrets(original, cfg)
	if original["B"] != "${A}" {
		t.Error("input map was mutated")
	}
}

func TestInterpolateSecrets_NoReferences_Unchanged(t *testing.T) {
	secrets := map[string]string{"FOO": "bar", "BAZ": "qux"}
	cfg := InterpolateConfig{Enabled: true}
	out, err := InterpolateSecrets(secrets, cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["FOO"] != "bar" || out["BAZ"] != "qux" {
		t.Error("values changed unexpectedly")
	}
}
