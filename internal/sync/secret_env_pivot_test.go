package sync

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestApplyPivot_Disabled_ReturnsOriginal(t *testing.T) {
	cfg := PivotConfig{Enabled: false}
	input := map[string]string{"DB_URL_production": "prod-url"}
	out := ApplyPivot(cfg, input)
	assert.Equal(t, input, out)
}

func TestApplyPivot_NoEnvVar_ReturnsOriginal(t *testing.T) {
	cfg := PivotConfig{Enabled: true, EnvVar: "", Suffixes: []string{"dev", "production"}}
	input := map[string]string{"DB_URL_production": "prod-url"}
	out := ApplyPivot(cfg, input)
	assert.Equal(t, input, out)
}

func TestApplyPivot_SelectsActiveSuffix(t *testing.T) {
	t.Setenv("APP_ENV", "production")
	cfg := PivotConfig{
		Enabled:  true,
		EnvVar:   "APP_ENV",
		Suffixes: []string{"dev", "staging", "production"},
	}
	input := map[string]string{
		"DB_URL_production": "prod-url",
		"DB_URL_dev":        "dev-url",
		"DB_URL_staging":    "staging-url",
	}
	out := ApplyPivot(cfg, input)
	assert.Equal(t, map[string]string{"DB_URL": "prod-url"}, out)
}

func TestApplyPivot_PassesThroughUnmatchedKeys(t *testing.T) {
	t.Setenv("APP_ENV", "dev")
	cfg := PivotConfig{
		Enabled:  true,
		EnvVar:   "APP_ENV",
		Suffixes: []string{"dev", "production"},
	}
	input := map[string]string{
		"API_KEY":       "abc123",
		"DB_URL_dev":    "dev-url",
		"DB_URL_production": "prod-url",
	}
	out := ApplyPivot(cfg, input)
	assert.Equal(t, "dev-url", out["DB_URL"])
	assert.Equal(t, "abc123", out["API_KEY"])
	assert.NotContains(t, out, "DB_URL_production")
}

func TestApplyPivot_DoesNotMutateInput(t *testing.T) {
	t.Setenv("APP_ENV", "staging")
	cfg := PivotConfig{
		Enabled:  true,
		EnvVar:   "APP_ENV",
		Suffixes: []string{"dev", "staging"},
	}
	input := map[string]string{
		"KEY_staging": "s",
		"KEY_dev":     "d",
	}
	copy := map[string]string{"KEY_staging": "s", "KEY_dev": "d"}
	ApplyPivot(cfg, input)
	assert.Equal(t, copy, input)
}

func TestPivotConfigFromEnv_Defaults(t *testing.T) {
	t.Setenv("VAULTPULL_PIVOT_ENABLED", "")
	t.Setenv("VAULTPULL_PIVOT_ENV_VAR", "")
	t.Setenv("VAULTPULL_PIVOT_SUFFIXES", "")
	cfg := PivotConfigFromEnv()
	assert.False(t, cfg.Enabled)
	assert.Empty(t, cfg.EnvVar)
	assert.Empty(t, cfg.Suffixes)
}

func TestPivotConfigFromEnv_Enabled(t *testing.T) {
	t.Setenv("VAULTPULL_PIVOT_ENABLED", "true")
	t.Setenv("VAULTPULL_PIVOT_ENV_VAR", "DEPLOY_ENV")
	t.Setenv("VAULTPULL_PIVOT_SUFFIXES", "dev, staging, production")
	cfg := PivotConfigFromEnv()
	assert.True(t, cfg.Enabled)
	assert.Equal(t, "DEPLOY_ENV", cfg.EnvVar)
	assert.Equal(t, []string{"dev", "staging", "production"}, cfg.Suffixes)
}
