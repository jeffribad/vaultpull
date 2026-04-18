package sync

import "testing"

func TestCastSecrets_Disabled_ReturnsOriginal(t *testing.T) {
	secrets := map[string]string{"FEATURE_FLAG": "yes", "PORT": "8080abc"}
	cfg := CastConfig{Enabled: false}
	out := CastSecrets(secrets, cfg)
	if out["FEATURE_FLAG"] != "yes" {
		t.Errorf("expected original value, got %q", out["FEATURE_FLAG"])
	}
}

func TestCastSecrets_BoolTrue(t *testing.T) {
	cases := []string{"true", "1", "yes", "TRUE", "YES"}
	for _, v := range cases {
		secrets := map[string]string{"FLAG": v}
		cfg := CastConfig{Enabled: true, BoolKeys: []string{"FLAG"}}
		out := CastSecrets(secrets, cfg)
		if out["FLAG"] != "true" {
			t.Errorf("value %q: expected \"true\", got %q", v, out["FLAG"])
		}
	}
}

func TestCastSecrets_BoolFalse(t *testing.T) {
	cases := []string{"false", "0", "no", "off"}
	for _, v := range cases {
		secrets := map[string]string{"FLAG": v}
		cfg := CastConfig{Enabled: true, BoolKeys: []string{"FLAG"}}
		out := CastSecrets(secrets, cfg)
		if out["FLAG"] != "false" {
			t.Errorf("value %q: expected \"false\", got %q", v, out["FLAG"])
		}
	}
}

func TestCastSecrets_IntStripsNonNumeric(t *testing.T) {
	secrets := map[string]string{"PORT": "8080abc"}
	cfg := CastConfig{Enabled: true, IntKeys: []string{"PORT"}}
	out := CastSecrets(secrets, cfg)
	if out["PORT"] != "8080" {
		t.Errorf("expected \"8080\", got %q", out["PORT"])
	}
}

func TestCastSecrets_IntNoDigits_FallsBackToZero(t *testing.T) {
	secrets := map[string]string{"PORT": "none"}
	cfg := CastConfig{Enabled: true, IntKeys: []string{"PORT"}}
	out := CastSecrets(secrets, cfg)
	if out["PORT"] != "0" {
		t.Errorf("expected \"0\", got %q", out["PORT"])
	}
}

func TestCastSecrets_IntStripsLeadingZeros(t *testing.T) {
	secrets := map[string]string{"PORT": "00080"}
	cfg := CastConfig{Enabled: true, IntKeys: []string{"PORT"}}
	out := CastSecrets(secrets, cfg)
	if out["PORT"] != "80" {
		t.Errorf("expected \"80\", got %q", out["PORT"])
	}
}

func TestCastSecrets_DoesNotMutateInput(t *testing.T) {
	secrets := map[string]string{"FLAG": "yes"}
	cfg := CastConfig{Enabled: true, BoolKeys: []string{"FLAG"}}
	CastSecrets(secrets, cfg)
	if secrets["FLAG"] != "yes" {
		t.Error("input map was mutated")
	}
}
