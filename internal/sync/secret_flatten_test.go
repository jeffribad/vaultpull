package sync

import (
	"testing"
)

func TestFlattenSecrets_Disabled_ReturnsOriginal(t *testing.T) {
	input := map[string]string{"db.host": "localhost", "db.port": "5432"}
	cfg := FlattenConfig{Enabled: false, Separator: "__", MaxDepth: 5}
	out, err := FlattenSecrets(input, cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["db.host"] != "localhost" {
		t.Errorf("expected original key preserved, got %v", out)
	}
}

func TestFlattenSecrets_FlattensDotsToSeparator(t *testing.T) {
	input := map[string]string{"db.host": "localhost"}
	cfg := FlattenConfig{Enabled: true, Separator: "__", MaxDepth: 5}
	out, err := FlattenSecrets(input, cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["DB__HOST"] != "localhost" {
		t.Errorf("expected DB__HOST=localhost, got %v", out)
	}
}

func TestFlattenSecrets_FlattensSlashesToSeparator(t *testing.T) {
	input := map[string]string{"app/config/debug": "true"}
	cfg := FlattenConfig{Enabled: true, Separator: "_", MaxDepth: 5}
	out, err := FlattenSecrets(input, cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["APP_CONFIG_DEBUG"] != "true" {
		t.Errorf("expected APP_CONFIG_DEBUG=true, got %v", out)
	}
}

func TestFlattenSecrets_RespectsMaxDepth(t *testing.T) {
	input := map[string]string{"a.b.c.d.e.f": "deep"}
	cfg := FlattenConfig{Enabled: true, Separator: "__", MaxDepth: 3}
	out, err := FlattenSecrets(input, cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["A__B__C"] != "deep" {
		t.Errorf("expected A__B__C=deep, got %v", out)
	}
}

func TestFlattenSecrets_CollisionReturnsError(t *testing.T) {
	input := map[string]string{
		"db.host": "localhost",
		"db/host": "remotehost",
	}
	cfg := FlattenConfig{Enabled: true, Separator: "__", MaxDepth: 5}
	_, err := FlattenSecrets(input, cfg)
	if err == nil {
		t.Error("expected collision error, got nil")
	}
}

func TestFlattenSecrets_DoesNotMutateInput(t *testing.T) {
	input := map[string]string{"x.y": "val"}
	cfg := FlattenConfig{Enabled: true, Separator: "__", MaxDepth: 5}
	_, _ = FlattenSecrets(input, cfg)
	if _, ok := input["X__Y"]; ok {
		t.Error("input map was mutated")
	}
	if input["x.y"] != "val" {
		t.Error("original key was removed from input")
	}
}
