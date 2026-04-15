package vault

import "strings"

// RolePolicy defines which secret keys are accessible for a given role.
type RolePolicy struct {
	Role    string
	Allowed []string // key prefixes or exact key names
}

// FilterByRole filters a secrets map based on the role's allowed keys.
// If no policy is found for the role, all keys are returned.
func FilterByRole(secrets map[string]string, role string, policies []RolePolicy) map[string]string {
	var policy *RolePolicy
	for i := range policies {
		if strings.EqualFold(policies[i].Role, role) {
			policy = &policies[i]
			break
		}
	}

	if policy == nil {
		// No policy found — return all secrets unchanged
		result := make(map[string]string, len(secrets))
		for k, v := range secrets {
			result[k] = v
		}
		return result
	}

	filtered := make(map[string]string)
	for k, v := range secrets {
		if matchesAny(k, policy.Allowed) {
			filtered[k] = v
		}
	}
	return filtered
}

// matchesAny returns true if the key matches any of the allowed patterns.
// A pattern matches if the key equals it exactly or starts with it followed by "_".
func matchesAny(key string, patterns []string) bool {
	for _, pattern := range patterns {
		if key == pattern {
			return true
		}
		if strings.HasPrefix(key, pattern+"_") {
			return true
		}
	}
	return false
}
