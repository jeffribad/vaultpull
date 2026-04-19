package sync

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRequiredConfigFromEnv_Defaults(t *testing.T) {
	t.Setenv("VAULTPULL_REQUIRED_KEYS", "")
	t.Setenv("VAULTPULL_REQUIRED_FAIL_FAST", "")

	cfg := RequiredConfigFromEnv()
	assert.False(t, cfg.Enabled)
	assert.Empty(t, cfg.Keys)
	assert.False(t, cfg.FailFast)
}

func TestRequiredConfigFromEnv_ParsesKeys(t *testing.T) {
	t.Setenv("VAULTPULL_REQUIRED_KEYS", "DB_HOST, DB_PORT , API_KEY")
	t.Setenv("VAULTPULL_REQUIRED_FAIL_FAST", "")

	cfg := RequiredConfigFromEnv()
	assert.True(t, cfg.Enabled)
	assert.Equal(t, []string{"DB_HOST", "DB_PORT", "API_KEY"}, cfg.Keys)
}

func TestRequiredConfigFromEnv_FailFast(t *testing.T) {
	t.Setenv("VAULTPULL_REQUIRED_KEYS", "SOME_KEY")
	t.Setenv("VAULTPULL_REQUIRED_FAIL_FAST", "true")

	cfg := RequiredConfigFromEnv()
	assert.True(t, cfg.FailFast)
}

func TestRequiredConfigFromEnv_NumericEnabled(t *testing.T) {
	t.Setenv("VAULTPULL_REQUIRED_KEYS", "SOME_KEY")
	t.Setenv("VAULTPULL_REQUIRED_FAIL_FAST", "1")

	cfg := RequiredConfigFromEnv()
	assert.True(t, cfg.FailFast)
}
