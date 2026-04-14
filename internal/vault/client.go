package vault

import (
	"errors"
	"fmt"

	vaultapi "github.com/hashicorp/vault/api"
)

// Client wraps the HashiCorp Vault API client.
type Client struct {
	logical *vaultapi.Logical
}

// NewClient creates a new Vault client using the provided address and token.
func NewClient(address, token string) (*Client, error) {
	if address == "" {
		return nil, errors.New("vault address must not be empty")
	}
	if token == "" {
		return nil, errors.New("vault token must not be empty")
	}

	cfg := vaultapi.DefaultConfig()
	cfg.Address = address

	raw, err := vaultapi.NewClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("creating vault api client: %w", err)
	}

	raw.SetToken(token)

	return &Client{logical: raw.Logical()}, nil
}

// ReadSecrets reads key/value secrets from the given KV path.
// Supports KV v2 by checking for a nested "data" map.
func (c *Client) ReadSecrets(path string) (map[string]string, error) {
	if path == "" {
		return nil, errors.New("secret path must not be empty")
	}

	secret, err := c.logical.Read(path)
	if err != nil {
		return nil, fmt.Errorf("reading secrets from %q: %w", path, err)
	}
	if secret == nil {
		return nil, fmt.Errorf("no secret found at path %q", path)
	}

	rawData := secret.Data

	// KV v2 wraps values inside a "data" key.
	if nested, ok := rawData["data"]; ok {
		if nestedMap, ok := nested.(map[string]interface{}); ok {
			rawData = nestedMap
		}
	}

	result := make(map[string]string, len(rawData))
	for k, v := range rawData {
		switch val := v.(type) {
		case string:
			result[k] = val
		default:
			result[k] = fmt.Sprintf("%v", v)
		}
	}

	return result, nil
}
