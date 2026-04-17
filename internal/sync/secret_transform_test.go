package sync

import (
	"testing"
)

func TestApplyTransforms_NoOp(t *testing.T) {
	secrets := map[string]string{"db_host": "localhost"}
	out, err := ApplyTransforms(secrets, TransformConfig{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["db_host"] != "localhost" {
		t.Errorf("expected key unchanged, got %v", out)
	}
}

func TestApplyTransforms_AddsPrefix(t *testing.T) {
	secrets := map[string]string{"host": "localhost"}
	out, err := ApplyTransforms(secrets, TransformConfig{Prefix: "APP_"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := out["APP_host"]; !ok {
		t.Errorf("expected key APP_host, got %v", out)
	}
}

func TestApplyTransforms_AddsSuffix(t *testing.T) {
	secrets := map[string]string{"host": "localhost"}
	out, err := ApplyTransforms(secrets, TransformConfig{Suffix: "_VAR"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := out["host_VAR"]; !ok {
		t.Errorf("expected key host_VAR, got %v", out)
	}
}

func TestApplyTransforms_UpperCase(t *testing.T) {
	secrets := map[string]string{"db_host": "localhost"}
	out, err := ApplyTransforms(secrets, TransformConfig{UpperCase: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := out["DB_HOST"]; !ok {
		t.Errorf("expected key DB_HOST, got %v", out)
	}
}

func TestApplyTransforms_LowerCase(t *testing.T) {
	secrets := map[string]string{"DB_HOST": "localhost"}
	out, err := ApplyTransforms(secrets, TransformConfig{LowerCase: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := out["db_host"]; !ok {
		t.Errorf("expected key db_host, got %v", out)
	}
}

func TestApplyTransforms_ConflictingCaseFlags(t *testing.T) {
	secrets := map[string]string{"key": "val"}
	_, err := ApplyTransforms(secrets, TransformConfig{UpperCase: true, LowerCase: true})
	if err == nil {
		t.Error("expected error for conflicting case flags")
	}
}
