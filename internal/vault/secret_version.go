package vault

import (
	"fmt"
	"strconv"
)

// SecretVersion holds metadata about a specific version of a KV secret.
type SecretVersion struct {
	Version     int
	CreatedTime string
	DeletionTime string
	Destroyed   bool
}

// GetSecretVersion reads the version metadata for a KV v2 secret at the given path.
// If version is 0, it returns metadata for the latest version.
func (c *Client) GetSecretVersion(path string, version int) (*SecretVersion, error) {
	if path == "" {
		return nil, fmt.Errorf("secret path must not be empty")
	}

	metaPath := "secret/metadata/" + path
	secret, err := c.logical.Read(metaPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read secret metadata at %q: %w", metaPath, err)
	}
	if secret == nil || secret.Data == nil {
		return nil, fmt.Errorf("no metadata found at path %q", metaPath)
	}

	versions, ok := secret.Data["versions"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected metadata format at path %q", metaPath)
	}

	target := version
	if target == 0 {
		currentVersion, _ := secret.Data["current_version"].(json.Number)
		cv, _ := strconv.Atoi(string(currentVersion))
		target = cv
	}

	versionKey := strconv.Itoa(target)
	versionData, ok := versions[versionKey].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("version %d not found at path %q", target, path)
	}

	sv := &SecretVersion{
		Version:      target,
		CreatedTime:  stringOrEmpty(versionData, "created_time"),
		DeletionTime: stringOrEmpty(versionData, "deletion_time"),
		Destroyed:    boolOrFalse(versionData, "destroyed"),
	}
	return sv, nil
}

func stringOrEmpty(m map[string]interface{}, key string) string {
	if v, ok := m[key].(string); ok {
		return v
	}
	return ""
}

func boolOrFalse(m map[string]interface{}, key string) bool {
	if v, ok := m[key].(bool); ok {
		return v
	}
	return false
}
