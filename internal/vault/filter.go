package vault

import "strings"

// FilterByRole returns only the secrets whose keys match the allowed keys
// defined for the given role. Role definitions map a role name to a set
// of allowed secret key prefixes or exact names.
//
// If the role is not found in the definitions map, an empty map is returned.
// An empty allowedKeys slice for a role means all secrets are permitted.
func FilterByRole(secrets map[string]string, role string, roleDefs map[string][]string) map[string]string {
	allowed, exists := roleDefs[role]
	if !exists {
		return map[string]string{}
	}

	// Empty allowed list means the role has access to everything.
	if len(allowed) == 0 {
		result := make(map[string]string, len(secrets))
		for k, v := range secrets {
			result[k] = v
		}
		return result
	}

	result := make(map[string]string)
	for key, val := range secrets {
		if matchesAny(key, allowed) {
			result[key] = val
		}
	}
	return result
}

// matchesAny returns true if key equals or has a prefix matching any entry
// in the patterns slice. Patterns ending with "*" are treated as prefix matches.
func matchesAny(key string, patterns []string) bool {
	for _, p := range patterns {
		if strings.HasSuffix(p, "*") {
			if strings.HasPrefix(key, strings.TrimSuffix(p, "*")) {
				return true
			}
		} else if key == p {
			return true
		}
	}
	return false
}
