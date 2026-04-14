package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoad_FromEnv(t *testing.T) {
	t.Setenv("VAULT_ADDR", "http://vault.example.com:8200")
	t.Setenv("VAULT_TOKEN", "s.testtoken")
	t.Setenv("VAULTPULL_SECRET_PATH", "secret/data/myapp")

	cfg, err := Load("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.VaultAddr != "http://vault.example.com:8200" {
		t.Errorf("expected vault addr from env, got %q", cfg.VaultAddr)
	}
	if cfg.VaultToken != "s.testtoken" {
		t.Errorf("expected vault token from env, got %q", cfg.VaultToken)
	}
	if cfg.SecretPath != "secret/data/myapp" {
		t.Errorf("expected secret path from env, got %q", cfg.SecretPath)
	}
	if cfg.OutputFile != ".env" {
		t.Errorf("expected default output file '.env', got %q", cfg.OutputFile)
	}
}

func TestLoad_FromFile(t *testing.T) {
	dir := t.TempDir()
	cfgPath := filepath.Join(dir, "config.yaml")

	content := []byte(`
vault_addr: "http://localhost:8200"
vault_token: "s.filetoken"
secret_path: "secret/data/svc"
output_file: "secrets.env"
roles:
  - admin
  - readonly
`)
	if err := os.WriteFile(cfgPath, content, 0600); err != nil {
		t.Fatalf("writing temp config: %v", err)
	}

	cfg, err := Load(cfgPath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.SecretPath != "secret/data/svc" {
		t.Errorf("expected secret_path 'secret/data/svc', got %q", cfg.SecretPath)
	}
	if cfg.OutputFile != "secrets.env" {
		t.Errorf("expected output_file 'secrets.env', got %q", cfg.OutputFile)
	}
	if len(cfg.Roles) != 2 || cfg.Roles[0] != "admin" {
		t.Errorf("unexpected roles: %v", cfg.Roles)
	}
}

func TestLoad_MissingToken(t *testing.T) {
	os.Unsetenv("VAULT_TOKEN")
	os.Unsetenv("VAULTPULL_VAULT_TOKEN")
	t.Setenv("VAULTPULL_SECRET_PATH", "secret/data/test")

	_, err := Load("")
	if err == nil {
		t.Fatal("expected validation error for missing vault_token")
	}
}

func TestLoad_MissingSecretPath(t *testing.T) {
	t.Setenv("VAULT_TOKEN", "s.tok")
	os.Unsetenv("VAULTPULL_SECRET_PATH")

	_, err := Load("")
	if err == nil {
		t.Fatal("expected validation error for missing secret_path")
	}
}
