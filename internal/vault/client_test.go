package vault_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/user/vaultpull/internal/vault"
)

func TestNewClient_MissingAddress(t *testing.T) {
	_, err := vault.NewClient("", "token")
	if err == nil {
		t.Fatal("expected error for empty address, got nil")
	}
}

func TestNewClient_MissingToken(t *testing.T) {
	_, err := vault.NewClient("http://127.0.0.1:8200", "")
	if err == nil {
		t.Fatal("expected error for empty token, got nil")
	}
}

func TestNewClient_Valid(t *testing.T) {
	c, err := vault.NewClient("http://127.0.0.1:8200", "test-token")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c == nil {
		t.Fatal("expected non-nil client")
	}
}

func TestReadSecrets_EmptyPath(t *testing.T) {
	c, _ := vault.NewClient("http://127.0.0.1:8200", "test-token")
	_, err := c.ReadSecrets("")
	if err == nil {
		t.Fatal("expected error for empty path, got nil")
	}
}

func TestReadSecrets_KVv2Response(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"data": {
				"data": {
					"DB_PASSWORD": "s3cr3t",
					"API_KEY": "abc123"
				},
				"metadata": {}
			}
		}`))
	}))
	defer server.Close()

	c, err := vault.NewClient(server.URL, "test-token")
	if err != nil {
		t.Fatalf("unexpected error creating client: %v", err)
	}

	secrets, err := c.ReadSecrets("secret/data/myapp")
	if err != nil {
		t.Fatalf("unexpected error reading secrets: %v", err)
	}

	if secrets["DB_PASSWORD"] != "s3cr3t" {
		t.Errorf("expected DB_PASSWORD=s3cr3t, got %q", secrets["DB_PASSWORD"])
	}
	if secrets["API_KEY"] != "abc123" {
		t.Errorf("expected API_KEY=abc123, got %q", secrets["API_KEY"])
	}
}
