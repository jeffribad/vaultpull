package sync

import (
	"os"
	"strings"
)

// EnvProfileConfig controls environment-profile-based secret filtering.
type EnvProfileConfig struct {
	Enabled  bool
	Profile  string
	Profiles []string
}

// EnvProfileConfigFromEnv loads profile config from environment variables.
func EnvProfileConfigFromEnv() EnvProfileConfig {
	raw := strings.TrimSpace(os.Getenv("VAULTPULL_PROFILE_ENABLED"))
	enabled := raw == "true" || raw == "1"

	profile := strings.TrimSpace(os.Getenv("VAULTPULL_PROFILE"))
	allowed := splitTrimmed(os.Getenv("VAULTPULL_PROFILES"), ",")

	return EnvProfileConfig{
		Enabled:  enabled,
		Profile:  profile,
		Profiles: allowed,
	}
}

// ApplyEnvProfile filters secrets to only those matching the active profile.
// Secrets are expected to have a "_profile" metadata key or a key suffix like KEY__prod.
// If disabled or no profile set, original map is returned unchanged.
func ApplyEnvProfile(cfg EnvProfileConfig, secrets map[string]string) map[string]string {
	if !cfg.Enabled || cfg.Profile == "" {
		return secrets
	}

	profile := strings.ToLower(cfg.Profile)
	suffix := "__" + profile

	result := make(map[string]string, len(secrets))
	for k, v := range secrets {
		lower := strings.ToLower(k)
		if strings.HasSuffix(lower, suffix) {
			newKey := k[:len(k)-len(suffix)]
			result[newKey] = v
		} else if !containsProfileSuffix(lower, cfg.Profiles) {
			result[k] = v
		}
	}
	return result
}

func containsProfileSuffix(key string, profiles []string) bool {
	for _, p := range profiles {
		if strings.HasSuffix(key, "__"+strings.ToLower(p)) {
			return true
		}
	}
	return false
}
