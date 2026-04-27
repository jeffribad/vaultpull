package sync

import (
	"os"
	"strings"
)

// PivotConfig controls key pivoting — remapping secrets based on a runtime
// environment variable that selects which key variant to use.
type PivotConfig struct {
	Enabled  bool
	EnvVar   string // OS env var whose value selects the active variant suffix
	Suffixes []string // known suffixes to strip from candidates
}

// PivotConfigFromEnv loads PivotConfig from environment variables.
//
//	VAULTPULL_PIVOT_ENABLED=true
//	VAULTPULL_PIVOT_ENV_VAR=APP_ENV          # e.g. "production"
//	VAULTPULL_PIVOT_SUFFIXES=dev,staging,production
func PivotConfigFromEnv() PivotConfig {
	enabled := strings.TrimSpace(os.Getenv("VAULTPULL_PIVOT_ENABLED"))
	envVar := strings.TrimSpace(os.Getenv("VAULTPULL_PIVOT_ENV_VAR"))
	suffixesRaw := strings.TrimSpace(os.Getenv("VAULTPULL_PIVOT_SUFFIXES"))

	var suffixes []string
	for _, s := range strings.Split(suffixesRaw, ",") {
		if t := strings.TrimSpace(s); t != "" {
			suffixes = append(suffixes, t)
		}
	}

	return PivotConfig{
		Enabled:  enabled == "true" || enabled == "1",
		EnvVar:   envVar,
		Suffixes: suffixes,
	}
}

// ApplyPivot selects the active variant of each key based on the pivot suffix
// resolved from the configured environment variable.
//
// Given secrets {DB_URL_production: "x", DB_URL_dev: "y"} and pivot suffix
// "production", the result will contain {DB_URL: "x"}, dropping other variants.
// Keys that do not match any known suffix are passed through unchanged.
func ApplyPivot(cfg PivotConfig, secrets map[string]string) map[string]string {
	if !cfg.Enabled || cfg.EnvVar == "" || len(cfg.Suffixes) == 0 {
		return secrets
	}

	active := strings.TrimSpace(os.Getenv(cfg.EnvVar))
	if active == "" {
		return secrets
	}

	result := make(map[string]string, len(secrets))

	for k, v := range secrets {
		matched := false
		for _, sfx := range cfg.Suffixes {
			sep := "_" + sfx
			if strings.HasSuffix(strings.ToLower(k), strings.ToLower(sep)) {
				matched = true
				if strings.EqualFold(sfx, active) {
					base := k[:len(k)-len(sep)]
					result[base] = v
				}
				break
			}
		}
		if !matched {
			result[k] = v
		}
	}

	return result
}
