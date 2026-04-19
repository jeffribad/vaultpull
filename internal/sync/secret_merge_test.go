package sync

import (
	"testing"
)

func TestMergeSecrets_Disabled_ReturnsVaultOnly(t *testing.T) {
	cfg := MergeConfig{Enabled: false}
	vault := map[string]string{"A": "vault-a"}
	base := map[string]string{"A": "local-a", "B": "local-b"}
	out := MergeSecrets(cfg, vault, base)
	if out["A"] != "vault-a" {
		t.Errorf("expected vault-a, got %s", out["A"])
	}
	if _, ok := out["B"]; ok {
		t.Error("expected B to be absent when disabled")
	}
}

func TestMergeSecrets_VaultWins(t *testing.T) {
	cfg := MergeConfig{Enabled: true, Strategy: "vault-wins"}
	vault := map[string]string{"A": "vault-a", "C": "vault-c"}
	base := map[string]string{"A": "local-a", "B": "local-b"}
	out := MergeSecrets(cfg, vault, base)
	if out["A"] != "vault-a" {
		t.Errorf("expected vault-a, got %s", out["A"])
	}
	if out["B"] != "local-b" {
		t.Errorf("expected local-b, got %s", out["B"])
	}
	if out["C"] != "vault-c" {
		t.Errorf("expected vault-c, got %s", out["C"])
	}
}

func TestMergeSecrets_LocalWins(t *testing.T) {
	cfg := MergeConfig{Enabled: true, Strategy: "local-wins"}
	vault := map[string]string{"A": "vault-a", "C": "vault-c"}
	base := map[string]string{"A": "local-a", "B": "local-b"}
	out := MergeSecrets(cfg, vault, base)
	if out["A"] != "local-a" {
		t.Errorf("expected local-a, got %s", out["A"])
	}
	if out["B"] != "local-b" {
		t.Errorf("expected local-b, got %s", out["B"])
	}
	if out["C"] != "vault-c" {
		t.Errorf("expected vault-c, got %s", out["C"])
	}
}

func TestMergeSecrets_Union(t *testing.T) {
	cfg := MergeConfig{Enabled: true, Strategy: "union"}
	vault := map[string]string{"A": "vault-a", "C": "vault-c"}
	base := map[string]string{"A": "local-a", "B": "local-b"}
	out := MergeSecrets(cfg, vault, base)
	if out["A"] != "local-a" {
		t.Errorf("union: base should win for existing key, got %s", out["A"])
	}
	if out["B"] != "local-b" {
		t.Errorf("expected local-b, got %s", out["B"])
	}
	if out["C"] != "vault-c" {
		t.Errorf("expected vault-c for vault-only key, got %s", out["C"])
	}
}

func TestMergeSecrets_DoesNotMutateInput(t *testing.T) {
	cfg := MergeConfig{Enabled: true, Strategy: "vault-wins"}
	vault := map[string]string{"A": "vault-a"}
	base := map[string]string{"B": "local-b"}
	MergeSecrets(cfg, vault, base)
	if len(vault) != 1 || len(base) != 1 {
		t.Error("input maps were mutated")
	}
}
