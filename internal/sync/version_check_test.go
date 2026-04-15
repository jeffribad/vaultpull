package sync

import (
	"errors"
	"testing"
)

type mockVersionChecker struct {
	result *SecretVersionResult
	err    error
}

func (m *mockVersionChecker) GetSecretVersion(path string, version int) (*SecretVersionResult, error) {
	return m.result, m.err
}

func TestCheckVersion_EmptyPath(t *testing.T) {
	checker := &mockVersionChecker{}
	_, err := CheckVersion(checker, "")
	if err == nil {
		t.Fatal("expected error for empty path")
	}
}

func TestCheckVersion_CheckerError(t *testing.T) {
	checker := &mockVersionChecker{
		err: errors.New("vault unreachable"),
	}
	_, err := CheckVersion(checker, "myapp/config")
	if err == nil {
		t.Fatal("expected error from checker")
	}
}

func TestCheckVersion_DestroyedSecret(t *testing.T) {
	checker := &mockVersionChecker{
		result: &SecretVersionResult{
			Version:   3,
			Destroyed: true,
		},
	}
	_, err := CheckVersion(checker, "myapp/config")
	if err == nil {
		t.Fatal("expected error for destroyed secret")
	}
}

func TestCheckVersion_Valid(t *testing.T) {
	checker := &mockVersionChecker{
		result: &SecretVersionResult{
			Version:     2,
			CreatedTime: "2024-03-15T10:00:00Z",
			Destroyed:   false,
		},
	}
	info, err := CheckVersion(checker, "myapp/config")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if info.Version != 2 {
		t.Errorf("expected version 2, got %d", info.Version)
	}
	if info.Path != "myapp/config" {
		t.Errorf("unexpected path: %s", info.Path)
	}
	if info.CheckedAt.IsZero() {
		t.Error("expected CheckedAt to be set")
	}
	if info.Destroyed {
		t.Error("expected Destroyed=false")
	}
}
