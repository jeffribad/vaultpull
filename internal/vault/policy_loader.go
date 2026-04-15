package vault

import (
	"encoding/json"
	"fmt"
	"os"
)

// policyFile represents the JSON structure for role policy configuration.
type policyFile struct {
	Policies []RolePolicy `json:"policies"`
}

// LoadPolicies reads role policies from a JSON file at the given path.
// Returns an empty slice if the file does not exist (non-strict mode).
func LoadPolicies(path string) ([]RolePolicy, error) {
	if path == "" {
		return []RolePolicy{}, nil
	}

	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return []RolePolicy{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("reading policy file %q: %w", path, err)
	}

	var pf policyFile
	if err := json.Unmarshal(data, &pf); err != nil {
		return nil, fmt.Errorf("parsing policy file %q: %w", path, err)
	}

	return pf.Policies, nil
}

// DefaultPoliciesFromEnv builds a minimal policy list from the
// VAULTPULL_ROLE and VAULTPULL_ALLOWED_KEYS environment variables.
// VAULTPULL_ALLOWED_KEYS should be a comma-separated list of key prefixes.
func DefaultPoliciesFromEnv() []RolePolicy {
	role := os.Getenv("VAULTPULL_ROLE")
	if role == "" {
		return nil
	}

	allowedRaw := os.Getenv("VAULTPULL_ALLOWED_KEYS")
	if allowedRaw == "" {
		return nil
	}

	keys := splitAndTrim(allowedRaw, ",")
	return []RolePolicy{{Role: role, Allowed: keys}}
}

func splitAndTrim(s, sep string) []string {
	var result []string
	for _, part := range splitString(s, sep) {
		trimmed := trimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

func splitString(s, sep string) []string {
	var parts []string
	start := 0
	for i := 0; i <= len(s)-len(sep); i++ {
		if s[i:i+len(sep)] == sep {
			parts = append(parts, s[start:i])
			start = i + len(sep)
		}
	}
	parts = append(parts, s[start:])
	return parts
}

func trimSpace(s string) string {
	start, end := 0, len(s)
	for start < end && (s[start] == ' ' || s[start] == '\t') {
		start++
	}
	for end > start && (s[end-1] == ' ' || s[end-1] == '\t') {
		end--
	}
	return s[start:end]
}
