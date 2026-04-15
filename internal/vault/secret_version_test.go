package vault

import (
	"testing"

	vaultapi "github.com/hashicorp/vault/api"
)

func TestGetSecretVersion_EmptyPath(t *testing.T) {
	c := &Client{}
	_, err := c.GetSecretVersion("", 0)
	if err == nil {
		t.Fatal("expected error for empty path, got nil")
	}
}

func TestGetSecretVersion_NoMetadata(t *testing.T) {
	mock := &mockLogical{
		readFn: func(path string) (*vaultapi.Secret, error) {
			return nil, nil
		},
	}
	c := &Client{logical: mock}
	_, err := c.GetSecretVersion("myapp/db", 0)
	if err == nil {
		t.Fatal("expected error for nil secret, got nil")
	}
}

func TestGetSecretVersion_ValidVersion(t *testing.T) {
	mock := &mockLogical{
		readFn: func(path string) (*vaultapi.Secret, error) {
			return &vaultapi.Secret{
				Data: map[string]interface{}{
					"current_version": json.Number("2"),
					"versions": map[string]interface{}{
						"1": map[string]interface{}{
							"created_time":  "2024-01-01T00:00:00Z",
							"deletion_time": "",
							"destroyed":     false,
						},
						"2": map[string]interface{}{
							"created_time":  "2024-06-01T00:00:00Z",
							"deletion_time": "",
							"destroyed":     false,
						},
					},
				},
			}, nil
		},
	}
	c := &Client{logical: mock}
	sv, err := c.GetSecretVersion("myapp/db", 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sv.Version != 1 {
		t.Errorf("expected version 1, got %d", sv.Version)
	}
	if sv.CreatedTime != "2024-01-01T00:00:00Z" {
		t.Errorf("unexpected created_time: %s", sv.CreatedTime)
	}
	if sv.Destroyed {
		t.Error("expected destroyed=false")
	}
}

func TestGetSecretVersion_MissingVersion(t *testing.T) {
	mock := &mockLogical{
		readFn: func(path string) (*vaultapi.Secret, error) {
			return &vaultapi.Secret{
				Data: map[string]interface{}{
					"current_version": json.Number("1"),
					"versions": map[string]interface{}{
						"1": map[string]interface{}{
							"created_time": "2024-01-01T00:00:00Z",
						},
					},
				},
			}, nil
		},
	}
	c := &Client{logical: mock}
	_, err := c.GetSecretVersion("myapp/db", 99)
	if err == nil {
		t.Fatal("expected error for missing version, got nil")
	}
}
